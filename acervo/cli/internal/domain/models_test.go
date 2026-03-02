package domain

import (
	"testing"
	"gopkg.in/yaml.v3"
)

func TestAgentUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		yamlData     string
		expectedID   string
		expectedName string
		expectedAtts []Attachment
		expectError  bool
	}{
		{
			name: "valid agent with attachments",
			yamlData: `id: agent-1
name: John Doe
attachment_1_type: image
attachment_1_url: john.jpg
attachment_1_label: Profile
attachment_2_type: pdf
attachment_2_url: resume.pdf
`,
			expectedID:   "agent-1",
			expectedName: "John Doe",
			expectedAtts: []Attachment{
				{Type: "image", URL: "john.jpg", Label: "Profile"},
				{Type: "pdf", URL: "resume.pdf"},
			},
			expectError: false,
		},
		{
			name: "agent with no attachments",
			yamlData: `id: agent-2
name: Jane Doe`,
			expectedID:   "agent-2",
			expectedName: "Jane Doe",
			expectedAtts: nil,
			expectError:  false,
		},
		{
			name: "invalid yaml syntax",
			yamlData: `id: agent-3
name: Invalid
attachment_1_type: image
  invalid_indent: true`,
			expectedID:   "",
			expectedName: "",
			expectedAtts: nil,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var agent Agent
			err := yaml.Unmarshal([]byte(tc.yamlData), &agent)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if agent.ID != tc.expectedID || agent.Name != tc.expectedName {
				t.Errorf("Expected ID %s, Name %s, got ID %s, Name %s", tc.expectedID, tc.expectedName, agent.ID, agent.Name)
			}

			if len(agent.Attachments) != len(tc.expectedAtts) {
				t.Fatalf("Expected %d attachments, got %d", len(tc.expectedAtts), len(agent.Attachments))
			}

			for i, expected := range tc.expectedAtts {
				actual := agent.Attachments[i]
				if actual.Type != expected.Type || actual.URL != expected.URL || actual.Label != expected.Label || actual.Category != expected.Category {
					t.Errorf("Attachment %d mismatch. Expected %+v, got %+v", i, expected, actual)
				}
			}
		})
	}
}

func TestAgentMarshal(t *testing.T) {
	tests := []struct {
		name         string
		agent        Agent
		expectedAtts []Attachment
		expectError  bool
	}{
		{
			name: "agent with attachments",
			agent: Agent{
				ID:   "agent-2",
				Name: "Jane Doe",
				Attachments: []Attachment{
					{Type: "image", URL: "jane.jpg", Label: "Avatar", Category: "documentation"},
					{Type: "link", URL: "https://example.com", Label: "Website"},
				},
			},
			expectedAtts: []Attachment{
				{Type: "image", URL: "jane.jpg", Label: "Avatar", Category: "documentation"},
				{Type: "link", URL: "https://example.com", Label: "Website"},
			},
			expectError: false,
		},
		{
			name: "agent without attachments",
			agent: Agent{
				ID:   "agent-3",
				Name: "No Attachments",
			},
			expectedAtts: nil,
			expectError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := yaml.Marshal(tc.agent)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var agent Agent
			if err := yaml.Unmarshal(data, &agent); err != nil {
				t.Fatalf("Failed to unmarshal back to Agent: %v", err)
			}

			if agent.ID != tc.agent.ID || agent.Name != tc.agent.Name {
				t.Errorf("Marshaled simple fields incorrectly. Expected ID: %s Name: %s, got ID: %s Name: %s", tc.agent.ID, tc.agent.Name, agent.ID, agent.Name)
			}

			if len(agent.Attachments) != len(tc.expectedAtts) {
				t.Fatalf("Expected %d attachments, got %d. Yaml string:\n%s", len(tc.expectedAtts), len(agent.Attachments), string(data))
			}

			for i, expected := range tc.expectedAtts {
				actual := agent.Attachments[i]
				if actual.Type != expected.Type || actual.URL != expected.URL || actual.Label != expected.Label || actual.Category != expected.Category {
					t.Errorf("Attachment %d mismatch. Expected %+v, got %+v", i, expected, actual)
				}
			}
		})
	}
}
