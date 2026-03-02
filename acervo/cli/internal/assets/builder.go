package assets

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dsoprea/go-exif/v2"
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

	// Ensure sourceDir is absolute for reliable relative path calculation later
	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("resolving absolute source dir: %w", err)
	}

	// Create source map once
	sourceMap, err := buildSourceMap(absSourceDir)
	if err != nil {
		return fmt.Errorf("building source map: %w", err)
	}

	var buildErrors []string
	expectedFiles := make(map[string]bool)

	fmt.Printf("📂 Processing %d references using source map of %d files...\n", len(refs), len(sourceMap))

	for _, ref := range refs {
		filename := filepath.Base(ref)
		baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

		sourcePath := findSourceFile(sourceMap, baseName)
		if sourcePath == "" {
			continue // Not one of our master images
		}

		// Determine destination path preserving the relative directory structure
		relPath, err := filepath.Rel(absSourceDir, sourcePath)
		if err != nil {
			buildErrors = append(buildErrors, fmt.Sprintf("calculating relative path for %s: %v", sourcePath, err))
			continue
		}

		relPathWebp := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + ".webp"
		destPath := filepath.Join(outputDir, relPathWebp)
		absDest, err := filepath.Abs(destPath)
		if err != nil {
			buildErrors = append(buildErrors, fmt.Sprintf("resolving abs path for %s: %v", relPathWebp, err))
			continue
		}
		expectedFiles[absDest] = true

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(absDest), 0755); err != nil {
			buildErrors = append(buildErrors, fmt.Sprintf("creating dir for %s: %v", relPathWebp, err))
			continue
		}

		if !force && isUpToDate(sourcePath, absDest) {
			continue
		}

		fmt.Printf("🔨 Building: %s -> %s\n", filepath.Base(sourcePath), relPathWebp)
		if err := optimizeImage(cwebpPath, sourcePath, absDest); err != nil {
			fmt.Printf("❌ Error optimizing %s: %v\n", relPathWebp, err)
			buildErrors = append(buildErrors, fmt.Sprintf("optimizing %s: %v", relPathWebp, err))
		} else {
			fmt.Printf("✅ Optimized: %s\n", relPathWebp)
		}
	}

	if len(buildErrors) > 0 {
		return fmt.Errorf("encountered %d error(s) during asset build:\n- %s", len(buildErrors), strings.Join(buildErrors, "\n- "))
	}

	// Cleanup phase: remove files and then empty directories
	fmt.Println("🧹 Cleaning up unused optimized assets and folders...")

	// Pass 1: Remove unused files
	err = filepath.WalkDir(outputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}
		if !expectedFiles[absPath] {
			fmt.Printf("🗑️  Removing unused asset: %s\n", path)
			return os.Remove(path)
		}
		return nil
	})

	// Pass 2: Remove empty directories (bottom-up)
	filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && path != outputDir {
			entries, _ := os.ReadDir(path)
			if len(entries) == 0 {
				fmt.Printf("📂 Removing empty folder: %s\n", path)
				os.Remove(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("⚠️  Warning during cleanup: %v\n", err)
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
	sourceMap := make(map[string]string)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil // Only process JPEG and PNG masters
		}

		fBase := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		// Normalized key: replace underscores with hyphens
		fNorm := strings.ReplaceAll(fBase, "_", "-")

		if _, ok := sourceMap[fNorm]; ok {
			return nil // Skip duplicate basis, keep the first one found
		}

		absSource, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("resolving absolute path for %s: %v", d.Name(), err)
		}
		sourceMap[fNorm] = absSource
		return nil
	})

	if err != nil {
		return nil, err
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
	args := []string{
		"-q", WebPQuality,
	}

	// 1. Resolve orientation and rotate if necessary
	orientation := getOrientation(srcPath)
	if orientation > 1 {
		fmt.Printf("🔄 Rotating image (Orientation: %d): %s\n", orientation, filepath.Base(srcPath))
		// For simplicity and to avoid too many temporary files, we'll decode, rotate,
		// and encode to a temporary JPEG that cwebp will then consume.
		rotatedPath, err := rotateAndSaveTemp(srcPath, orientation)
		if err != nil {
			return fmt.Errorf("rotating image: %w", err)
		}
		defer os.Remove(rotatedPath)
		srcPath = rotatedPath
	}

	// 2. Decode config to check width of (possibly rotated) source
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	config, _, err := image.DecodeConfig(srcFile)
	if err != nil {
		return fmt.Errorf("decoding image config for %s: %w", srcPath, err)
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

func getOrientation(path string) int {
	rawExif, err := exif.SearchFileAndExtractExif(path)
	if err != nil {
		return 1
	}

	entries, err := exif.GetFlatExifData(rawExif)
	if err != nil {
		return 1
	}

	for _, entry := range entries {
		if entry.TagName == "Orientation" {
			if val, ok := entry.Value.([]uint16); ok && len(val) > 0 {
				return int(val[0])
			}
			if val, ok := entry.Value.(int); ok {
				return val
			}
		}
	}
	return 1
}

func rotateAndSaveTemp(srcPath string, orientation int) (string, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	var rotated image.Image
	switch orientation {
	case 3: // 180 degrees
		rotated = rotate180(img)
	case 6: // 90 degrees CW
		rotated = rotate90(img)
	case 8: // 270 degrees CW (90 CCW)
		rotated = rotate270(img)
	default:
		return srcPath, nil
	}

	tmpFile, err := os.CreateTemp("", "acervo-rotated-*.png")
	if err != nil {
		return "", err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close() // Close so we can write with png.Encode

	out, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, rotated); err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	return tmpPath, nil
}

func rotate90(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newImg.Set(bounds.Max.Y-y-1, x, img.At(x, y))
		}
	}
	return newImg
}

func rotate180(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newImg.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
		}
	}
	return newImg
}

func rotate270(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newImg.Set(y, bounds.Max.X-x-1, img.At(x, y))
		}
	}
	return newImg
}
