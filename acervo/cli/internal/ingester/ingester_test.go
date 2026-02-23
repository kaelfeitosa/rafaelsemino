package ingester

import (
	"testing"
)

func TestApplyArgs(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		args     []string
		expected map[string]interface{}
	}{
		{
			name: "Boolean true",
			data: map[string]interface{}{"featured": false},
			args: []string{"featured=true"},
			expected: map[string]interface{}{"featured": true},
		},
		{
			name: "Boolean false",
			data: map[string]interface{}{"featured": true},
			args: []string{"featured=false"},
			expected: map[string]interface{}{"featured": false},
		},
		{
			name: "New Boolean true",
			data: map[string]interface{}{},
			args: []string{"featured=true"},
			expected: map[string]interface{}{"featured": true},
		},
		{
			name: "String value",
			data: map[string]interface{}{"title": "Old"},
			args: []string{"title=New"},
			expected: map[string]interface{}{"title": "New"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyArgs(tt.data, tt.args)
			for k, v := range tt.expected {
				if tt.data[k] != v {
					t.Errorf("expected %v for key %s, got %v", v, k, tt.data[k])
				}
			}
		})
	}
}
