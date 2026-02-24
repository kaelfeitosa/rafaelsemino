package ingester

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Ingester remains relatively dynamic because it needs to apply KV args to templates
// which might not map 1:1 to domain structs without heavy reflection logic.
// However, we ensure it respects the folder structure.

func Create(entityType string, slug string, args []string) error {
	targetDir := fmt.Sprintf("../entities/%ss", entityType)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

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

	applyArgs(data, args)

	newYaml, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	body := strings.TrimSpace(string(parts[2]))
	finalContent := fmt.Sprintf("---\n%s---\n%s", newYaml, body)

	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.md", id))
	if err := os.WriteFile(targetPath, []byte(finalContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Entidade criada: %s\n", targetPath)
	return nil
}

func Update(entityType string, id string, args []string) error {
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
			if val == "true" {
				data[key] = true
			} else if val == "false" {
				data[key] = false
			} else {
				data[key] = val
			}
		}
	}
}
