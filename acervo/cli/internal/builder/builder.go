package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"acervo/internal/domain"

	"gopkg.in/yaml.v3"
)

type TemplateData struct {
	TeatroFeatured         []domain.Work
	TeatroList             []domain.Work
	AudiovisualFeatured    []domain.Work
	AudiovisualList        []domain.Work
	PesquisaFeatured       []domain.Work
	PesquisaList           []domain.Work
	CulturaPopularFeatured []domain.Work
	CulturaPopularList     []domain.Work
	JogosDigitaisFeatured  []domain.Work
	JogosDigitaisList      []domain.Work
	ProducaoFeatured       []domain.Work
	ProducaoList           []domain.Work
	EnsinoFeatured         []domain.Work
	EnsinoList             []domain.Work
	OutrosFeatured         []domain.Work
	OutrosList             []domain.Work
}

func parseYear(y interface{}) int {
	switch v := y.(type) {
	case int:
		return v
	case string:
		var year int
		fmt.Sscanf(v, "%d", &year)
		return year
	default:
		return 0
	}
}

func isFeatured(w domain.Work) bool {
	for _, a := range w.Attachments {
		if a.Category == "documentation" || a.Category == "poster" {
			return true
		}
	}
	// Also explicitly check the boolean flag just in case
	if w.Featured {
		return true
	}
	return false
}

// BuildSite reads Works from Markdown, groups them by Medium, and executes the HTML template.
func BuildSite(entitiesDir, templatePath, outputPath string) error {
	fmt.Println("🚀 Inciando o SSG (Static Site Generator)...")

	var allWorks []domain.Work

	err := filepath.Walk(entitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			rel, _ := filepath.Rel(entitiesDir, path)
			parentDir := filepath.Base(filepath.Dir(rel))

			if parentDir == "works" {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				parts := bytes.SplitN(content, []byte("---"), 3)
				if len(parts) >= 3 {
					var work domain.Work
					if err := yaml.Unmarshal(parts[1], &work); err == nil {
						work.Description = strings.TrimSpace(string(parts[2]))
						allWorks = append(allWorks, work)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("erro lendo entities: %w", err)
	}

	sort.Slice(allWorks, func(i, j int) bool {
		yi := parseYear(allWorks[i].Year)
		yj := parseYear(allWorks[j].Year)
		return yi > yj
	})

	data := TemplateData{}
	for _, w := range allWorks {
		feat := isFeatured(w)
		switch w.Medium {
		case "teatro":
			if feat {
				data.TeatroFeatured = append(data.TeatroFeatured, w)
			} else {
				data.TeatroList = append(data.TeatroList, w)
			}
		case "audiovisual":
			if feat {
				data.AudiovisualFeatured = append(data.AudiovisualFeatured, w)
			} else {
				data.AudiovisualList = append(data.AudiovisualList, w)
			}
		case "pesquisa_academica", "pesquisa_artistica":
			if feat {
				data.PesquisaFeatured = append(data.PesquisaFeatured, w)
			} else {
				data.PesquisaList = append(data.PesquisaList, w)
			}
		case "cultura_popular":
			if feat {
				data.CulturaPopularFeatured = append(data.CulturaPopularFeatured, w)
			} else {
				data.CulturaPopularList = append(data.CulturaPopularList, w)
			}
		case "jogos_digitais":
			if feat {
				data.JogosDigitaisFeatured = append(data.JogosDigitaisFeatured, w)
			} else {
				data.JogosDigitaisList = append(data.JogosDigitaisList, w)
			}
		case "producao_cultural":
			if feat {
				data.ProducaoFeatured = append(data.ProducaoFeatured, w)
			} else {
				data.ProducaoList = append(data.ProducaoList, w)
			}
		case "ensino", "ensino_e_formacao", "formacao":
			if feat {
				data.EnsinoFeatured = append(data.EnsinoFeatured, w)
			} else {
				data.EnsinoList = append(data.EnsinoList, w)
			}
		default:
			if feat {
				data.OutrosFeatured = append(data.OutrosFeatured, w)
			} else {
				data.OutrosList = append(data.OutrosList, w)
			}
		}
	}

	fmt.Printf("✅ %d Works carregados na memória.\n", len(allWorks))

	// Define template functions
	funcMap := template.FuncMap{
		"mod": func(i, j int) int {
			if j == 0 {
				return 0
			}
			return i % j
		},
		"safeURL": func(s string) template.URL {
			// Resolve nested path by searching in the source images directory
			baseName := strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))
			// Normalized lookup: handle underscores vs hyphens
			normalizedBase := strings.ReplaceAll(baseName, "_", "-")

			// Heuristic: walk the source dir to find where this image lives
			sourceDir := "../media/images"
			var relFolder string
			filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					f := d.Name()
					fBase := strings.TrimSuffix(f, filepath.Ext(f))
					fNorm := strings.ReplaceAll(fBase, "_", "-")
					if fNorm == normalizedBase {
						rel, _ := filepath.Rel(sourceDir, path)
						relFolder = filepath.ToSlash(filepath.Dir(rel))
						return fmt.Errorf("found") // Stop walking
					}
				}
				return nil
			})

			path := ""
			if relFolder != "" && relFolder != "." {
				path = relFolder + "/" + baseName + ".webp"
			} else {
				path = baseName + ".webp"
			}
			return template.URL("images/optimized/" + path)
		},
		"getDocs": func(attachments []domain.Attachment) []domain.Attachment {
			var res []domain.Attachment
			for _, a := range attachments {
				if a.Category == "documentation" {
					res = append(res, a)
				}
			}
			return res
		},
		"getPosters": func(attachments []domain.Attachment) []domain.Attachment {
			var res []domain.Attachment
			for _, a := range attachments {
				if a.Category == "poster" {
					res = append(res, a)
				}
			}
			return res
		},
		"hasOccurrences": func(occ []domain.Occurrence) bool {
			return len(occ) > 0
		},
		"getWork": func(id string) domain.Work {
			for _, w := range allWorks {
				if w.ID == id {
					return w
				}
			}
			return domain.Work{}
		},
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("erro ao compilar o template: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro criando arquivo html final: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("erro injetando variáveis no template: %w", err)
	}

	fmt.Printf("🎉 Build concluído! Front-end gerado em: %s\n", outputPath)
	return nil
}
