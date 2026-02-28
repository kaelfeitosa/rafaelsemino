package auditor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"acervo/internal/domain"
	"acervo/internal/utils"

	"gopkg.in/yaml.v3"
)

var (
	re     = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	reWiki = regexp.MustCompile(`!\[\[(.*?)\]\]`)
)

func Audit(entitiesDir string, imagesDir string, htmlPath string) error {
	agents := make(map[string]bool)
	works := make(map[string]bool)
	actions := make(map[string]bool)
	var allActions []domain.Action
	var allWorks []domain.Work
	referencedImages := make(map[string]bool)

	// Single pass walk to collect all entities and actions
	err := filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			return nil
		}

		relPath, err := filepath.Rel(entitiesDir, path)
		if err != nil {
			fmt.Printf("[WARNING] Falha ao obter caminho relativo para %s: %v\n", path, err)
			return nil
		}
		parentDir := filepath.Base(filepath.Dir(relPath))

		if parentDir == "agents" {
			var agent domain.Agent
			if err := yaml.Unmarshal(parts[1], &agent); err == nil {
				agents[agent.ID] = true
				for _, att := range agent.Attachments {
					if att.Type == "image" && att.Src != "" {
						referencedImages[filepath.Base(att.Src)] = true
					}
				}
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Agent em %s: %v\n", path, err)
			}
		} else if parentDir == "works" {
			var work domain.Work
			if err := yaml.Unmarshal(parts[1], &work); err == nil {
				works[work.ID] = true
				allWorks = append(allWorks, work)
				for _, att := range work.Attachments {
					if att.Type == "image" && att.Src != "" {
						referencedImages[filepath.Base(att.Src)] = true
					}
				}
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Work em %s: %v\n", path, err)
			}
		} else if parentDir == "actions" {
			var action domain.Action
			if err := yaml.Unmarshal(parts[1], &action); err == nil {
				actions[action.ID] = true
				allActions = append(allActions, action)
				for _, att := range action.Attachments {
					if att.Type == "image" && att.Src != "" {
						referencedImages[filepath.Base(att.Src)] = true
					}
				}
			} else {
				fmt.Printf("[WARNING] Falha ao analisar Action em %s: %v\n", path, err)
			}
		}

		// Also parse markdown body for referenced images
		if len(parts) == 3 {
			matches := re.FindAllSubmatch(parts[2], -1)
			for _, match := range matches {
				if len(match) > 1 {
					imgSrc := string(match[1])
					// Handle wiki links e.g., ![[image.jpg|300]] or standard URLs
					if strings.HasPrefix(imgSrc, "[") && strings.HasSuffix(imgSrc, "]") {
						imgSrc = strings.TrimPrefix(strings.TrimSuffix(imgSrc, "]"), "[")
						imgSrc = strings.Split(imgSrc, "|")[0]
					}
					referencedImages[filepath.Base(imgSrc)] = true
				}
			}

			// Handle obsidian wiki image syntax: ![[image.jpg]] or ![[image.jpg|300]]
			matchesWiki := reWiki.FindAllSubmatch(parts[2], -1)
			for _, match := range matchesWiki {
				if len(match) > 1 {
					imgSrc := string(match[1])
					imgSrc = strings.Split(imgSrc, "|")[0]
					referencedImages[filepath.Base(imgSrc)] = true
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	brokenLinks := 0
	// Validate relations using collected data
	for _, action := range allActions {
		pb := utils.CleanWikilink(action.PerformedBy)
		if pb != "" && !agents[pb] {
			fmt.Printf("[BROKEN LINK] Action %s aponta para Agent %s inexistente\n", action.ID, pb)
			brokenLinks++
		}

		wid := utils.CleanWikilink(action.WorkID)
		if wid != "" && !works[wid] {
			fmt.Printf("[BROKEN LINK] Action %s aponta para Work %s inexistente\n", action.ID, wid)
			brokenLinks++
		}
	}

	// 2. Entity Image Audit (Markdown Attachments)
	fmt.Println("--- ENTITY IMAGE AUDIT ---")
	missingImages := 0
	for img := range referencedImages {
		imgPath := filepath.Join(imagesDir, img)
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			fmt.Printf("[MISSING IMAGE] Imagem referenciada em entitidade não existe: %s\n", img)
			missingImages++
		}
	}

	// 3. Frontend Asset Audit (index.html References)
	fmt.Println("--- FRONTEND ASSET AUDIT ---")
	missingMasters := 0
	htmlContent, err := os.ReadFile(htmlPath)
	if err == nil {
		findRefs := func(content string) []string {
			var refs []string
			parts := strings.Split(content, "images/optimized/")
			for i := 1; i < len(parts); i++ {
				end := strings.IndexAny(parts[i], ` "'>`)
				if end != -1 {
					refs = append(refs, parts[i][:end])
				}
			}
			return refs
		}

		refs := findRefs(string(htmlContent))
		sourceMap, _ := buildHeuristicSourceMap(imagesDir)

		for _, ref := range refs {
			baseName := strings.TrimSuffix(ref, filepath.Ext(ref))
			normalizedBase := strings.ReplaceAll(baseName, "_", "-")
			if _, ok := sourceMap[normalizedBase]; !ok {
				fmt.Printf("[MISSING MASTER] Referência em index.html: %s (Nenhum master encontrado em %s)\n", ref, imagesDir)
				missingMasters++
			}
		}
	} else {
		fmt.Printf("[WARNING] Falha ao ler index.html para auditoria de assets: %v\n", err)
	}

	// 4. Dangling & Naming Audit
	fmt.Println("--- DANGLING & NAMING AUDIT ---")
	danglingImages := 0
	namingViolations := 0

	// Parse optional ignored images from environment variable (comma-separated)
	ignoredImages := make(map[string]bool)
	ignoredEnv := os.Getenv("ACERVO_IGNORE_IMAGES")
	if ignoredEnv != "" {
		for _, img := range strings.Split(ignoredEnv, ",") {
			trimmedImg := strings.TrimSpace(img)
			if trimmedImg != "" {
				ignoredImages[trimmedImg] = true
			}
		}
	} else {
		// Default to ignoring the protected test artifact
		ignoredImages["test-robust.jpeg"] = true
	}

	imageFiles, err := os.ReadDir(imagesDir)
	if err == nil {
		for _, f := range imageFiles {
			if f.IsDir() {
				continue
			}
			name := f.Name()

			// Skip specifically ignored artifacts
			if ignoredImages[name] {
				continue
			}

			// Check if referenced in entities
			if !referencedImages[name] {
				fmt.Printf("[DANGLING IMAGE] Imagem não referenciada por nenhuma entidade: %s\n", name)
				danglingImages++
			}

			// Check naming convention
			validPrefix := false
			expectedPrefixes := []string{"action-", "work-", "agent-"}
			for _, p := range expectedPrefixes {
				if strings.HasPrefix(name, p) {
					validPrefix = true
					break
				}
			}
			if !validPrefix {
				fmt.Printf("[NAMING ERROR] Imagem %s não possui prefixo válido (action-, work-, agent-)\n", name)
				namingViolations++
			}
		}
	}

	fmt.Printf("=== AUDIT REPORT ===\n")
	fmt.Printf("Broken Links: %d\n", brokenLinks)
	fmt.Printf("Missing Images (Entities): %d\n", missingImages)
	fmt.Printf("Missing Masters (Frontend): %d\n", missingMasters)
	fmt.Printf("Dangling Images: %d\n", danglingImages)
	fmt.Printf("Naming Violations: %d\n", namingViolations)

	if brokenLinks > 0 || missingImages > 0 || missingMasters > 0 {
		return fmt.Errorf("audit failed with reference errors")
	}
	return nil
}

// buildHeuristicSourceMap is a helper for the auditor
func buildHeuristicSourceMap(dir string) (map[string]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sourceMap := make(map[string]string)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fBase := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		fNorm := strings.ReplaceAll(fBase, "_", "-")
		sourceMap[fNorm] = f.Name()
	}
	return sourceMap, nil
}
