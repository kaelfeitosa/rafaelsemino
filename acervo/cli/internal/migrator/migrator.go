package migrator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func MigrateAttachments(entitiesDir string) error {
	return filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			parts := bytes.SplitN(content, []byte("---"), 3)
			if len(parts) < 3 {
				return nil // Skip files without valid frontmatter
			}

			var rootNode yaml.Node
			if err := yaml.Unmarshal(parts[1], &rootNode); err != nil {
				return fmt.Errorf("error parsing YAML in %s: %w", path, err)
			}

			if len(rootNode.Content) > 0 && rootNode.Content[0].Kind == yaml.MappingNode {
				mapNode := rootNode.Content[0]
				var newContentNodes []*yaml.Node

				var oldAttachments []map[string]interface{}

				// Identify and extract the 'attachments' key
				for i := 0; i < len(mapNode.Content); i += 2 {
					keyNode := mapNode.Content[i]
					valNode := mapNode.Content[i+1]

					if keyNode.Value == "attachments" {
						// Extract to simple map
						var attList []interface{}
						if err := valNode.Decode(&attList); err == nil {
							for _, a := range attList {
								if amap, ok := a.(map[string]interface{}); ok {
									oldAttachments = append(oldAttachments, amap)
								}
							}
						}
					} else {
						// Also check for Occurrences that might have nested attachments
						if keyNode.Value == "occurrences" && valNode.Kind == yaml.SequenceNode {
							for _, occNode := range valNode.Content {
								if occNode.Kind == yaml.MappingNode {
									var newOccNodes []*yaml.Node
									var occAttachments []map[string]interface{}

									for j := 0; j < len(occNode.Content); j += 2 {
										oKey := occNode.Content[j]
										oVal := occNode.Content[j+1]

										if oKey.Value == "attachments" {
											var occAttList []interface{}
											if err := oVal.Decode(&occAttList); err == nil {
												for _, a := range occAttList {
													if amap, ok := a.(map[string]interface{}); ok {
														occAttachments = append(occAttachments, amap)
													}
												}
											}
										} else {
											newOccNodes = append(newOccNodes, oKey, oVal)
										}
									}

									// Rebuild occurrence with flattened attachment
									for idx, att := range occAttachments {
										seqIdx := idx + 1
										if t, ok := att["type"].(string); ok && t != "" {
											newOccNodes = append(newOccNodes,
												&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_type", seqIdx)},
												&yaml.Node{Kind: yaml.ScalarNode, Value: t})
										}
										if u, ok := att["url"].(string); ok && u != "" {
											newOccNodes = append(newOccNodes,
												&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_url", seqIdx)},
												&yaml.Node{Kind: yaml.ScalarNode, Value: u})
										}
										if l, ok := att["label"].(string); ok && l != "" {
											newOccNodes = append(newOccNodes,
												&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_label", seqIdx)},
												&yaml.Node{Kind: yaml.ScalarNode, Value: l})
										}
										if c, ok := att["caption"].(string); ok && c != "" {
											newOccNodes = append(newOccNodes,
												&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_label", seqIdx)},
												&yaml.Node{Kind: yaml.ScalarNode, Value: c})
										}
										if c, ok := att["category"].(string); ok && c != "" {
											newOccNodes = append(newOccNodes,
												&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_category", seqIdx)},
												&yaml.Node{Kind: yaml.ScalarNode, Value: c})
										}
									}
									occNode.Content = newOccNodes
								}
							}
						}

						newContentNodes = append(newContentNodes, keyNode, valNode)
					}
				}

				// Append flattened attachments at root level
				for idx, att := range oldAttachments {
					seqIdx := idx + 1
					if t, ok := att["type"].(string); ok && t != "" {
						newContentNodes = append(newContentNodes,
							&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_type", seqIdx)},
							&yaml.Node{Kind: yaml.ScalarNode, Value: t})
					}
					if u, ok := att["url"].(string); ok && u != "" {
						newContentNodes = append(newContentNodes,
							&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_url", seqIdx)},
							&yaml.Node{Kind: yaml.ScalarNode, Value: u})
					}
					if l, ok := att["label"].(string); ok && l != "" {
						newContentNodes = append(newContentNodes,
							&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_label", seqIdx)},
							&yaml.Node{Kind: yaml.ScalarNode, Value: l})
					}
					if c, ok := att["caption"].(string); ok && c != "" {
						newContentNodes = append(newContentNodes,
							&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_label", seqIdx)},
							&yaml.Node{Kind: yaml.ScalarNode, Value: c})
					}
					if c, ok := att["category"].(string); ok && c != "" {
						newContentNodes = append(newContentNodes,
							&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_category", seqIdx)},
							&yaml.Node{Kind: yaml.ScalarNode, Value: c})
					}
				}

				if len(oldAttachments) > 0 {
					mapNode.Content = newContentNodes

					var out bytes.Buffer
					encoder := yaml.NewEncoder(&out)
					encoder.SetIndent(2)
					if err := encoder.Encode(&rootNode); err != nil {
						return fmt.Errorf("failed to encode migrated frontmatter for %s: %w", path, err)
					}
					encoder.Close()

					newYaml := out.String()
					newFileContent := "---\n" + newYaml + "---" + string(parts[2])
					if err := os.WriteFile(path, []byte(newFileContent), 0644); err != nil {
						return fmt.Errorf("failed to write migrated file %s: %w", path, err)
					}
					fmt.Printf("✅ Migrated attachments in %s\n", path)
				}
			}
		}
		return nil
	})
}
