package syncer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncImagesYamlToBody_Sanitization(t *testing.T) {
	// Create a temporary workspace for testing
	dir := t.TempDir()
	worksDir := filepath.Join(dir, "works")
	os.MkdirAll(worksDir, 0755)

	filePath := filepath.Join(worksDir, "work-xss-test.md")
	content := `---
id: work-xss
attachment_1_url: 'foo")](http://malicious.com) .jpg'
attachment_1_type: image
attachment_1_label: 'Hacked [Label]'
---
Body text`
	os.WriteFile(filePath, []byte(content), 0644)

	err := SyncImages(dir, "yaml-to-body")
	if err != nil {
		t.Fatalf("SyncImages failed: %v", err)
	}

	res, _ := os.ReadFile(filePath)
	resStr := string(res)

	// Check that '[' and ']' are escaped in the label
	if !strings.Contains(resStr, `![Hacked \[Label\]]`) {
		t.Errorf("Expected markdown caption to be escaped. Got: \n%s", resStr)
	}

	// Check that ')' in the URL is replaced with '%29' and '"' with '%22'
	if !strings.Contains(resStr, `../../media/images/foo%22%29](http://malicious.com%29 .jpg`) {
		t.Errorf("Expected markdown URL to be escaped. Got: \n%s", resStr)
	}
}

func TestSyncImagesBodyToYaml_LegacyAndFlatMapping(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, "agents")
	os.MkdirAll(agentsDir, 0755)

	content2 := `---
id: agent-test2
attachments:
  - type: image
    url: legacy2.jpg
    caption: Legacy Caption
    role: legacy_role
attachment_5_type: image
attachment_5_url: flat6.jpg
attachment_5_label: Flat Label
attachment_5_caption: Ignore Me
---
![Body Image](../../media/images/body2.jpg)
![Legacy Image](../../media/images/legacy2.jpg)
![Flat Image](../../media/images/flat6.jpg)`
	filePath2 := filepath.Join(agentsDir, "agent-test2.md")
	os.WriteFile(filePath2, []byte(content2), 0644)

	SyncImages(dir, "body-to-yaml")

	res2, _ := os.ReadFile(filePath2)
	resStr2 := string(res2)

	// Validate legacy mapping
	if !strings.Contains(resStr2, "attachment_1_label: Legacy Caption") {
		t.Errorf("Expected legacy caption to be mapped to attachment_1_label. Output:\n%s", resStr2)
	}
	if !strings.Contains(resStr2, "attachment_1_category: legacy_role") {
		t.Errorf("Expected legacy role to be mapped to attachment_1_category. Output:\n%s", resStr2)
	}

	// Validate flat precedence mapping
	if !strings.Contains(resStr2, "attachment_2_label: Flat Label") {
		t.Errorf("Expected Flat Label to be maintained for attachment_2. Output:\n%s", resStr2)
	}
	if strings.Contains(resStr2, "Ignore Me") {
		t.Errorf("Did not expect 'Ignore Me' to persist since label should take precedence. Output:\n%s", resStr2)
	}

	// Validate body appending. The injected string from node AST renderer uses "body2.jpg" (DoubleQuotedStyle)
	if !strings.Contains(resStr2, "attachment_3_url: \"body2.jpg\"") {
		t.Errorf("Expected \"body2.jpg\" to be appended. Output:\n%s", resStr2)
	}
}
