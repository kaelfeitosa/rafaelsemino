package syncer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestManual(t *testing.T) {
	tempDir := "/tmp/test_syncer"
	agentDir := filepath.Join(tempDir, "agents")
	os.MkdirAll(agentDir, 0755)

	content := `---
id: agent-1
name: Test Agent
attachment_1_type: image
attachment_1_url: test.jpg
attachment_1_label: Profile
attachment_2_type: pdf
attachment_2_url: test.pdf
---
Body text`

	filePath := filepath.Join(agentDir, "agent-1.md")
	os.WriteFile(filePath, []byte(content), 0644)

	err := SyncImages(tempDir, "yaml-to-body")
	if err != nil {
		fmt.Printf("SyncImages failed: %v", err)
	}

	newContent, _ := os.ReadFile(filePath)
	fmt.Println("RESULT:\n" + string(newContent))
}
