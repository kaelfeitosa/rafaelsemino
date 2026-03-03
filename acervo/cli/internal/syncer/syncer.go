package syncer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"acervo/internal/domain"

	"gopkg.in/yaml.v3"
)

var (
	mdImageRegex   = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	wikiImageRegex = regexp.MustCompile(`!\[\[(.*?)\]\]`)
	flattenedRe    = regexp.MustCompile(`^attachment_(\d+)_(.+)$`)
)

func SyncImages(entitiesDir string, mode string) error {
	if mode != "yaml-to-body" && mode != "body-to-yaml" {
		return fmt.Errorf("invalid mode: %s. Must be 'yaml-to-body' or 'body-to-yaml'", mode)
	}

	// Sanitize against path traversal vulnerabilities
	safeEntitiesDir := filepath.Clean(entitiesDir)

	return filepath.Walk(safeEntitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			rel, _ := filepath.Rel(safeEntitiesDir, path)
			parentDir := filepath.Base(filepath.Dir(rel))
			if parentDir != "actions" && parentDir != "works" && parentDir != "agents" {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			parts := bytes.SplitN(content, []byte("---"), 3)
			if len(parts) < 3 {
				return nil
			}

			var rootNode yaml.Node
			if err := yaml.Unmarshal(parts[1], &rootNode); err != nil {
				return fmt.Errorf("error parsing YAML in %s: %w", path, err)
			}

			var tempMap map[string]interface{}
			yaml.Unmarshal(parts[1], &tempMap)

			// Simplify extraction using domain layer utility
			domainAttachments := domain.ExtractAttachments(tempMap)

			body := string(parts[2])

			bodyImages := make(map[string]bool)

			mdMatches := mdImageRegex.FindAllStringSubmatch(body, -1)
			for _, match := range mdMatches {
				imgSrc := match[2]
				if idx := strings.Index(imgSrc, "media/images/"); idx != -1 {
					imgSrc = imgSrc[idx+len("media/images/"):]
				} else if strings.HasPrefix(imgSrc, "../") || strings.HasPrefix(imgSrc, "./") {
					imgSrc = filepath.Base(imgSrc)
				}
				bodyImages[imgSrc] = true
			}

			wikiMatches := wikiImageRegex.FindAllStringSubmatch(body, -1)
			for _, match := range wikiMatches {
				imgSrc := match[1]
				if idx := strings.Index(imgSrc, "|"); idx != -1 {
					imgSrc = imgSrc[:idx]
				}
				bodyImages[imgSrc] = true
			}

			if mode == "yaml-to-body" {
				updatedBody := body
				changed := false
				for _, att := range domainAttachments {
					if att.Type == "image" && att.URL != "" {
						if !bodyImages[att.URL] {
							if !changed {
								updatedBody = strings.TrimRight(updatedBody, "\n") + "\n\n"
							}
							caption := att.Label
							if caption == "" {
								caption = "Image"
							}
							// Sanitize fields to prevent markdown injection
							safeCaption := strings.ReplaceAll(caption, "\n", " ")
							safeCaption = strings.ReplaceAll(safeCaption, "[", "\\[")
							safeCaption = strings.ReplaceAll(safeCaption, "]", "\\]")

							safeSrcStr := strings.ReplaceAll(att.URL, "\n", "")
							safeSrcStr = strings.ReplaceAll(safeSrcStr, ")", "%29")
							safeSrcStr = strings.ReplaceAll(safeSrcStr, "\"", "%22")
							safeSrcStr = strings.ReplaceAll(safeSrcStr, "<", "%3C")
							safeSrcStr = strings.ReplaceAll(safeSrcStr, ">", "%3E")

							updatedBody += fmt.Sprintf("![%s](../../media/images/%s)\n", safeCaption, safeSrcStr)
							changed = true
							bodyImages[att.URL] = true
						}
					}
				}
				if changed {
					newContent := string(parts[0]) + "---" + string(parts[1]) + "---" + updatedBody
					if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %w", path, err)
					}
					fmt.Printf("🔄 Updated body of %s\n", path)
				}
			} else if mode == "body-to-yaml" {
				if len(rootNode.Content) > 0 && rootNode.Content[0].Kind == yaml.MappingNode {
					mapNode := rootNode.Content[0]

					parsedAtts, otherNodes := parseAttachmentsFromNode(mapNode)
					keptAtts, yamlImages := filterKeptAttachments(parsedAtts, bodyImages)
					keptAtts, appendedChanged := appendMissingBodyImages(keptAtts, bodyImages, yamlImages)
					newContent, reindexChanged := reindexAttachments(keptAtts)

					// If there's a difference in length or elements, or explicitly changed
					needsSave := appendedChanged || reindexChanged || len(parsedAtts) != len(keptAtts)
					if attsList, ok := tempMap["attachments"].([]interface{}); ok && len(attsList) > 0 {
						needsSave = true
					}
					if needsSave {
						mapNode.Content = append(otherNodes, newContent...)

						var out bytes.Buffer
						encoder := yaml.NewEncoder(&out)
						encoder.SetIndent(2)
						if err := encoder.Encode(&rootNode); err != nil {
							return fmt.Errorf("failed to encode frontmatter for %s: %w", path, err)
						}
						encoder.Close()

						newYaml := out.String()
						newContentStr := "---\n" + newYaml + "---" + string(parts[2])
						if err := os.WriteFile(path, []byte(newContentStr), 0644); err != nil {
							return fmt.Errorf("failed to write %s: %w", path, err)
						}
						fmt.Printf("🔄 Updated YAML of %s\n", path)
					}
				}
			}
		}
		return nil
	})
}

