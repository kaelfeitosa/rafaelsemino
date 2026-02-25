package validator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"acervo/internal/domain"

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

			// Robust entity type detection based on parent directory name
			rel, err := filepath.Rel(entitiesDir, path)
			if err != nil {
				return fmt.Errorf("ERRO: falha ao obter caminho relativo para %s: %w", path, err)
			}
			parentDir := filepath.Base(filepath.Dir(rel))

			if parentDir == "agents" {
				var agent domain.Agent
				if err := yaml.Unmarshal(parts[1], &agent); err != nil {
					return fmt.Errorf("ERRO: Agent em %s tem YAML inválido: %w", path, err)
				}
				if agent.ID == "" {
					return fmt.Errorf("ERRO: Agent %s sem id", path)
				}
				if agent.Name == "" {
					return fmt.Errorf("ERRO: Agent %s sem name", agent.ID)
				}
				if agent.Kind != "person" && agent.Kind != "collective" {
					return fmt.Errorf("ERRO: Agent %s com kind inválido: '%s'. Deve ser 'person' ou 'collective'", agent.ID, agent.Kind)
				}
			} else if parentDir == "works" {
				var work domain.Work
				if err := yaml.Unmarshal(parts[1], &work); err != nil {
					return fmt.Errorf("ERRO: Work em %s tem YAML inválido: %w", path, err)
				}
				if work.ID == "" {
					return fmt.Errorf("ERRO: Work %s sem id", path)
				}
				if work.Title == "" {
					return fmt.Errorf("ERRO: Work %s sem title", work.ID)
				}
				validWorkTypes := map[string]bool{
					"teatro": true, "jogo": true, "filme": true, "roteiro": true, "performance": true, "outro": true,
				}
				if work.Type == "" || !validWorkTypes[work.Type] {
					return fmt.Errorf("ERRO: Work %s tem tipo inválido: '%s'. Tipos permitidos: teatro | jogo | filme | roteiro | performance | outro", work.ID, work.Type)
				}
			} else if parentDir == "actions" {
				var action domain.Action
				if err := yaml.Unmarshal(parts[1], &action); err != nil {
					return fmt.Errorf("ERRO: Action em %s tem YAML inválido: %w", path, err)
				}
				if action.ID == "" {
					return fmt.Errorf("ERRO: Action %s sem id", path)
				}
				if action.Title == "" {
					return fmt.Errorf("ERRO: Action %s sem title", action.ID)
				}
				validActionTypes := map[string]bool{
					"criacao": true, "exibicao": true, "formacao": true, "avaliacao": true, "curadoria": true, "premiacao": true, "outro": true,
				}
				if action.Type == "" || !validActionTypes[action.Type] {
					return fmt.Errorf("ERRO: Action %s tem type inválido: '%s'. Tipos permitidos: criacao | exibicao | formacao | avaliacao | curadoria | premiacao | outro", action.ID, action.Type)
				}
				validContextKinds := map[string]bool{
					"festival": true, "mostra": true, "curso": true, "oficina": true, "residencia": true, "premiacao": true, "entrevista": true, "outro": true,
				}
				if action.Kind != "" && !validContextKinds[action.Kind] {
					return fmt.Errorf("ERRO: Action %s tem kind inválido: '%s'. Tipos permitidos: festival | mostra | curso | oficina | residencia | premiacao | entrevista | outro", action.ID, action.Kind)
				}
				if action.PerformedBy == "" {
					return fmt.Errorf("ERRO: Action %s sem performed_by", action.ID)
				}
				if action.MyRole == "" {
					return fmt.Errorf("ERRO: Action %s sem my_role", action.ID)
				}
				if action.Label == "" {
					return fmt.Errorf("ERRO: Action %s sem label", action.ID)
				}
				if action.DateStart == "" {
					return fmt.Errorf("ERRO: Action %s sem date_start", action.ID)
				}
				dateRegex := regexp.MustCompile(`^\d{4}(-\d{2})?(-\d{2})?$`)
				if !dateRegex.MatchString(action.DateStart) {
					return fmt.Errorf("ERRO: Action %s tem formato de date_start inválido: '%s'", action.ID, action.DateStart)
				}
				if action.DateEnd != "" {
					if !dateRegex.MatchString(action.DateEnd) {
						return fmt.Errorf("ERRO: Action %s tem formato de date_end inválido: '%s'", action.ID, action.DateEnd)
					}
				}
			}
		}
		return nil
	})
}
