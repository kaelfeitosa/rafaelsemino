package ingester

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func applyArgs(data map[string]interface{}, args []string) {
	for _, a := range args {
		parts := strings.SplitN(a, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// Check if the field is a slice in the template/existing data
			isSlice := false
			if v, ok := data[key]; ok {
				switch v.(type) {
				case []interface{}, []string:
					isSlice = true
				}
			}

			if isSlice {
				// Strip surrounding brackets if present
				if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
					val = strings.TrimSuffix(strings.TrimPrefix(val, "["), "]")
				}
				if val == "" {
					data[key] = []string{}
				} else {
					items := strings.Split(val, ",")
					var arr []string
					for i := range items {
						arr = append(arr, strings.TrimSpace(items[i]))
					}
					data[key] = arr
				}
			} else {
				data[key] = val
			}
		}
	}
}

func Create(entityType, slug string, args []string) error {
	templatePath := fmt.Sprintf("../templates/%s.md", entityType)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template não encontrado para %s", entityType)
	}

	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return fmt.Errorf("template mal formatado")
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(parts[1], &data); err != nil {
		return err
	}

	id := fmt.Sprintf("%s-%s", entityType, slug)
	data["id"] = id
	data["type"] = entityType

	applyArgs(data, args)

	newYaml, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	body := string(parts[2])

	// If it's a record, ensure we append the embed natively
	if entityType == "record" {
		if p, ok := data["path"].(string); ok && p != "" && !strings.Contains(p, "<arquivo>") {
			filename := filepath.Base(p)
			embed := fmt.Sprintf("![[%s]]", filename)
			if !strings.Contains(body, embed) {
				body = strings.TrimSpace(body) + "\n\n" + embed + "\n"
			}
		}
	}

	targetDir := fmt.Sprintf("../entities/%ss", entityType)
	os.MkdirAll(targetDir, 0755)

	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.md", id))

	finalContent := fmt.Sprintf("---\n%s---\n%s", newYaml, strings.TrimSpace(body))
	err = os.WriteFile(targetPath, []byte(finalContent), 0644)
	if err == nil {
		fmt.Printf("✅ Entidade criada: %s\n", id)
	}
	return err
}

func Update(entityType, id string, args []string) error {
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

	applyArgs(data, args)

	newYaml, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	finalContent := fmt.Sprintf("---\n%s---\n%s", newYaml, strings.TrimSpace(string(parts[2])))
	err = os.WriteFile(targetPath, []byte(finalContent), 0644)
	if err == nil {
		fmt.Printf("✅ Entidade atualizada: %s\n", id)
	}
	return err
}
