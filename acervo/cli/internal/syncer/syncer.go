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

	return filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			rel, _ := filepath.Rel(entitiesDir, path)
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

			attMap := make(map[int]map[string]interface{})
			for k, v := range tempMap {
				matches := flattenedRe.FindStringSubmatch(k)
				if len(matches) == 3 {
					idx, err := strconv.Atoi(matches[1])
					if err != nil {
						continue
					}
					field := matches[2]

					if attMap[idx] == nil {
						attMap[idx] = make(map[string]interface{})
					}
					attMap[idx][field] = v
				}
			}

			if atts, ok := tempMap["attachments"].([]interface{}); ok {
				idx := 1
				for k := range attMap {
					if k >= idx {
						idx = k + 1
					}
				}
				for _, a := range atts {
					if m, ok := a.(map[string]interface{}); ok {
						attMap[idx] = m
						idx++
					}
				}
			}

			var indices []int
			for k := range attMap {
				indices = append(indices, k)
			}
			sort.Ints(indices)

			var rawAttachments []map[string]interface{}
			for _, idx := range indices {
				rawAttachments = append(rawAttachments, attMap[idx])
			}

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
				for _, att := range rawAttachments {
					typeStr, _ := att["type"].(string)
					srcStr, _ := att["url"].(string)
					if typeStr == "image" && srcStr != "" {
						if !bodyImages[srcStr] {
							if !changed {
								updatedBody = strings.TrimRight(updatedBody, "\n") + "\n\n"
							}
							caption, _ := att["label"].(string)
							if caption == "" {
								caption, _ = att["caption"].(string)
							}
							if caption == "" {
								caption = "Image"
							}
							updatedBody += fmt.Sprintf("![%s](../../media/images/%s)\n", caption, srcStr)
							changed = true
							bodyImages[srcStr] = true
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

					var newContent []*yaml.Node
					yamlImages := make(map[string]bool)

					// Rebuild the yaml map, keeping non-attachment keys
					// For attachments, we keep them if they are in bodyImages (or if they are not images)

					type attData struct {
						idx int
						nodes []*yaml.Node // key, value, key, value...
						url string
						isImage bool
					}

					parsedAtts := make(map[int]*attData)
					var otherNodes []*yaml.Node

					for i := 0; i < len(mapNode.Content); i += 2 {
						keyNode := mapNode.Content[i]
						valNode := mapNode.Content[i+1]

						if keyNode.Value == "attachments" {
							// For migration: parse legacy attachments block and integrate
							if valNode.Kind == yaml.SequenceNode {
								idx := 1
								for _, itemNode := range valNode.Content {
									// find next available idx
									for parsedAtts[idx] != nil { idx++ }

									if itemNode.Kind == yaml.MappingNode {
										ad := &attData{idx: idx, nodes: make([]*yaml.Node, 0)}
										for j := 0; j < len(itemNode.Content); j += 2 {
											k := itemNode.Content[j]
											v := itemNode.Content[j+1]
											if k.Value == "type" && v.Value == "image" { ad.isImage = true }
											if k.Value == "url" { ad.url = v.Value }
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

							if field == "type" && valNode.Value == "image" { ad.isImage = true }
							if field == "url" { ad.url = valNode.Value }

							ad.nodes = append(ad.nodes, &yaml.Node{Kind: yaml.ScalarNode, Value: field}, valNode)
						} else {
							otherNodes = append(otherNodes, keyNode, valNode)
						}
					}

					// Filter existing attachments
					var keptAtts []*attData
					for _, idx := range indices {
						ad := parsedAtts[idx]
						if ad == nil { continue }

						if ad.isImage && ad.url != "" {
							yamlImages[ad.url] = true
							if bodyImages[ad.url] {
								keptAtts = append(keptAtts, ad)
							}
						} else {
							keptAtts = append(keptAtts, ad)
						}
					}

					// Add missing body images
					changed := false
					// To avoid any gap issues, since they will be reindexed by keptAtts anyway,
					// we just need a unique high index to append. The sort will handle it.
					nextIdx := 99999
					for imgName := range bodyImages {
						if !yamlImages[imgName] {
							newAd := &attData{
								idx: nextIdx,
								url: imgName,
								isImage: true,
								nodes: []*yaml.Node{
									{Kind: yaml.ScalarNode, Value: "label"}, {Kind: yaml.ScalarNode, Value: "Image"},
									{Kind: yaml.ScalarNode, Value: "category"}, {Kind: yaml.ScalarNode, Value: "documentation"},
									{Kind: yaml.ScalarNode, Value: "url"}, {Kind: yaml.ScalarNode, Value: imgName},
									{Kind: yaml.ScalarNode, Value: "type"}, {Kind: yaml.ScalarNode, Value: "image"},
								},
							}
							keptAtts = append(keptAtts, newAd)
							nextIdx++
							changed = true
						}
					}

					// Re-index and reconstruct nodes
					sort.Slice(keptAtts, func(i, j int) bool { return keptAtts[i].idx < keptAtts[j].idx })

					for i, ad := range keptAtts {
						newIdx := i + 1
						if ad.idx != newIdx { changed = true } // Needs re-indexing

						for j := 0; j < len(ad.nodes); j += 2 {
							field := ad.nodes[j].Value
							valNode := ad.nodes[j+1]

							keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_%s", newIdx, field)}
							newContent = append(newContent, keyNode, valNode)
						}
					}

					// If there's a difference in length or elements, or explicitly changed
					needsSave := changed || len(parsedAtts) != len(keptAtts)
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
