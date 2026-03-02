package syncer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	mdImageRegex   = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	wikiImageRegex = regexp.MustCompile(`!\[\[(.*?)\]\]`)
	attKeyRegex    = regexp.MustCompile(`^attachment_(\d+)_(.+)$`)
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

			// We need the raw map to easily check types/src without navigating the complex Node tree
			var tempMap map[string]interface{}
			yaml.Unmarshal(parts[1], &tempMap)

			// Extract flattened attachments
			type flatAttachment struct {
				idx      int
				Type     string
				URL      string
				Label    string
				Category string
			}

			rawAtts := make(map[int]*flatAttachment)
			for k, v := range tempMap {
				matches := attKeyRegex.FindStringSubmatch(k)
				if len(matches) == 3 {
					idx, _ := strconv.Atoi(matches[1])
					field := matches[2]
					if rawAtts[idx] == nil {
						rawAtts[idx] = &flatAttachment{idx: idx}
					}
					vStr := fmt.Sprintf("%v", v)
					switch field {
					case "type":
						rawAtts[idx].Type = vStr
					case "url":
						rawAtts[idx].URL = vStr
					case "label", "caption": // fallback caption
						rawAtts[idx].Label = vStr
					case "category":
						rawAtts[idx].Category = vStr
					}
				}
			}

			var rawAttachments []*flatAttachment
			for _, att := range rawAtts {
				rawAttachments = append(rawAttachments, att)
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
				// Normalize Obsidian embed targets (e.g., "folder/image.jpg|300")
				if idx := strings.Index(imgSrc, "|"); idx != -1 {
					imgSrc = imgSrc[:idx]
				}
				bodyImages[imgSrc] = true
			}

			if mode == "yaml-to-body" {
				updatedBody := body
				changed := false
				for _, att := range rawAttachments {
					if att.Type == "image" && att.URL != "" {
						if !bodyImages[att.URL] {
							if !changed {
								updatedBody = strings.TrimRight(updatedBody, "\n") + "\n\n"
							}
							caption := att.Label
							if caption == "" {
								caption = "Image"
							}
							updatedBody += fmt.Sprintf("![%s](../../media/images/%s)\n", caption, att.URL)
							changed = true
							bodyImages[att.URL] = true // Prevent duplicates from multiple YAML entries for the same src
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

					var newContentNodes []*yaml.Node
					var yamlImages = make(map[string]bool)

					// First pass: identify which existing image attachments should be kept
					// Non-image attachments are always kept.
					// We will rebuild the list of valid attachments to re-index them properly.
					validAtts := make([]*flatAttachment, 0)

					// Iterate index from 1 to 1000 to maintain order
					maxIndexFound := 0
					for idx := 1; idx < 1000; idx++ {
						if att, ok := rawAtts[idx]; ok {
							maxIndexFound = idx
							isImage := att.Type == "image"

							if isImage && att.URL != "" {
								if bodyImages[att.URL] {
									validAtts = append(validAtts, att)
									yamlImages[att.URL] = true
								}
							} else {
								// Keep non-image
								validAtts = append(validAtts, att)
							}
						}
					}

					// Second pass: append missing body images
					changed := false
					for imgName := range bodyImages {
						if !yamlImages[imgName] {
							maxIndexFound++
							validAtts = append(validAtts, &flatAttachment{
								idx: maxIndexFound,
								Type: "image",
								URL: imgName,
								Label: "Image",
								Category: "documentation",
							})
							changed = true
						}
					}

					// Rebuild the map. Copy over non-attachment keys.
					for i := 0; i < len(mapNode.Content); i += 2 {
						keyNode := mapNode.Content[i]
						valNode := mapNode.Content[i+1]

						if !attKeyRegex.MatchString(keyNode.Value) {
							newContentNodes = append(newContentNodes, keyNode, valNode)
						}
					}

					// Re-add the flattened keys with sequential indexing based on validAtts
					for newIdx, att := range validAtts {
						seqIdx := newIdx + 1

						if att.Type != "" {
							newContentNodes = append(newContentNodes,
								&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_type", seqIdx)},
								&yaml.Node{Kind: yaml.ScalarNode, Value: att.Type},
							)
						}
						if att.URL != "" {
							newContentNodes = append(newContentNodes,
								&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_url", seqIdx)},
								&yaml.Node{Kind: yaml.ScalarNode, Value: att.URL},
							)
						}
						if att.Label != "" {
							newContentNodes = append(newContentNodes,
								&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_label", seqIdx)},
								&yaml.Node{Kind: yaml.ScalarNode, Value: att.Label},
							)
						}
						if att.Category != "" {
							newContentNodes = append(newContentNodes,
								&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_category", seqIdx)},
								&yaml.Node{Kind: yaml.ScalarNode, Value: att.Category},
							)
						}
					}

					// Check if total count of keys changed, or if any missing were added
					numOldAttKeys := 0
					for k := range tempMap {
						if attKeyRegex.MatchString(k) {
							numOldAttKeys++
						}
					}

					numNewAttKeys := len(newContentNodes) - (len(mapNode.Content) - numOldAttKeys)

					if numNewAttKeys != numOldAttKeys || changed {
						mapNode.Content = newContentNodes

						var out bytes.Buffer
						encoder := yaml.NewEncoder(&out)
						encoder.SetIndent(2)
						if err := encoder.Encode(&rootNode); err != nil {
							return fmt.Errorf("failed to encode frontmatter for %s: %w", path, err)
						}
						encoder.Close()

						newYaml := out.String()
						newContent := "---\n" + newYaml + "---" + string(parts[2])
						if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
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
