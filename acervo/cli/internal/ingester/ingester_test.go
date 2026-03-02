package ingester

import (
	"reflect"
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
			name:     "Boolean true",
			data:     map[string]interface{}{"featured": false},
			args:     []string{"featured=true"},
			expected: map[string]interface{}{"featured": true},
		},
		{
			name:     "Boolean false",
			data:     map[string]interface{}{"featured": true},
			args:     []string{"featured=false"},
			expected: map[string]interface{}{"featured": false},
		},
		{
			name:     "New Boolean true",
			data:     map[string]interface{}{},
			args:     []string{"featured=true"},
			expected: map[string]interface{}{"featured": true},
		},
		{
			name:     "String value",
			data:     map[string]interface{}{"title": "Old"},
			args:     []string{"title=New"},
			expected: map[string]interface{}{"title": "New"},
		},
		{
			name:     "Existing string true",
			data:     map[string]interface{}{"description": "true"},
			args:     []string{"description=false"},
			expected: map[string]interface{}{"description": "false"},
		},
		{
			name:     "Array string valid JSON",
			data:     map[string]interface{}{"tags": []string{"a", "b"}},
			args:     []string{`tags=["c", "d"]`},
			expected: map[string]interface{}{"tags": []interface{}{"c", "d"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyArgs(tt.data, tt.args)
			for k, v := range tt.expected {
				if !reflect.DeepEqual(tt.data[k], v) {
					t.Errorf("expected %v for key %s, got %v", v, k, tt.data[k])
				}
			}
		})
	}
}
