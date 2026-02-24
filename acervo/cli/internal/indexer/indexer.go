package indexer

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

			var data EntityMap
			if err := yaml.Unmarshal(parts[1], &data); err != nil {
				return err
			}

			id, _ := data["id"].(string)

			// Determine Type
			typ := ""
			if strings.Contains(path, "/agents/") {
				typ = "agent"
			} else if strings.Contains(path, "/works/") {
				typ = "work"
			} else if strings.Contains(path, "/actions/") {
				typ = "action"
			}

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

			// Store as proper JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				// Fallback to empty JSON object if marshalling fails
				jsonData = []byte("{}")
			}

			_, err = db.Exec("INSERT INTO entities VALUES(?,?,?,?,?,?)", id, typ, title, absolutePath, featured, string(jsonData))
			if err != nil {
				return fmt.Errorf("falha ao inserir entidade %s: %w", id, err)
			}

			// Index Relations
			switch typ {
			case "action":
				if pb, ok := data["performed_by"].(string); ok && pb != "" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "performed_by", cleanWikilink(pb))
				}
				if wid, ok := data["work_id"].(string); ok && wid != "" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "work_id", cleanWikilink(wid))
				}
			case "work":
				// Any relations for Work? Maybe created_by if it exists in frontmatter
				if cb, ok := data["created_by"].(string); ok && cb != "" {
					db.Exec("INSERT INTO relations VALUES(?,?,?)", id, "created_by", cleanWikilink(cb))
				}
			case "agent":
				// founded_by_me, etc.
			}
		}
		return nil
	})
}
