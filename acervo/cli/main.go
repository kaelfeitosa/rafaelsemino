package main

import (
	"fmt"
	"os"

	"bytes"
	"path/filepath"
	"strings"

	"acervo/internal/domain"
	"gopkg.in/yaml.v3"


	"acervo/internal/assets"
	"acervo/internal/auditor"
	"acervo/internal/builder"
	"acervo/internal/indexer"
	"acervo/internal/ingester"
	"acervo/internal/metadata"
	"acervo/internal/reposter"
	"acervo/internal/setup"
	"acervo/internal/validator"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "acervo",
		Short: "Acervo CLI for managing the editorial portfolio database (Action-centric)",
	}

	var validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates all markdown entities (Actions, Works, Agents)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := validator.ValidateEntities("../entities"); err != nil {
				fmt.Println("❌ Validation failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Validação OK")
		},
	}

	var verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Executes strict Syntax Validation and Graph Integrity Auditing",
		Run: func(cmd *cobra.Command, args []string) {
			if err := validator.ValidateEntities("../entities"); err != nil {
				fmt.Println("❌ Validation failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Validação Sintática OK")

			if err := auditor.Audit("../entities", "../media/images", "../../frontend/index.html"); err != nil {
				fmt.Println("❌ Audit failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Auditoria de Grafo OK")
		},
	}

	var reindexCmd = &cobra.Command{
		Use:   "reindex",
		Short: "Rebuilds the SQLite database from markdown entries (auto-verifies first)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Executando pré-verificação obrigatória...")
			if err := validator.ValidateEntities("../entities"); err != nil {
				fmt.Println("❌ Validation blocked reindex:", err)
				os.Exit(1)
			}
			/*
				if err := auditor.Audit("../entities", "../media/images", "../../frontend/index.html"); err != nil {
					fmt.Println("❌ Audit blocked reindex:", err)
					os.Exit(1)
				}
			*/

			if err := indexer.Reindex("../entities", "../db.sqlite"); err != nil {
				fmt.Println("❌ Reindex failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Index gerado em ../db.sqlite")
		},
	}

	var ingestCreateCmd = &cobra.Command{
		Use:   "create [entityType] [slug] [key=value...]",
		Short: "Creates a new markdown entity strictly from templates (action, work, agent)",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			entityType := args[0]
			slug := args[1]
			kvArgs := args[2:]
			if err := ingester.Create(entityType, slug, kvArgs); err != nil {
				fmt.Println("❌ Create failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Entity created successfully.")
		},
	}

	var ingestUpdateCmd = &cobra.Command{
		Use:   "update [entityType] [id] [key=value...]",
		Short: "Updates an existing markdown entity's YAML properties",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			entityType := args[0]
			id := args[1]
			kvArgs := args[2:]
			if err := ingester.Update(entityType, id, kvArgs); err != nil {
				fmt.Println("❌ Update failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Entity updated successfully.")
		},
	}

	var ingestCmd = &cobra.Command{
		Use:   "ingest",
		Short: "Ingester CMS Headless engine",
	}
	ingestCmd.AddCommand(ingestCreateCmd, ingestUpdateCmd)

	var setFocusCmd = &cobra.Command{
		Use:   "set-focus [imagePath] [x] [y]",
		Short: "Sets the XMP focus point for an image",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			x, errX := metadata.GetFloat64(args[1])
			y, errY := metadata.GetFloat64(args[2])
			if errX != nil || errY != nil {
				fmt.Println("❌ Coordenadas inválidas. Use números entre 0.0 e 1.0")
				os.Exit(1)
			}
			if err := metadata.SetFocus(path, x, y); err != nil {
				fmt.Println("❌ Erro ao definir foco:", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Foco definido em (%.4f, %.4f) para %s\n", x, y, path)
		},
	}

	var hooksCmd = &cobra.Command{
		Use:   "hooks",
		Short: "Installs Git pre-commit hooks",
		Run: func(cmd *cobra.Command, args []string) {
			if err := setup.InstallHooks(); err != nil {
				fmt.Println("❌ Erro ao instalar hooks:", err)
				os.Exit(1)
			}
		},
	}

	var (
		htmlPath  string
		sourceDir string
		outputDir string
		force     bool
	)

	var buildAssetsCmd = &cobra.Command{
		Use:   "build-assets",
		Short: "Scans HTML for images and generates optimized WebP assets using cwebp",
		Long: `Scans the specified HTML file for 'images/optimized/' references and generates
WebP assets from source master images.
Default paths assume the command is run from 'acervo/cli'.
Use absolute paths or adjust flags if running from elsewhere.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚀 Starting asset optimization...")
			if err := assets.BuildAssets(htmlPath, sourceDir, outputDir, force); err != nil {
				fmt.Println("❌ Error building assets:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Assets optimized successfully.")
		},
	}
	buildAssetsCmd.Flags().StringVar(&htmlPath, "html", "../../frontend/index.html", "Path to HTML file to scan (relative to execution dir or absolute)")
	buildAssetsCmd.Flags().StringVar(&sourceDir, "source", "../media/images", "Directory containing source master images (relative to execution dir or absolute)")
	buildAssetsCmd.Flags().StringVar(&outputDir, "output", "../../frontend/images/optimized", "Directory to output optimized WebP assets (relative to execution dir or absolute)")
	buildAssetsCmd.Flags().BoolVar(&force, "force", false, "Force regeneration of all assets even if they are up to date")

	var (
		repo        string
		pullNumber  string
		reviewID    string
		mentionUser string
		token       string
	)

	var repostReviewCmd = &cobra.Command{
		Use:   "repost-review",
		Short: "Reposts a GitHub PR review with a specific user mention",
		Run: func(cmd *cobra.Command, args []string) {
			if token == "" {
				token = os.Getenv("USER_PAT")
			}
			if token == "" {
				fmt.Fprintln(os.Stderr, "Error: USER_PAT environment variable or --token must be set")
				os.Exit(1)
			}

			rep := reposter.NewReviewReposter(token, repo, pullNumber)
			if err := rep.RepostReview(reviewID, mentionUser); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Review reposted successfully.")
		},
	}
	repostReviewCmd.Flags().StringVar(&repo, "repo", "", "GitHub repository (owner/repo)")
	repostReviewCmd.Flags().StringVar(&pullNumber, "pull-number", "", "Pull request number")
	repostReviewCmd.Flags().StringVar(&reviewID, "review-id", "", "ID of the review to repost")
	repostReviewCmd.Flags().StringVar(&mentionUser, "mention-user", "@jules", "User to mention in the reposted review")
	repostReviewCmd.Flags().StringVar(&token, "token", "", "GitHub PAT (optional, uses USER_PAT env var if not set)")
	repostReviewCmd.MarkFlagRequired("repo")
	repostReviewCmd.MarkFlagRequired("pull-number")
	repostReviewCmd.MarkFlagRequired("review-id")

	var buildSiteCmd = &cobra.Command{
		Use:   "build-site",
		Short: "Generates the static frontend (index.html) from markdown Works",
		Run: func(cmd *cobra.Command, args []string) {
			if err := builder.BuildSite("../entities", "../../frontend/index.tmpl", "../../frontend/index.html"); err != nil {
				fmt.Println("❌ Build failed:", err)
				os.Exit(1)
			}
		},
	}


	var migrateAttachmentsCmd = &cobra.Command{
		Use:   "migrate-attachments",
		Short: "Migrates nested YAML attachments to flattened attachment_X_Y keys",
		Run: func(cmd *cobra.Command, args []string) {
			// Actually, just running syncer body-to-yaml is practically doing this
			// since we added migration logic into syncer.go itself for `attachments` map parsing.
			// Let's implement a clean pass using yaml.v3 Unmarshal/Marshal which our models now support.
			fmt.Println("🚀 Starting attachment migration...")
			err := filepath.Walk("../entities", func(path string, info os.FileInfo, err error) error {
				if err != nil { return err }
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
					content, err := os.ReadFile(path)
					if err != nil { return err }
					parts := bytes.SplitN(content, []byte("---"), 3)
					if len(parts) >= 3 {
						rel, _ := filepath.Rel("../entities", path)
						parentDir := filepath.Base(filepath.Dir(rel))

						var frontmatter map[string]interface{}
						if err := yaml.Unmarshal(parts[1], &frontmatter); err != nil { return err }

						// If the file has 'attachments' or already needs formatting
						if parentDir == "agents" {
							var agent domain.Agent
							if err := yaml.Unmarshal(parts[1], &agent); err != nil {
								return fmt.Errorf("failed to unmarshal agent %s: %w", path, err)
							}
							newYaml, err := yaml.Marshal(agent)
							if err != nil {
								return fmt.Errorf("failed to marshal agent %s: %w", path, err)
							}
							newContent := string(parts[0]) + "---\n" + string(newYaml) + "---" + string(parts[2])
							if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
								return fmt.Errorf("failed to write %s: %w", path, err)
							}
						} else if parentDir == "works" {
							var work domain.Work
							if err := yaml.Unmarshal(parts[1], &work); err != nil {
								return fmt.Errorf("failed to unmarshal work %s: %w", path, err)
							}
							newYaml, err := yaml.Marshal(work)
							if err != nil {
								return fmt.Errorf("failed to marshal work %s: %w", path, err)
							}
							newContent := string(parts[0]) + "---\n" + string(newYaml) + "---" + string(parts[2])
							if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
								return fmt.Errorf("failed to write %s: %w", path, err)
							}
						}
					}
				}
				return nil
			})
			if err != nil {
				fmt.Println("❌ Migration failed:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Migration completed successfully.")
		},
	}

	rootCmd.AddCommand(migrateAttachmentsCmd, validateCmd, reindexCmd, verifyCmd, ingestCmd, hooksCmd, setFocusCmd, buildAssetsCmd, repostReviewCmd, buildSiteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
