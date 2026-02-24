package ingester

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

	targetDir := fmt.Sprintf("../entities/%ss", entityType)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	templatePath := fmt.Sprintf("../templates/%s.md", entityType)
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template não encontrado para %s", entityType)
	}

	// Parse template
	tmpl, err := template.New(entityType).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("invalid template syntax: %w", err)
	}

	// Prepare data
	id := fmt.Sprintf("%s-%s", entityType, slug)
	data := IngesterData{
		ID:          id,
		ContentBody: "Detalhes aqui.", // Default body
	}

	// Simple args parsing to populate struct fields
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) == 2 {
			key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			switch key {
			case "title":
				data.Title = val
			case "name":
				data.Name = val
			case "kind":
				data.Kind = val
			case "description":
				data.Description = val
			case "performed_by":
				data.PerformedBy = val
			case "my_role":
				data.MyRole = val
			case "work_id":
				data.WorkID = val
			case "context_label":
				data.Context.Label = val
			case "context_kind":
				data.Context.Kind = val
			case "context_location":
				data.Context.Location = val
			case "context_year":
				data.Context.Year = val
			case "date_start":
				data.DateStart = val
			case "date_end":
				data.DateEnd = val
			case "type":
				data.Type = val
			case "year":
				data.Year = val
			case "founded_by_me":
				data.FoundedByMe = (strings.ToLower(val) == "true")
			case "active_since":
				data.ActiveSince = val
			}
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.md", id))
	if err := os.WriteFile(targetPath, buf.Bytes(), 0644); err != nil {
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

			// Handle Array format [a, b, c]
			if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
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
			// Smart bool parsing: only if existing value is bool OR it's a known boolean field
			isBoolField := (key == "founded_by_me" || key == "featured")
			if existing, ok := data[key]; ok {
				if _, isBool := existing.(bool); isBool {
					isBoolField = true
				}
			}

			if isBoolField && (vLower == "true" || vLower == "false") {
				data[key] = (vLower == "true")
			} else if strings.Contains(val, ",") {
				// Handle structured slices
				parts := strings.Split(val, ",")
				var slice []interface{}
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if strings.Contains(p, "|") {
						objParts := strings.Split(p, "|")
						obj := make(map[string]string)
						if key == "collaborators" {
							obj["name"] = objParts[0]
							if len(objParts) > 1 {
								obj["role"] = objParts[1]
							}
						} else if key == "attachments" {
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
						} else {
							obj["value"] = objParts[0]
							if len(objParts) > 1 {
								obj["role"] = objParts[1]
							}
						}
						slice = append(slice, obj)
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
