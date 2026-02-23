package indexer

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

type EntityMap map[string]interface{}

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

	if _, err := db.Exec(`
		CREATE TABLE entities(
			id TEXT PRIMARY KEY,
			type TEXT,
			title TEXT,
			path TEXT,
			featured INTEGER DEFAULT 0
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
				return nil // Skip empty files or wrong formats
			}

			var data EntityMap
			if err := yaml.Unmarshal(parts[1], &data); err != nil {
				return err
			}

			id, _ := data["id"].(string)
			typ, _ := data["type"].(string)

			title := ""
			if val, ok := data["title"].(string); ok {
				title = val
			} else if val, ok := data["name"].(string); ok {
				title = val
			}

			absolutePath, _ := filepath.Abs(path)

			featured := 0
			if isFeatured(data["featured"]) {
				featured = 1
			}

			_, err = db.Exec("INSERT INTO entities VALUES(?,?,?,?,?)", id, typ, title, absolutePath, featured)
			if err != nil {
				return fmt.Errorf("falha ao inserir entidade %s: %w", id, err)
			}

			switch typ {
			case "event":
				if orgs, ok := data["organizers"]; ok {
					if orgList, okList := orgs.([]interface{}); okList {
						for _, item := range orgList {
							if target, isStr := item.(string); isStr && target != "" {
								db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "organizer", cleanWikilink(target))
							}
						}
					}
				}
			case "work":
				if c, ok := data["created_by"].(string); ok && c != "" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "created_by", cleanWikilink(c))
				}
			case "participation":
				if a, ok := data["agent"].(string); ok && a != "" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "agent", cleanWikilink(a))
				}
				if e, ok := data["event"].(string); ok && e != "" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "event", cleanWikilink(e))
				}
				if w, ok := data["work"].(string); ok && w != "" && w != "null" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "work", cleanWikilink(w))
				}
			case "record":
				if rel, ok := data["related_to"]; ok {
					if relStr, okStr := rel.(string); okStr && relStr != "" {
						db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "related_to", cleanWikilink(relStr))
					} else if relList, okList := rel.([]interface{}); okList {
						for _, item := range relList {
							if target, isStr := item.(string); isStr && target != "" {
								db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "related_to", cleanWikilink(target))
							}
						}
					}
				}
			}
		}
		return nil
	})
}
