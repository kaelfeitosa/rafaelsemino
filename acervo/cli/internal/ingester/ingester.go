package ingester

import (
	"acervo/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func Create(entityType string, slug string, args []string) error {
	// Security: Sanitize inputs to prevent path traversal
	entityType = filepath.Base(entityType)
	slug = filepath.Base(slug)

	// Validate entityType against allow-list
	allowedTypes := map[string]bool{"action": true, "work": true, "agent": true}
	if !allowedTypes[entityType] {
		return fmt.Errorf("tipo de entidade inválido: %s", entityType)
	}

	id := fmt.Sprintf("%s-%s", entityType, slug)
	var frontmatter interface{}
	var body string = "Detalhes específicos."

	// Helper to ensure wikilink format
	toWikilink := func(s string) string {
		s = strings.TrimSpace(s)
		if s == "" {
			return ""
		}
		if !strings.HasPrefix(s, "[[") {
			s = "[[" + s
		}
		if !strings.HasSuffix(s, "]]") {
			s = s + "]]"
		}
		return s
	}

	switch entityType {
	case "action":
		data := &domain.Action{
			ID: id,
		}
		for _, arg := range args {
			kv := strings.SplitN(arg, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			switch key {
			case "title":
				data.Title = val
			case "type":
				data.Type = val
			case "kind":
				data.Kind = val
			case "label":
				data.Label = val
			case "location":
				data.Location = val
			case "performed_by":
				data.PerformedBy = toWikilink(val)
			case "my_role":
				data.MyRole = val
			case "work_id":
				data.WorkID = toWikilink(val)
			case "context_label":
				data.Label = val
			case "context_kind":
				data.Kind = val
			case "context_location":
				data.Location = val
			case "date_start":
				data.DateStart = val
			case "date_end":
				data.DateEnd = val
			case "description":
				data.Description = val
			case "featured":
				data.Featured = (strings.ToLower(val) == "true")
			case "collaborators":
				if val != "" {
					if err := json.Unmarshal([]byte(val), &data.Collaborators); err != nil {
						return fmt.Errorf("falha ao analisar collaborators (esperado JSON): %w", err)
					}
				}
			case "attachments":
				if val != "" {
					if err := json.Unmarshal([]byte(val), &data.Attachments); err != nil {
						return fmt.Errorf("falha ao analisar attachments (esperado JSON): %w", err)
					}
				}
			}
		}
		frontmatter = data
	case "agent":
		data := &domain.Agent{
			ID: id,
		}
		for _, arg := range args {
			kv := strings.SplitN(arg, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			switch key {
			case "name":
				data.Name = val
			case "kind":
				data.Kind = val
			case "description":
				data.Description = val
			case "founded_by_me":
				data.FoundedByMe = (strings.ToLower(val) == "true")
			case "active_since":
				if val != "" {
					y, err := strconv.Atoi(val)
					if err != nil {
						return fmt.Errorf("valor inválido para active_since: '%s' não é um número", val)
					}
					data.ActiveSince = y
				}
			case "featured":
				data.Featured = (strings.ToLower(val) == "true")
			}
		}
		frontmatter = data
	case "work":
		data := &domain.Work{
			ID: id,
		}
		for _, arg := range args {
			kv := strings.SplitN(arg, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			switch key {
			case "title":
				data.Title = val
			case "type":
				data.Type = val
			case "description":
				data.Description = val
			case "year":
				if val != "" {
					y, err := strconv.Atoi(val)
					if err != nil {
						return fmt.Errorf("valor inválido para year: '%s' não é um número", val)
					}
					data.Year = y
				}
			case "featured":
				data.Featured = (strings.ToLower(val) == "true")
			}
		}
		frontmatter = data
	}

	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	targetDir := fmt.Sprintf("../entities/%ss", entityType)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.md", id))
	finalContent := fmt.Sprintf("---\n%s---\n%s\n", string(yamlData), body)

	if err := os.WriteFile(targetPath, []byte(finalContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Entidade criada: %s\n", targetPath)
	return nil
}

func Update(entityType string, id string, args []string) error {
	// Security: Sanitize inputs
	entityType = filepath.Base(entityType)
	id = filepath.Base(id)

	allowedTypes := map[string]bool{"action": true, "work": true, "agent": true}
	if !allowedTypes[entityType] {
		return fmt.Errorf("tipo de entidade inválido: %s", entityType)
	}

	targetDir := fmt.Sprintf("../entities/%ss", entityType)
	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.md", id))

	content, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("arquivo não encontrado: %s", targetPath)
	}

	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return fmt.Errorf("arquivo mal formatado")
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(parts[1], &data); err != nil {
		return err
	}

	// Update allows generic key setting
	if err := applyArgs(data, args); err != nil {
		return err
	}

	newYaml, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	finalContent := fmt.Sprintf("---\n%s---\n%s", newYaml, strings.TrimSpace(string(parts[2])))
	if err := os.WriteFile(targetPath, []byte(finalContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Entidade atualizada: %s\n", id)
	return nil
}
func applyArgs(data map[string]interface{}, args []string) error {
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) == 2 {
			key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			// Validate year
			if key == "year" || key == "active_since" {
				if val != "" {
					if _, err := strconv.Atoi(val); err != nil {
						return fmt.Errorf("valor inválido para %s: '%s' não é um número", key, val)
					}
				}
			}

			// Handle JSON for structured data (starts with [ or {)
			if (strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]")) ||
				(strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}")) {
				var jsonVal interface{}
				if err := json.Unmarshal([]byte(val), &jsonVal); err != nil {
					return fmt.Errorf("falha ao analisar %s (esperado JSON): %w", key, err)
				}
				data[key] = jsonVal
				continue
			}

			vLower := strings.ToLower(val)
			isBoolField := (key == "founded_by_me" || key == "featured")
			if existing, ok := data[key]; ok {
				if _, isBool := existing.(bool); isBool {
					isBoolField = true
				}
			}

			if isBoolField && (vLower == "true" || vLower == "false") {
				data[key] = (vLower == "true")
			} else {
				data[key] = val
			}
		}
	}
	return nil
}
