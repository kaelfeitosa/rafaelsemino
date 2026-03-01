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

	validWorkMediums = map[string]bool{
		"teatro": true, "audiovisual": true, "pesquisa_academica": true, "pesquisa_artistica": true, "pesquisa": true, "ensino": true, "formacao": true, "exposicao": true, "performance": true, "jogo": true, "livro": true, "empresa": true, "outro": true, "producao_cultural": true, "curadoria_e_juri": true, "cultura_popular": true, "ensino_e_formacao": true, "literatura": true,
	}
	validOccurrenceTypes = map[string]bool{
		"apresentacao": true, "residencia": true, "oficina": true, "publicacao_ou_apresentacao": true, "lancamento": true, "premio": true, "exposicao": true,
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
				if work.Medium == "" || !validWorkMediums[work.Medium] {
					return fmt.Errorf("ERRO: Work %s tem medium inválido: '%s'. Tipos permitidos: %s", work.ID, work.Medium, getKeys(validWorkMediums))
				}

				// Optional: validate occurrences inside works
				dateRegex := regexp.MustCompile(`^\d{4}(-\d{2})?(-\d{2})?$`)
				for i, occ := range work.Occurrences {
					if occ.Type != "" && !validOccurrenceTypes[occ.Type] {
						return fmt.Errorf("ERRO: Occurrence %d no Work %s tem type inválido: '%s'", i+1, work.ID, occ.Type)
					}
					if occ.StartDate != "" && !dateRegex.MatchString(occ.StartDate) {
						return fmt.Errorf("ERRO: Occurrence %d no Work %s tem formato de start_date inválido: '%s'", i+1, work.ID, occ.StartDate)
					}
					if occ.EndDate != "" && !dateRegex.MatchString(occ.EndDate) {
						return fmt.Errorf("ERRO: Occurrence %d no Work %s tem formato de end_date inválido: '%s'", i+1, work.ID, occ.EndDate)
					}
				}

				if err := validateAttachmentsSync(path, work.Attachments, parts[2]); err != nil {
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
		if idx := strings.Index(imgSrc, "media/images/"); idx != -1 {
			imgSrc = imgSrc[idx+len("media/images/"):]
		} else if strings.HasPrefix(imgSrc, "../") || strings.HasPrefix(imgSrc, "./") {
			imgSrc = filepath.Base(imgSrc)
		}
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
		bodyImages[imgSrc] = true
	}

	yamlImages := make(map[string]bool)
	for _, att := range attachments {
		if att.Type == "image" && att.URL != "" {
			yamlImages[att.URL] = true
			if !bodyImages[att.URL] {
				return fmt.Errorf("ERRO (%s): Imagem '%s' está em attachments YAML mas não no corpo do Markdown. Execute 'acervo sync-images --mode=yaml-to-body'", path, att.URL)
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
