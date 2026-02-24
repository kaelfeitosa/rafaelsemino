package validator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ValidateEntities(entitiesDir string) error {
	return filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			parts := bytes.SplitN(content, []byte("---"), 3)
			if len(parts) < 3 {
				return fmt.Errorf("ERRO: %s não possui frontmatter YAML válido", path)
			}

			var base struct {
				ID   string `yaml:"id"`
				Type string `yaml:"type,omitempty"` // Legacy check, can be removed or kept for safety
			}
			if err := yaml.Unmarshal(parts[1], &base); err != nil {
				return fmt.Errorf("ERRO: %s tem YAML inválido: %w", path, err)
			}

			if base.ID == "" {
				return fmt.Errorf("ERRO: %s sem id", path)
			}

			// Identify type by folder or content structure
			if strings.Contains(path, "/agents/") {
				var agent struct {
					Name string `yaml:"name"`
					Kind string `yaml:"kind"`
				}
				if err := yaml.Unmarshal(parts[1], &agent); err != nil {
					return err
				}
				if agent.Name == "" {
					return fmt.Errorf("ERRO: Agent %s sem name", base.ID)
				}
				if agent.Kind != "person" && agent.Kind != "collective" {
					// Soft warning or error? Enforce schema.
					// return fmt.Errorf("ERRO: Agent %s com kind inválido: %s", base.ID, agent.Kind)
				}
			} else if strings.Contains(path, "/works/") {
				var work struct {
					Title string `yaml:"title"`
				}
				if err := yaml.Unmarshal(parts[1], &work); err != nil {
					return err
				}
				if work.Title == "" {
					return fmt.Errorf("ERRO: Work %s sem title", base.ID)
				}
			} else if strings.Contains(path, "/actions/") {
				var action struct {
					Title       string `yaml:"title"`
					PerformedBy string `yaml:"performed_by"`
					MyRole      string `yaml:"my_role"`
					Context     struct {
						Label string `yaml:"label"`
					} `yaml:"context"`
				}
				if err := yaml.Unmarshal(parts[1], &action); err != nil {
					return err
				}
				if action.PerformedBy == "" {
					return fmt.Errorf("ERRO: Action %s sem performed_by", base.ID)
				}
				if action.MyRole == "" {
					return fmt.Errorf("ERRO: Action %s sem my_role", base.ID)
				}
				if action.Context.Label == "" {
					// Ideally context should be present
					// return fmt.Errorf("ERRO: Action %s sem context.label", base.ID)
				}
			}
		}
		return nil
	})
}
