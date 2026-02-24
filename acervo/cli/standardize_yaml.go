package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Minimal Action struct for ordering
type OrderedAction struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Kind        string `yaml:"kind"`
	PerformedBy string `yaml:"performed_by"`
	MyRole      string `yaml:"my_role"`
	WorkID      string `yaml:"work_id"`
	Context     struct {
		Label    string `yaml:"label"`
		Kind     string `yaml:"kind"`
		Location string `yaml:"location"`
		Year     int    `yaml:"year,omitempty"`
	} `yaml:"context"`
	DateStart     string        `yaml:"date_start"`
	DateEnd       string        `yaml:"date_end,omitempty"`
	Description   string        `yaml:"description"`
	Collaborators []interface{} `yaml:"collaborators,omitempty"`
	Attachments   []interface{} `yaml:"attachments,omitempty"`
	Featured      bool          `yaml:"featured,omitempty"`
}

func main() {
	types := []string{"actions", "agents", "works"}
	for _, t := range types {
		entitiesDir := filepath.Join("../entities", t)
		files, err := ioutil.ReadDir(entitiesDir)
		if err != nil {
			fmt.Printf("Error reading dir %s: %v\n", entitiesDir, err)
			continue
		}

		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			path := filepath.Join(entitiesDir, f.Name())
			content, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				continue
			}

			parts := bytes.SplitN(content, []byte("---"), 3)
			if len(parts) < 3 {
				continue
			}

			var data map[string]interface{}
			if err := yaml.Unmarshal(parts[1], &data); err != nil {
				fmt.Printf("Error unmarshaling %s: %v\n", path, err)
				continue
			}

			// Marshal back with default ordering
			newFrontmatter, err := yaml.Marshal(&data)
			if err != nil {
				fmt.Printf("Error marshaling %s: %v\n", path, err)
				continue
			}

			finalContent := fmt.Sprintf("---\n%s---\n%s", string(newFrontmatter), string(parts[2]))
			if err := ioutil.WriteFile(path, []byte(finalContent), 0644); err != nil {
				fmt.Printf("Error writing file %s: %v\n", path, err)
			} else {
				fmt.Printf("✅ Standardized: %s/%s\n", t, f.Name())
			}
		}
	}
}
