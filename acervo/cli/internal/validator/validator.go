package validator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Entity struct {
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
}

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

			// Extract YAML frontmatter
			parts := bytes.SplitN(content, []byte("---"), 3)
			if len(parts) < 3 {
				return fmt.Errorf("ERRO: %s não possui frontmatter YAML válido", path)
			}

			var entity struct {
				ID    string `yaml:"id"`
				Type  string `yaml:"type"`
				Title string `yaml:"title"`
				Name  string `yaml:"name"`
				Date  string `yaml:"date"`
				Agent string `yaml:"agent"`
				Path  string `yaml:"path"`
				URL   string `yaml:"url"`
			}
			if err := yaml.Unmarshal(parts[1], &entity); err != nil {
				return fmt.Errorf("ERRO: %s tem YAML inválido: %w", path, err)
			}

			if entity.ID == "" || entity.Type == "" {
				return fmt.Errorf("ERRO: %s sem id ou type", path)
			}

			switch entity.Type {
			case "agent":
				if entity.Name == "" {
					return fmt.Errorf("ERRO: %s (agent) sem name", path)
				}
			case "event":
				if entity.Name == "" {
					return fmt.Errorf("ERRO: %s (event) sem name", path)
				}
			case "work":
				if entity.Title == "" {
					return fmt.Errorf("ERRO: %s (work) sem title", path)
				}
			case "participation":
				if entity.Agent == "" {
					return fmt.Errorf("ERRO: %s (participation) sem agent", path)
				}
			case "record":
				if entity.Path == "" && entity.URL == "" {
					return fmt.Errorf("ERRO: %s (record) sem path e sem url", path)
				}
			}
		}
		return nil
	})
}
