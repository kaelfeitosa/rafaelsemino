package domain

import (
	"reflect"
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
	expected := []Attachment{
		{Type: "image", URL: "foo.jpg", Label: "Foo Label", Category: "documentation"},
		{Type: "video", URL: "bar.mp4"},
	}
	if !reflect.DeepEqual(w.Attachments, expected) {
		t.Errorf("Attachments do not match expected.\nGot:  %#v\nWant: %#v", w.Attachments, expected)
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
