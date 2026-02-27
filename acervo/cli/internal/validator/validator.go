package validator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"acervo/internal/domain"

	"gopkg.in/yaml.v3"
)

var (
	mdImageRegex   = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	wikiImageRegex = regexp.MustCompile(`!\[\[(.*?)\]\]`)

	validWorkTypes = map[string]bool{
		"teatro": true, "jogo": true, "filme": true, "roteiro": true, "performance": true, "outro": true,
	}
	validActionCategories = map[string]bool{
		"criacao": true, "exibicao": true, "formacao": true, "avaliacao": true, "curadoria": true, "premiacao": true, "outro": true,
	}
	validActionFormats = map[string]bool{
		"festival": true, "mostra": true, "curso": true, "oficina": true, "residencia": true, "premiacao": true, "entrevista": true, "outro": true,
	}
)

func getKeys(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, " | ")
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
				if err := validateAttachmentsSync(path, agent.Attachments, parts[2]); err != nil {
					return err
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
				if work.Type == "" || !validWorkTypes[work.Type] {
					return fmt.Errorf("ERRO: Work %s tem tipo inválido: '%s'. Tipos permitidos: %s", work.ID, work.Type, getKeys(validWorkTypes))
				}
				if err := validateAttachmentsSync(path, work.Attachments, parts[2]); err != nil {
					return err
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
				if action.Category == "" || !validActionCategories[action.Category] {
					return fmt.Errorf("ERRO: Action %s tem category inválido: '%s'. Tipos permitidos: %s", action.ID, action.Category, getKeys(validActionCategories))
				}
				if action.Format == "" || !validActionFormats[action.Format] {
					return fmt.Errorf("ERRO: Action %s tem format inválido: '%s'. Tipos permitidos: %s", action.ID, action.Format, getKeys(validActionFormats))
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
				if err := validateAttachmentsSync(path, action.Attachments, parts[2]); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func validateAttachmentsSync(path string, attachments []domain.Attachment, body []byte) error {
	bodyStr := string(body)
	bodyImages := make(map[string]bool)

	// Standard markdown images: ![alt](path)
	mdMatches := mdImageRegex.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range mdMatches {
		imgSrc := match[2]
		imgSrc = filepath.Base(imgSrc)
		bodyImages[imgSrc] = true
	}

	// Obsidian Wiki links: ![[image.jpg]]
	wikiMatches := wikiImageRegex.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range wikiMatches {
		imgSrc := match[1]
		// Normalize Obsidian embed targets (e.g., "folder/image.jpg|300")
		if idx := strings.Index(imgSrc, "|"); idx != -1 {
			imgSrc = imgSrc[:idx]
		}
		imgSrc = filepath.Base(imgSrc)
		bodyImages[imgSrc] = true
	}

	yamlImages := make(map[string]bool)
	for _, att := range attachments {
		if att.Type == "image" && att.Src != "" {
			yamlImages[att.Src] = true
			if !bodyImages[att.Src] {
				return fmt.Errorf("ERRO (%s): Imagem '%s' está em attachments YAML mas não no corpo do Markdown. Execute 'acervo sync-images --mode=yaml-to-body'", path, att.Src)
			}
		}
	}

	for imgSrc := range bodyImages {
		if !yamlImages[imgSrc] {
			return fmt.Errorf("ERRO (%s): Imagem '%s' está no corpo do Markdown mas não no attachments YAML. Execute 'acervo sync-images --mode=body-to-yaml'", path, imgSrc)
		}
	}

	return nil
}
