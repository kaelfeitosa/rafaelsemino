package indexer

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"acervo/internal/domain"

	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

func cleanWikilink(s string) string {
	s = strings.TrimPrefix(s, "[[")
	s = strings.TrimSuffix(s, "]]")
	return s
}

func isFeatured(v interface{}) bool {
	if val, ok := v.(bool); ok {
		return val
	}
	if val, ok := v.(string); ok {
		return strings.ToLower(val) == "true"
	}
	return false
}

func Reindex(entitiesDir, dbPath string) error {
	os.Remove(dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// New Schema compatible with the Action-centric model
	if _, err := db.Exec(`
		CREATE TABLE entities(
			id TEXT PRIMARY KEY,
			type TEXT,
			title TEXT,
			path TEXT,
			featured INTEGER DEFAULT 0,
			json_data TEXT
		);
		CREATE TABLE relations(
			src TEXT,
			rel TEXT,
			dst TEXT
		);
	`); err != nil {
		return err
	}

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
				return nil
			}

			var id, typ, title string
			var featured int
			var jsonData []byte

			// Robust entity type detection based on parent directory name
			rel, _ := filepath.Rel(entitiesDir, path)
			parentDir := filepath.Base(filepath.Dir(rel))

			if parentDir == "agents" {
				var data domain.Agent
				if err := yaml.Unmarshal(parts[1], &data); err != nil {
					return err
				}
				id = data.ID
				typ = "agent"
				title = data.Name
				if data.Featured {
					featured = 1
				}
				jsonData, err = json.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to marshal agent %s: %w", id, err)
				}

			} else if parentDir == "works" {
				var data domain.Work
				if err := yaml.Unmarshal(parts[1], &data); err != nil {
					return err
				}
				id = data.ID
				typ = "work"
				title = data.Title
				if data.Featured {
					featured = 1
				}
				jsonData, err = json.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to marshal work %s: %w", id, err)
				}

				if cb := cleanWikilink(data.CreatedBy); cb != "" {
					if _, err := db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "created_by", cb); err != nil {
						return fmt.Errorf("failed to insert relation created_by for %s: %w", id, err)
					}
				}

			} else if parentDir == "actions" {
				var data domain.Action
				if err := yaml.Unmarshal(parts[1], &data); err != nil {
					return err
				}
				id = data.ID
				typ = "action"
				title = data.Title
				if data.Featured {
					featured = 1
				}
				jsonData, err = json.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to marshal action %s: %w", id, err)
				}

				if pb := cleanWikilink(data.PerformedBy); pb != "" {
					if _, err := db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "performed_by", pb); err != nil {
						return fmt.Errorf("failed to insert relation performed_by for %s: %w", id, err)
					}
				}
				if wid := cleanWikilink(data.WorkID); wid != "" {
					if _, err := db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "work_id", wid); err != nil {
						return fmt.Errorf("failed to insert relation work_id for %s: %w", id, err)
					}
				}
			}

			if id != "" {
				absolutePath, _ := filepath.Abs(path)
				_, err = db.Exec("INSERT INTO entities VALUES(?,?,?,?,?,?)", id, typ, title, absolutePath, featured, string(jsonData))
				if err != nil {
					return fmt.Errorf("falha ao inserir entidade %s: %w", id, err)
				}
			}
		}
		return nil
	})
}
