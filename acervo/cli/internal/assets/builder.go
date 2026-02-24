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
	WebPQuality     = "80"
	WebPMetadata    = "xmp"
	MaxImageWidth   = 1920
)

// BuildAssets scans HTML files for optimized image references and generates them
func BuildAssets(htmlPath, sourceDir, outputDir string) error {
	// Check for cwebp availability
	cwebpPath, err := exec.LookPath("cwebp")
	if err != nil {
		fmt.Println("❌ Error: 'cwebp' tool not found in PATH.")
		fmt.Println("   Please install WebP tools:")
		fmt.Println("   - macOS: brew install webp")
		fmt.Println("   - Linux: sudo apt-get install webp")
		fmt.Println("   - Windows: Download from https://developers.google.com/speed/webp/docs/precompiled")
		return fmt.Errorf("cwebp dependency missing")
	}

	refs, err := scanHTMLForImages(htmlPath)
	if err != nil {
		return fmt.Errorf("scanning HTML: %w", err)
	}

	if len(refs) == 0 {
		fmt.Println("No optimized images found in HTML.")
		return nil
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

	for _, ref := range refs {
		if !strings.Contains(ref, "images/optimized/") {
			continue
		}

		filename := filepath.Base(ref)
		baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

		sourcePath := findSourceFile(sourceMap, baseName)
		if sourcePath == "" {
			fmt.Printf("⚠️  Source not found for: %s\n", baseName)
			continue
		}

		destPath := filepath.Join(outputDir, filename)
		if isUpToDate(sourcePath, destPath) {
			continue
		}

		fmt.Printf("🔨 Building: %s -> %s\n", filepath.Base(sourcePath), filename)
		if err := optimizeImage(cwebpPath, sourcePath, destPath); err != nil {
			errorMsg := fmt.Sprintf("optimizing %s: %v", filename, err)
			fmt.Printf("❌ Error %s\n", errorMsg)
			buildErrors = append(buildErrors, errorMsg)
		} else {
			fmt.Printf("✅ Optimized: %s\n", filename)
		}
	}

	if len(buildErrors) > 0 {
		return fmt.Errorf("encountered %d error(s) during asset build:\n- %s", len(buildErrors), strings.Join(buildErrors, "\n- "))
	}

	return nil
}

func scanHTMLForImages(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Updated regex to handle whitespace around '='
	re := regexp.MustCompile(`src\s*=\s*["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	var refs []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			ref := match[1]
			if !seen[ref] {
				refs = append(refs, ref)
				seen[ref] = true
			}
		}
	}
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

		fBase := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		// Normalized key: replace underscores with hyphens
		fNorm := strings.ReplaceAll(fBase, "_", "-")

		sourceMap[fNorm] = filepath.Join(dir, f.Name())
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

func optimizeImage(cwebpCmd, srcPath, destPath string) error {
	// 1. Decode config to check width
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	// Defer closing immediately to ensure resource cleanup
	defer srcFile.Close()

	// We only need to decode config, not the whole image
	config, _, err := image.DecodeConfig(srcFile)
	if err != nil {
		return fmt.Errorf("decoding image config: %w", err)
	}

	args := []string{
		"-q", WebPQuality,
		"-metadata", WebPMetadata,
		"-quiet",
	}

	if config.Width > MaxImageWidth {
		// Calculate height to maintain aspect ratio (0 = auto)
		args = append(args, "-resize", strconv.Itoa(MaxImageWidth), "0")
	}

	args = append(args, srcPath, "-o", destPath)

	// Execute cwebp
	cmd := exec.Command(cwebpCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cwebp failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
