package auditor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"acervo/internal/domain"
	"acervo/internal/utils"

	"gopkg.in/yaml.v3"
)

func Audit(entitiesDir string, imagesDir string) error {
	agents := make(map[string]bool)
	works := make(map[string]bool)
	actions := make(map[string]bool)
	var allActions []domain.Action
	var allWorks []domain.Work
	referencedImages := make(map[string]bool)

	// Single pass walk to collect all entities and actions
	err := filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			return nil
		}

		relPath, err := filepath.Rel(entitiesDir, path)
		if err != nil {
			fmt.Printf("[WARNING] Falha ao obter caminho relativo para %s: %v\n", path, err)
			return nil
		}
		parentDir := filepath.Base(filepath.Dir(relPath))

		if parentDir == "agents" {
			var agent struct {
				ID string `yaml:"id"`
			}
			if err := yaml.Unmarshal(parts[1], &agent); err == nil {
				agents[agent.ID] = true
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Agent em %s: %v\n", path, err)
			}
		} else if parentDir == "works" {
			var work domain.Work
			if err := yaml.Unmarshal(parts[1], &work); err == nil {
				works[work.ID] = true
				allWorks = append(allWorks, work)
				for _, att := range work.Attachments {
					if att.Type == "image" && att.Src != "" {
						referencedImages[att.Src] = true
					}
				}
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Work em %s: %v\n", path, err)
			}
		} else if parentDir == "actions" {
			var action domain.Action
			if err := yaml.Unmarshal(parts[1], &action); err == nil {
				actions[action.ID] = true
				allActions = append(allActions, action)
				for _, att := range action.Attachments {
					if att.Type == "image" && att.Src != "" {
						referencedImages[att.Src] = true
					}
				}
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Action em %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	brokenLinks := 0
	// Validate relations using collected data
	for _, action := range allActions {
		pb := utils.CleanWikilink(action.PerformedBy)
		if pb != "" && !agents[pb] {
			fmt.Printf("[BROKEN LINK] Action %s aponta para Agent %s inexistente\n", action.ID, pb)
			brokenLinks++
		}

		wid := utils.CleanWikilink(action.WorkID)
		if wid != "" && !works[wid] {
			fmt.Printf("[BROKEN LINK] Action %s aponta para Work %s inexistente\n", action.ID, wid)
			brokenLinks++
		}
	}

	// Image Auditing
	fmt.Println("--- IMAGE AUDIT ---")
	missingImages := 0
	for img := range referencedImages {
		imgPath := filepath.Join(imagesDir, img)
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			fmt.Printf("[MISSING IMAGE] Imagem referenciada não existe: %s\n", img)
			missingImages++
		}
	}

	danglingImages := 0
	namingViolations := 0
	imageFiles, err := os.ReadDir(imagesDir)
	if err == nil {
		for _, f := range imageFiles {
			if f.IsDir() {
				continue
			}
			name := f.Name()

			// 1. Check if referenced
			if !referencedImages[name] {
				fmt.Printf("[DANGLING IMAGE] Imagem não referenciada por nenhuma entidade: %s\n", name)
				danglingImages++
			}

			// 2. Check naming convention
			validPrefix := false
			expectedPrefixes := []string{"action-", "work-", "agent-"}
			for _, p := range expectedPrefixes {
				if strings.HasPrefix(name, p) {
					validPrefix = true
					// Check if it matches a known entity ID
					parts := strings.Split(name, "-")
					if len(parts) >= 2 {
						// This is heuristic, but let's see if we can find a matching ID
						// e.g. action-slug-001.jpeg -> action-slug
						// Find longest matching prefix that is an ID
						match := false
						currentPrefix := ""
						for i := 0; i < len(parts)-1; i++ {
							if i == 0 {
								currentPrefix = parts[i]
							} else {
								currentPrefix += "-" + parts[i]
							}
							if p == "action-" && actions[currentPrefix] {
								match = true
								break
							}
							if p == "work-" && works[currentPrefix] {
								match = true
								break
							}
							if p == "agent-" && agents[currentPrefix] {
								match = true
								break
							}
						}

						// Special case: if name starts with a valid ID entirely (minus extension)
						idWithoutExt := strings.TrimSuffix(name, filepath.Ext(name))
						if actions[idWithoutExt] || works[idWithoutExt] || agents[idWithoutExt] {
							match = true
						}

						if !match {
							// Try to see if it's a "test" or "doc" or something
							// We'll be strict for now as requested
							fmt.Printf("[NAMING WARNING] Imagem %s não parece seguir o ID de nenhuma entidade conhecida\n", name)
							namingViolations++
						}
					}
					break
				}
			}
			if !validPrefix {
				fmt.Printf("[NAMING ERROR] Imagem %s não possui prefixo válido (action-, work-, agent-)\n", name)
				namingViolations++
			}
		}
	}

	fmt.Printf("=== AUDIT REPORT ===\n")
	fmt.Printf("Broken Links: %d\n", brokenLinks)
	fmt.Printf("Missing Images: %d\n", missingImages)
	fmt.Printf("Dangling Images: %d\n", danglingImages)
	fmt.Printf("Naming Violations: %d\n", namingViolations)

	if brokenLinks > 0 || missingImages > 0 {
		return fmt.Errorf("audit failed with reference errors")
	}
	return nil
}
