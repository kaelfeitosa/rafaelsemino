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
	actions := make(map[string]bool)

	// Collect all IDs first
	filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		content, _ := os.ReadFile(path)
		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			return nil
		}

		// Simple ID extraction
		var base struct {
			ID string `yaml:"id"`
		}
		yaml.Unmarshal(parts[1], &base)

		if strings.Contains(path, "/agents/") {
			agents[base.ID] = true
		} else if strings.Contains(path, "/works/") {
			works[base.ID] = true
		} else if strings.Contains(path, "/actions/") {
			actions[base.ID] = true
		}
		return nil
	})

	brokenLinks := 0

	// Check relations
	err := filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		content, _ := os.ReadFile(path)
		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			return nil
		}

		if strings.Contains(path, "/actions/") {
			var action domain.Action
			if err := yaml.Unmarshal(parts[1], &action); err != nil {
				return nil
			}

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

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("=== AUDIT REPORT ===\nBroken Links: %d\n", brokenLinks)
	if brokenLinks > 0 {
		return fmt.Errorf("audit failed with broken links")
	}
	return nil
}
