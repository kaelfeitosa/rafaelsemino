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

// IngesterData holds all potential fields for template rendering
type IngesterData struct {
	ID          string
	Title       string
	Name        string
	Kind        string
	Description string
	PerformedBy string
	MyRole      string
	WorkID      string
	Context     struct {
		Label    string
		Kind     string
		Location string
		Year     string // Keep as string for template, can be parsed later or left empty
	}
	DateStart   string
	DateEnd     string
	Type        string
	Year        string
	FoundedByMe bool
	ActiveSince string
	ContentBody string
}

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
			if len(kv) == 2 {
				key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
				switch key {
				case "title":
					data.Title = val
				case "kind":
					data.Kind = val
				case "performed_by":
					data.PerformedBy = toWikilink(val)
				case "my_role":
					data.MyRole = val
				case "work_id":
					data.WorkID = toWikilink(val)
				case "context_label":
					data.Context.Label = val
				case "context_kind":
					data.Context.Kind = val
				case "context_location":
					data.Context.Location = val
				case "context_year":
					if val != "" {
						y, err := strconv.Atoi(val)
						if err != nil {
							return fmt.Errorf("valor inválido para context_year: '%s' não é um número", val)
						}
						data.Context.Year = y
					}
				case "date_start":
					data.DateStart = val
				case "date_end":
					data.DateEnd = val
				case "description":
					data.Description = val
				case "featured":
					data.Featured = (strings.ToLower(val) == "true")
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
			if len(kv) == 2 {
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
		}
		frontmatter = data
	case "work":
		data := &domain.Work{
			ID: id,
		}
		for _, arg := range args {
			kv := strings.SplitN(arg, "=", 2)
			if len(kv) == 2 {
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
	applyArgs(data, args)

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
func applyArgs(data map[string]interface{}, args []string) {
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) == 2 {
			key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

			// Validate year/context_year
			if key == "year" || key == "context_year" || key == "active_since" {
				if val != "" {
					if _, err := strconv.Atoi(val); err != nil {
						fmt.Printf("[WARNING] Valor inválido para %s: '%s' não é um número. Ignorando.\n", key, val)
						continue
					}
				}
			}

			// Handle JSON for structured data (starts with [ or {)
			if (strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]")) ||
				(strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}")) {
				var jsonVal interface{}
				if err := json.Unmarshal([]byte(val), &jsonVal); err == nil {
					data[key] = jsonVal
					continue
				}
				// Fallback to custom parsing if JSON fails (e.g. simple list [a,b])
				val = strings.TrimPrefix(val, "[")
				val = strings.TrimSuffix(val, "]")
				parts := strings.Split(val, ",")
				var slice []string
				for _, p := range parts {
					slice = append(slice, strings.TrimSpace(p))
				}
				data[key] = slice
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
			} else if strings.Contains(val, ",") {
				// Handle legacy structured slices (backwards compatibility)
				parts := strings.Split(val, ",")
				var slice []interface{}
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if strings.Contains(p, "|") {
						objParts := strings.Split(p, "|")
						switch key {
						case "collaborators":
							obj := make(map[string]string)
							obj["name"] = objParts[0]
							if len(objParts) > 1 {
								obj["role"] = objParts[1]
							}
							slice = append(slice, obj)
						case "attachments":
							obj := make(map[string]string)
							obj["type"] = objParts[0]
							if len(objParts) > 1 {
								obj["role"] = objParts[1]
							}
							if len(objParts) > 2 {
								obj["src"] = objParts[2]
							}
							if len(objParts) > 3 {
								obj["caption"] = objParts[3]
							}
							slice = append(slice, obj)
						default:
							slice = append(slice, p)
						}
					} else {
						slice = append(slice, p)
					}
				}
				data[key] = slice
			} else {
				data[key] = val
			}
		}
	}
}
