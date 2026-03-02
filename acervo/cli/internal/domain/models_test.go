package domain

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestWorkMarshalUnmarshal(t *testing.T) {
	yamlStr := `id: work-1
title: Test Work
medium: teatro
attachment_1_type: image
attachment_1_url: foo.jpg
attachment_1_label: Foo Label
attachment_1_category: documentation
attachment_2_type: video
attachment_2_url: bar.mp4
`
	var w Work
	if err := yaml.Unmarshal([]byte(yamlStr), &w); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if w.ID != "work-1" {
		t.Errorf("Expected ID 'work-1', got '%s'", w.ID)
	}
	if len(w.Attachments) != 2 {
		t.Fatalf("Expected 2 attachments, got %d", len(w.Attachments))
	}
	if w.Attachments[0].URL != "foo.jpg" || w.Attachments[0].Type != "image" || w.Attachments[0].Label != "Foo Label" || w.Attachments[0].Category != "documentation" {
		t.Errorf("Unexpected first attachment: %+v", w.Attachments[0])
	}
	if w.Attachments[1].URL != "bar.mp4" || w.Attachments[1].Type != "video" {
		t.Errorf("Unexpected second attachment: %+v", w.Attachments[1])
	}

	out, err := yaml.Marshal(&w)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	outStr := string(out)
	if !strings.Contains(outStr, "attachment_1_url: foo.jpg") {
		t.Errorf("Marshal output missing attachment_1_url: \n%s", outStr)
	}
	if !strings.Contains(outStr, "attachment_2_type: video") {
		t.Errorf("Marshal output missing attachment_2_type: \n%s", outStr)
	}
}
