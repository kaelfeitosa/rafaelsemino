package auditor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func cleanWikilink(s string) string {
	s = strings.TrimPrefix(s, "[[")
	s = strings.TrimSuffix(s, "]]")
	return s
}

func Audit(entitiesDir string) error {
	participations := make(map[string]bool)
	recordsWithRelated := make(map[string]bool)
	allEntities := make(map[string]bool)
	relations := make(map[string][]string) // source -> []targets

	filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			parts := bytes.SplitN(content, []byte("---"), 3)
			if len(parts) < 3 {
				return nil
			}

			var data map[string]interface{}
			if err := yaml.Unmarshal(parts[1], &data); err != nil {
				return nil
			}

			typ, _ := data["type"].(string)
			id, _ := data["id"].(string)

			if id != "" {
				allEntities[id] = true
			}

			if typ == "participation" {
				participations[id] = true
			}

			// Map related_to
			if rel, ok := data["related_to"]; ok {
				if relStr, okStr := rel.(string); okStr && relStr != "" {
					cleanRel := cleanWikilink(relStr)
					relations[id] = append(relations[id], cleanRel)
					if typ == "record" {
						recordsWithRelated[cleanRel] = true
					}
				} else if relList, okList := rel.([]interface{}); okList {
					for _, item := range relList {
						if target, isStr := item.(string); isStr && target != "" {
							cleanTarget := cleanWikilink(target)
							relations[id] = append(relations[id], cleanTarget)
							if typ == "record" {
								recordsWithRelated[cleanTarget] = true
							}
						}
					}
				}
			}
		}
		return nil
	})

	fmt.Println("=== ACERVO STRUCTURAL AUDIT ===")

	// 1. Broken Links Check
	brokenLinks := 0
	for sourceID, targets := range relations {
		for _, targetID := range targets {
			if !allEntities[targetID] {
				fmt.Printf("[BROKEN LINK] %s aponta para %s que não existe!\n", sourceID, targetID)
				brokenLinks++
			}
		}
	}

	// 2. Participations without records
	fmt.Println("\nParticipations sem records:")
	unrecorded := 0
	for pid := range participations {
		if !recordsWithRelated[pid] {
			fmt.Println(" -", pid)
			unrecorded++
		}
	}

	fmt.Printf("\nResultado Final: %d broken links, %d participations sem arquivo.\n", brokenLinks, unrecorded)

	if brokenLinks > 0 {
		return fmt.Errorf("o Grafo do Acervo contém links quebrados intransponíveis")
	}

	return nil
}
