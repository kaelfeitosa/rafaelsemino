package auditor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"acervo/internal/domain"

	"gopkg.in/yaml.v3"
)

func cleanWikilink(s string) string {
	s = strings.TrimPrefix(s, "[[")
	s = strings.TrimSuffix(s, "]]")
	return s
}

func Audit(entitiesDir string) error {
	agents := make(map[string]bool)
	works := make(map[string]bool)
	var allActions []domain.Action

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

		if strings.Contains(path, "/agents/") {
			var agent struct {
				ID string `yaml:"id"`
			}
			if err := yaml.Unmarshal(parts[1], &agent); err == nil {
				agents[agent.ID] = true
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Agent em %s: %v\n", path, err)
			}
		} else if strings.Contains(path, "/works/") {
			var work struct {
				ID string `yaml:"id"`
			}
			if err := yaml.Unmarshal(parts[1], &work); err == nil {
				works[work.ID] = true
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Work em %s: %v\n", path, err)
			}
		} else if strings.Contains(path, "/actions/") {
			var action domain.Action
			if err := yaml.Unmarshal(parts[1], &action); err == nil {
				allActions = append(allActions, action)
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
		pb := cleanWikilink(action.PerformedBy)
		if pb != "" && !agents[pb] {
			fmt.Printf("[BROKEN LINK] Action %s aponta para Agent %s inexistente\n", action.ID, pb)
			brokenLinks++
		}

		wid := cleanWikilink(action.WorkID)
		if wid != "" && !works[wid] {
			fmt.Printf("[BROKEN LINK] Action %s aponta para Work %s inexistente\n", action.ID, wid)
			brokenLinks++
		}
	}

	fmt.Printf("=== AUDIT REPORT ===\nBroken Links: %d\n", brokenLinks)
	if brokenLinks > 0 {
		return fmt.Errorf("audit failed with broken links")
	}
	return nil
}
