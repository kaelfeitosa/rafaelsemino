package assets

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

	for _, ref := range refs {
		if !strings.Contains(ref, "images/optimized/") {
			continue
		}

		filename := filepath.Base(ref)
		baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

		sourcePath := findSourceFile(sourceDir, baseName)
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
			fmt.Printf("❌ Error optimizing %s: %v\n", filename, err)
		} else {
			fmt.Printf("✅ Optimized: %s\n", filename)
		}
	}

	return nil
}

func scanHTMLForImages(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`src=["']([^"']+)["']`)
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

func findSourceFile(dir, baseName string) string {
	exts := []string{".jpg", ".jpeg", ".png"}

	// 1. Try exact match
	for _, ext := range exts {
		path := filepath.Join(dir, baseName+ext)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 2. Try normalized match (hyphens vs underscores)
	files, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	normalizedTarget := strings.ReplaceAll(baseName, "_", "-")

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fBase := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		fNorm := strings.ReplaceAll(fBase, "_", "-")

		if fNorm == normalizedTarget {
			return filepath.Join(dir, f.Name())
		}
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
	// We only need to decode config, not the whole image
	// But image.DecodeConfig needs registered formats
	config, _, err := image.DecodeConfig(srcFile)
	srcFile.Close()

	if err != nil {
		return fmt.Errorf("decoding image config: %w", err)
	}

	maxWidth := 1920
	args := []string{
		"-q", "80",          // Quality 80
		"-metadata", "xmp",  // Copy XMP metadata
		"-quiet",            // Reduce noise
	}

	if config.Width > maxWidth {
		// Calculate height to maintain aspect ratio (0 = auto)
		args = append(args, "-resize", strconv.Itoa(maxWidth), "0")
	}

	args = append(args, srcPath, "-o", destPath)

	// Execute cwebp
	cmd := exec.Command(cwebpCmd, args...)
	// Capture stderr for debugging if it fails
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cwebp failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// Side-effect imports to register formats for DecodeConfig
func init() {
	// Register formats is handled by blank imports in import block
	// but image/jpeg and image/png must be there
	_ = jpeg.Decode
	_ = png.Decode
}
