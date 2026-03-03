package domain

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestWorkMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		yamlStr  string
		expected []Attachment
	}{
		{
			name: "Flattened Format",
			yamlStr: `id: work-1
title: Test Work
medium: teatro
attachment_1_type: image
attachment_1_url: foo.jpg
attachment_1_label: Foo Label
attachment_1_category: documentation
attachment_2_type: video
attachment_2_url: bar.mp4
`,
			expected: []Attachment{
				{Type: "image", URL: "foo.jpg", Label: "Foo Label", Category: "documentation"},
				{Type: "video", URL: "bar.mp4"},
			},
		},
		{
			name: "Legacy Array Format",
			yamlStr: `id: work-2
title: Test Work 2
medium: teatro
attachments:
  - type: image
    url: baz.jpg
    caption: Baz Label
    role: documentation
`,
			expected: []Attachment{
				{Type: "image", URL: "baz.jpg", Label: "Baz Label", Category: "documentation"},
			},
		},
		{
			name: "Mixed Format",
			yamlStr: `id: work-3
title: Test Work 3
medium: teatro
attachments:
  - type: image
    url: legacy.jpg
    caption: Legacy Label
attachment_2_type: video
attachment_2_url: flat.mp4
`,
			expected: []Attachment{
				{Type: "image", URL: "legacy.jpg", Label: "Legacy Label", Category: ""},
				{Type: "video", URL: "flat.mp4"},
			},
		},
		{
			name: "Fallback Priority (Label over Caption)",
			yamlStr: `id: work-4
title: Test Work 4
medium: teatro
attachment_1_type: image
attachment_1_url: img.jpg
attachment_1_label: Correct Label
attachment_1_caption: Wrong Label
attachment_1_category: Correct Category
attachment_1_role: Wrong Category
`,
			expected: []Attachment{
				{Type: "image", URL: "img.jpg", Label: "Correct Label", Category: "Correct Category"},
			},
		},
		{
			name: "Fallback (Caption when Label is empty)",
			yamlStr: `id: work-5
title: Test Work 5
medium: teatro
attachment_1_type: image
attachment_1_url: img.jpg
attachment_1_label: ""
attachment_1_caption: Fallback Label
attachment_1_category: ""
attachment_1_role: Fallback Category
`,
			expected: []Attachment{
				{Type: "image", URL: "img.jpg", Label: "Fallback Label", Category: "Fallback Category"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w Work
			if err := yaml.Unmarshal([]byte(tt.yamlStr), &w); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if !reflect.DeepEqual(w.Attachments, tt.expected) {
				t.Errorf("Attachments do not match expected.\nGot:  %#v\nWant: %#v", w.Attachments, tt.expected)
			}

			out, err := yaml.Marshal(&w)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var w2 Work
			if err := yaml.Unmarshal(out, &w2); err != nil {
				t.Fatalf("Failed to unmarshal marshaled output: %v", err)
			}

			// Validate roundtrip guarantees the same attachments structure
			if !reflect.DeepEqual(w.Attachments, w2.Attachments) {
				t.Errorf("Roundtrip failed. Got: %#v\nWant: %#v", w2.Attachments, w.Attachments)
			}
		})
	}
}
