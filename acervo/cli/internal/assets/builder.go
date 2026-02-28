package assets

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	WebPQuality   = "80"
	WebPMetadata  = "xmp"
	MaxImageWidth = 1920
)

// BuildAssets scans HTML files for optimized image references and generates them
func BuildAssets(htmlPath, sourceDir, outputDir string, force bool) error {
	refs, err := scanHTMLForImages(htmlPath)
	if err != nil {
		return fmt.Errorf("scanning HTML: %w", err)
	}

	if len(refs) == 0 {
		fmt.Printf("No optimized images found in HTML file: %s.\n", htmlPath)
		return nil
	}

	// Check for cwebp availability
	cwebpPath := "cwebp"
	// Look for cwebp.exe in the current working directory (where we run the command)
	if _, err := os.Stat("cwebp.exe"); err == nil {
		if absPath, err := filepath.Abs("cwebp.exe"); err == nil {
			cwebpPath = absPath
		} else {
			// If we found cwebp.exe but couldn't get its absolute path, log a warning and fall back to just "cwebp.exe"
			fmt.Printf("⚠️  Warning: Found 'cwebp.exe' but could not resolve absolute path: %v. Using 'cwebp.exe' directly.\n", err)
			cwebpPath = "cwebp.exe"
		}
	} else {
		// cwebp.exe not found in current directory, try looking in PATH
		var lookPathErr error
		cwebpPath, lookPathErr = exec.LookPath("cwebp")
		if lookPathErr != nil {
			fmt.Println("❌ Error: 'cwebp' tool not found in PATH or current directory.")
			fmt.Println("   Please install WebP tools:")
			fmt.Println("   - macOS: brew install webp")
			fmt.Println("   - Linux: sudo apt-get install webp")
			fmt.Println("   - Windows: Download from https://developers.google.com/speed/webp/docs/precompiled")
			return fmt.Errorf("cwebp dependency missing: %w", lookPathErr)
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Create source map once
	sourceMap, err := buildSourceMap(sourceDir)
	if err != nil {
		return fmt.Errorf("building source map: %w", err)
	}

	var buildErrors []string

	fmt.Printf("📂 Processing %d references using source map of %d files...\n", len(refs), len(sourceMap))

	for _, ref := range refs {
		filename := filepath.Base(ref)
		baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

		sourcePath := findSourceFile(sourceMap, baseName)
		if sourcePath == "" {
			continue // Not one of our master images
		}

		// Determine destination path. We want to preserve the name but use .webp
		relPath := baseName + ".webp"
		destPath := filepath.Join(outputDir, relPath)
		absDest, err := filepath.Abs(destPath)
		if err != nil {
			buildErrors = append(buildErrors, fmt.Sprintf("resolving abs path for %s: %v", relPath, err))
			continue
		}

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(absDest), 0755); err != nil {
			buildErrors = append(buildErrors, fmt.Sprintf("creating dir for %s: %v", relPath, err))
			continue
		}

		if !force && isUpToDate(sourcePath, absDest) {
			continue
		}

		fmt.Printf("🔨 Building: %s -> %s\n", filepath.Base(sourcePath), relPath)
		if err := optimizeImage(cwebpPath, sourcePath, absDest); err != nil {
			fmt.Printf("❌ Error optimizing %s: %v\n", relPath, err)
			buildErrors = append(buildErrors, fmt.Sprintf("optimizing %s: %v", relPath, err))
		} else {
			fmt.Printf("✅ Optimized: %s\n", relPath)
		}
	}

	if len(buildErrors) > 0 {
		return fmt.Errorf("encountered %d error(s) during asset build:\n- %s", len(buildErrors), strings.Join(buildErrors, "\n- "))
	}

	return nil
}

func scanHTMLForImages(path string) ([]string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("⚠️  Warning: could not resolve absolute path for %s: %v\n", path, err)
		absPath = path // Fallback to original
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("reading HTML at %s: %w", absPath, err)
	}

	// Refined regex: non-greedy, handles src and url() with various quoting
	re := regexp.MustCompile(`(?:src\s*=\s*["'](?P<src>[^"']+)["']|url\s*\(\s*["']?(?P<url>[^"'\)]+)["']?\s*\))`)
	matches := re.FindAllSubmatch(content, -1)

	var refs []string
	seen := make(map[string]bool)

	for _, match := range matches {
		for i, name := range re.SubexpNames() {
			if (name == "src" || name == "url") && i < len(match) && match[i] != nil {
				ref := strings.TrimSpace(string(match[i]))
				if ref != "" && !seen[ref] {
					refs = append(refs, ref)
					seen[ref] = true
				}
			}
		}
	}
	fmt.Fprintf(os.Stderr, "🔍 DEBUG: All unique references found: %v\n", refs)
	return refs, nil
}

// buildSourceMap scans the source directory and creates a map of normalized filenames to full paths
func buildSourceMap(dir string) (map[string]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	sourceMap := make(map[string]string)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue // Only process JPEG and PNG masters
		}

		fBase := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		// Normalized key: replace underscores with hyphens
		fNorm := strings.ReplaceAll(fBase, "_", "-")

		if existing, ok := sourceMap[fNorm]; ok {
			return nil, fmt.Errorf("filename collision: %s and %s both map to %s", filepath.Base(existing), f.Name(), fNorm)
		}

		absSource, err := filepath.Abs(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("resolving absolute path for %s: %v", f.Name(), err)
		}
		sourceMap[fNorm] = absSource
	}
	return sourceMap, nil
}

func findSourceFile(sourceMap map[string]string, baseName string) string {
	// Normalized lookup
	normalizedBase := strings.ReplaceAll(baseName, "_", "-")
	if path, ok := sourceMap[normalizedBase]; ok {
		return path
	}
	return ""
}

func isUpToDate(src, dest string) bool {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return false
	}
	destInfo, err := os.Stat(dest)
	if err != nil {
		return false
	}
	return destInfo.ModTime().After(srcInfo.ModTime())
}

func optimizeImage(cwebpCmd, srcPath, destPath string) (err error) {
	// 1. Decode config to check width
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	// Defer closing with combined error handling
	defer func() {
		if cerr := srcFile.Close(); cerr != nil {
			if err == nil {
				err = fmt.Errorf("closing source file: %w", cerr)
			} else {
				// Don't overwrite the primary error, but wrap it
				err = fmt.Errorf("%w; additionally failed to close source file: %v", err, cerr)
			}
		}
	}()

	// We only need to decode config, not the whole image
	config, format, err := image.DecodeConfig(srcFile)
	if err != nil {
		return fmt.Errorf("decoding image config for %s: %w", srcPath, err)
	}
	_ = format // For now we don't use it, but good for debug

	args := []string{
		"-q", WebPQuality,
	}

	if config.Width > MaxImageWidth {
		args = append(args, "-resize", strconv.Itoa(MaxImageWidth), "0")
	}

	args = append(args, srcPath, "-o", destPath)

	// Execute cwebp
	fmt.Printf("🔨 Running: %s %v\n", cwebpCmd, args)
	cmd := exec.Command(cwebpCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cwebp failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
