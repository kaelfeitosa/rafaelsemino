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

	// New Schema compatible with the Work+Occurrence model
	if _, err := db.Exec(`
		CREATE TABLE entities(
			id TEXT PRIMARY KEY,
			type TEXT,
			title TEXT,
			path TEXT,
			featured INTEGER DEFAULT 0,
			json_data TEXT
		);
		CREATE TABLE occurrences(
			work_id TEXT,
			title TEXT,
			type TEXT,
			start_date TEXT,
			end_date TEXT,
			context TEXT,
			role TEXT,
			FOREIGN KEY(work_id) REFERENCES entities(id)
		);
		CREATE TABLE relations(
			src TEXT,
			rel TEXT,
			dst TEXT
		);
	`); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
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
			rel, err := filepath.Rel(entitiesDir, path)
			if err != nil {
				fmt.Printf("[WARNING] Falha ao obter caminho relativo para %s: %v\n", path, err)
				return nil
			}
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

				// Index occurrences
				for _, occ := range data.Occurrences {
					if _, err := tx.Exec("INSERT INTO occurrences VALUES(?,?,?,?,?,?,?)",
						id, occ.Title, occ.Type, occ.StartDate, occ.EndDate, occ.Context, occ.Role); err != nil {
						return fmt.Errorf("failed to insert occurrence for work %s: %w", id, err)
					}
				}
			}

			if id != "" {
				relPath, err := filepath.Rel(entitiesDir, path)
				if err != nil {
					return fmt.Errorf("falha ao obter caminho relativo para %s: %w", path, err)
				}
				_, err = tx.Exec("INSERT INTO entities VALUES(?,?,?,?,?,?)", id, typ, title, relPath, featured, string(jsonData))
				if err != nil {
					return fmt.Errorf("falha ao inserir entidade %s: %w", id, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return tx.Commit()
}
