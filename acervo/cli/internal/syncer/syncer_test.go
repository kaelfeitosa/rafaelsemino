package syncer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncImages(t *testing.T) {
	tests := []struct {
		name           string
		mode           string
		entityType     string
		fileName       string
		initialContent string
		expectContains []string
		expectNotContains []string
	}{
		{
			name:       "Yaml to Body",
			mode:       "yaml-to-body",
			entityType: "agents",
			fileName:   "agent-1.md",
			initialContent: `---
id: agent-1
name: Test Agent
attachment_1_type: image
attachment_1_url: test.jpg
attachment_1_label: Profile
attachment_2_type: pdf
attachment_2_url: test.pdf
---
Body text`,
			expectContains: []string{
				"![Profile](../../media/images/test.jpg)",
			},
			expectNotContains: []string{
				"test.pdf",
			},
		},
		{
			name:       "Body to Yaml",
			mode:       "body-to-yaml",
			entityType: "works",
			fileName:   "work-1.md",
			initialContent: `---
id: work-1
title: Test Work
attachment_1_type: pdf
attachment_1_url: doc.pdf
attachment_2_type: image
attachment_2_url: removed.jpg
---
Body text
![Image 1](../../media/images/added.jpg)
`,
			expectContains: []string{
				"attachment_1_type: pdf",
				"attachment_2_url: added.jpg",
			},
			expectNotContains: []string{
				"removed.jpg",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			entityDir := filepath.Join(tempDir, tc.entityType)
			os.MkdirAll(entityDir, 0755)

			filePath := filepath.Join(entityDir, tc.fileName)
			os.WriteFile(filePath, []byte(tc.initialContent), 0644)

			err := SyncImages(tempDir, tc.mode)
			if err != nil {
				t.Fatalf("SyncImages failed: %v", err)
			}

			newContent, _ := os.ReadFile(filePath)
			contentStr := string(newContent)

			var targetStr string
			parts := strings.Split(contentStr, "---")

			if tc.mode == "yaml-to-body" {
				if len(parts) >= 3 {
					targetStr = parts[2]
				} else {
					targetStr = contentStr
				}
			} else {
				if len(parts) >= 2 {
					targetStr = parts[1]
				} else {
					targetStr = contentStr
				}
			}

			for _, expected := range tc.expectContains {
				if !strings.Contains(targetStr, expected) {
					t.Errorf("Expected content to contain '%s', but it didn't.\nTarget Section:\n%s", expected, targetStr)
				}
			}

			for _, notExpected := range tc.expectNotContains {
				if strings.Contains(targetStr, notExpected) {
					t.Errorf("Expected content to NOT contain '%s', but it did.\nTarget Section:\n%s", notExpected, targetStr)
				}
			}
		})
	}
}
