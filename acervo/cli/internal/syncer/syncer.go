package syncer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	mdImageRegex   = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	wikiImageRegex = regexp.MustCompile(`!\[\[(.*?)\]\]`)
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

			var rawAttachments []map[string]interface{}
			if atts, ok := tempMap["attachments"].([]interface{}); ok {
				for _, a := range atts {
					if m, ok := a.(map[string]interface{}); ok {
						rawAttachments = append(rawAttachments, m)
					}
				}
			}

			body := string(parts[2])

			bodyImages := make(map[string]bool)

			mdMatches := mdImageRegex.FindAllStringSubmatch(body, -1)
			for _, match := range mdMatches {
				imgSrc := match[2]
				imgSrc = filepath.Base(imgSrc)
				bodyImages[imgSrc] = true
			}

			wikiMatches := wikiImageRegex.FindAllStringSubmatch(body, -1)
			for _, match := range wikiMatches {
				imgSrc := match[1]
				bodyImages[imgSrc] = true
			}

			if mode == "yaml-to-body" {
				updatedBody := body
				changed := false
				for _, att := range rawAttachments {
					typeStr, _ := att["type"].(string)
					srcStr, _ := att["src"].(string)
					if typeStr == "image" && srcStr != "" {
						if !bodyImages[srcStr] {
							if !changed {
								updatedBody = strings.TrimRight(updatedBody, "\n") + "\n\n"
							}
							caption, _ := att["caption"].(string)
							if caption == "" {
								caption = "Image"
							}
							updatedBody += fmt.Sprintf("![%s](../../media/images/%s)\n", caption, srcStr)
							changed = true
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
				var newAttachments []*yaml.Node
				yamlImages := make(map[string]bool)

				// Rebuild using yaml.Node to preserve everything perfectly
				if len(rootNode.Content) > 0 && rootNode.Content[0].Kind == yaml.MappingNode {
					mapNode := rootNode.Content[0]

					var attachmentsNode *yaml.Node
					var attachmentsKeyIdx int = -1

					for i := 0; i < len(mapNode.Content); i += 2 {
						keyNode := mapNode.Content[i]
						if keyNode.Value == "attachments" {
							attachmentsNode = mapNode.Content[i+1]
							attachmentsKeyIdx = i
							break
						}
					}

					if attachmentsNode != nil && attachmentsNode.Kind == yaml.SequenceNode {
						for _, itemNode := range attachmentsNode.Content {
							// Check if it's an image
							isImage := false
							src := ""
							if itemNode.Kind == yaml.MappingNode {
								for j := 0; j < len(itemNode.Content); j += 2 {
									k := itemNode.Content[j].Value
									v := itemNode.Content[j+1].Value
									if k == "type" && v == "image" {
										isImage = true
									}
									if k == "src" {
										src = v
									}
								}
							}

							if isImage && src != "" {
								yamlImages[src] = true
								if bodyImages[src] {
									newAttachments = append(newAttachments, itemNode)
								}
							} else {
								// keep non-image attachments exactly as they are
								newAttachments = append(newAttachments, itemNode)
							}
						}
					}

					// Add missing body images
					changed := false
					for imgName := range bodyImages {
						if !yamlImages[imgName] {
							// Create new mapping node for this image
							newNode := &yaml.Node{Kind: yaml.MappingNode}
							newNode.Content = append(newNode.Content,
								&yaml.Node{Kind: yaml.ScalarNode, Value: "caption"},
								&yaml.Node{Kind: yaml.ScalarNode, Value: "Image"},
								&yaml.Node{Kind: yaml.ScalarNode, Value: "role"},
								&yaml.Node{Kind: yaml.ScalarNode, Value: "documentation"},
								&yaml.Node{Kind: yaml.ScalarNode, Value: "src"},
								&yaml.Node{Kind: yaml.ScalarNode, Value: imgName},
								&yaml.Node{Kind: yaml.ScalarNode, Value: "type"},
								&yaml.Node{Kind: yaml.ScalarNode, Value: "image"},
							)
							newAttachments = append(newAttachments, newNode)
							changed = true
						}
					}

					if (attachmentsNode != nil && len(attachmentsNode.Content) != len(newAttachments)) || changed {
						if len(newAttachments) == 0 && attachmentsKeyIdx != -1 {
							// Remove attachments entirely
							mapNode.Content = append(mapNode.Content[:attachmentsKeyIdx], mapNode.Content[attachmentsKeyIdx+2:]...)
						} else {
							if attachmentsNode == nil {
								// Need to add 'attachments' key
								mapNode.Content = append(mapNode.Content,
									&yaml.Node{Kind: yaml.ScalarNode, Value: "attachments"},
									&yaml.Node{Kind: yaml.SequenceNode, Content: newAttachments},
								)
							} else {
								attachmentsNode.Content = newAttachments
							}
						}

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