type attData struct {
	idx     int
	nodes   []*yaml.Node
	url     string
	isImage bool
}

func parseAttachmentsFromNode(mapNode *yaml.Node) (map[int]*attData, []*yaml.Node) {
	parsedAtts := make(map[int]*attData)
	var otherNodes []*yaml.Node

	for i := 0; i < len(mapNode.Content); i += 2 {
		keyNode := mapNode.Content[i]
		valNode := mapNode.Content[i+1]

		if keyNode.Value == "attachments" {
			if valNode.Kind == yaml.SequenceNode {
				idx := 1
				for _, itemNode := range valNode.Content {
					for parsedAtts[idx] != nil {
						idx++
					}
					if itemNode.Kind == yaml.MappingNode {
						ad := &attData{idx: idx, nodes: make([]*yaml.Node, 0)}
						for j := 0; j < len(itemNode.Content); j += 2 {
							k := itemNode.Content[j]
							v := itemNode.Content[j+1]
							if k.Value == "type" && v.Value == "image" {
								ad.isImage = true
							}
							if k.Value == "url" {
								ad.url = v.Value
							}
							ad.nodes = append(ad.nodes, k, v)
						}
						parsedAtts[idx] = ad
					}
				}
			}
			continue
		}

		matches := flattenedRe.FindStringSubmatch(keyNode.Value)
		if len(matches) == 3 {
			idx, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			field := matches[2]

			if parsedAtts[idx] == nil {
				parsedAtts[idx] = &attData{idx: idx, nodes: make([]*yaml.Node, 0)}
			}
			ad := parsedAtts[idx]

			if field == "type" && valNode.Value == "image" {
				ad.isImage = true
			}
			if field == "url" {
				ad.url = valNode.Value
			}

			ad.nodes = append(ad.nodes, &yaml.Node{Kind: yaml.ScalarNode, Value: field}, valNode)
		} else {
			otherNodes = append(otherNodes, keyNode, valNode)
		}
	}
	return parsedAtts, otherNodes
}

func filterKeptAttachments(parsedAtts map[int]*attData, bodyImages map[string]bool) ([]*attData, map[string]bool) {
	var keptAtts []*attData
	yamlImages := make(map[string]bool)

	var parsedIndices []int
	for idx := range parsedAtts {
		parsedIndices = append(parsedIndices, idx)
	}
	sort.Ints(parsedIndices)

	for _, idx := range parsedIndices {
		ad := parsedAtts[idx]
		if ad == nil {
			continue
		}

		if ad.isImage && ad.url != "" {
			yamlImages[ad.url] = true
			if bodyImages[ad.url] {
				keptAtts = append(keptAtts, ad)
			}
		} else {
			keptAtts = append(keptAtts, ad)
		}
	}
	return keptAtts, yamlImages
}

func appendMissingBodyImages(keptAtts []*attData, bodyImages map[string]bool, yamlImages map[string]bool) ([]*attData, bool) {
	changed := false
	nextIdx := 0
	for _, ad := range keptAtts {
		if ad.idx > nextIdx {
			nextIdx = ad.idx
		}
	}
	nextIdx++
	for imgName := range bodyImages {
		if !yamlImages[imgName] {
			newAd := &attData{
				idx:     nextIdx,
				url:     imgName,
				isImage: true,
				nodes: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "label"}, {Kind: yaml.ScalarNode, Value: "Image"},
					{Kind: yaml.ScalarNode, Value: "category"}, {Kind: yaml.ScalarNode, Value: "registro"},
					{Kind: yaml.ScalarNode, Value: "url"}, {Kind: yaml.ScalarNode, Value: imgName, Style: yaml.DoubleQuotedStyle},
					{Kind: yaml.ScalarNode, Value: "type"}, {Kind: yaml.ScalarNode, Value: "image"},
				},
			}
			keptAtts = append(keptAtts, newAd)
			nextIdx++
			changed = true
		}
	}
	return keptAtts, changed
}

func reindexAttachments(keptAtts []*attData) ([]*yaml.Node, bool) {
	var newContent []*yaml.Node
	changed := false

	sort.Slice(keptAtts, func(i, j int) bool { return keptAtts[i].idx < keptAtts[j].idx })

	for i, ad := range keptAtts {
		newIdx := i + 1
		if ad.idx != newIdx {
			changed = true
		} // Needs re-indexing

		// Collect fields to apply precedence logic
		fields := make(map[string]*yaml.Node)
		for j := 0; j < len(ad.nodes); j += 2 {
			fields[ad.nodes[j].Value] = ad.nodes[j+1]
		}

		// Apply precedence: label > caption
		if _, hasLabel := fields["label"]; !hasLabel {
			if captionNode, hasCaption := fields["caption"]; hasCaption {
				fields["label"] = captionNode
				changed = true
			}
		}
		delete(fields, "caption")

		// Apply precedence: category > role
		if _, hasCategory := fields["category"]; !hasCategory {
			if roleNode, hasRole := fields["role"]; hasRole {
				fields["category"] = roleNode
				changed = true
			}
		}
		delete(fields, "role")

		// Rebuild nodes using a stable order (e.g., type, url, label, category, etc.)
		orderedKeys := []string{"type", "url", "label", "category"}
		writtenKeys := make(map[string]bool)

		for _, k := range orderedKeys {
			if v, exists := fields[k]; exists {
				keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_%s", newIdx, k)}
				newContent = append(newContent, keyNode, v)
				writtenKeys[k] = true
			}
		}

		// Write any remaining unknown fields
		for k, v := range fields {
			if !writtenKeys[k] {
				keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_%s", newIdx, k)}
				newContent = append(newContent, keyNode, v)
			}
		}
	}
	return newContent, changed
}
