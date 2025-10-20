package jsonrpc

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_allowedToolGate_Replace(t *testing.T) {
	allowed := []string{"foo", "bar"}

	type testcase struct {
		name     string
		input    string
		expected string
		err      error
	}

	tests := []testcase{
		{
			name:     "filter tools",
			input:    `{"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"foo"},{"name":"bar"},{"name":"baz"}]}}`,
			expected: `{"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"foo"},{"name":"bar"}]}}`,
		},
		{
			name:     "preserve non tool response",
			input:    `{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05"}}`,
			expected: `{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAllowedToolGate(allowed)
			got, err := a.Replace(tt.input)
			if !errors.Is(err, tt.err) {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}

			must := func(s string) map[string]json.RawMessage {
				var m map[string]json.RawMessage
				if err := json.Unmarshal([]byte(s), &m); err != nil {
					t.Fatalf("failed to unmarshal string: %v", err)
				}
				return m
			}

			if diff := cmp.Diff(must(tt.expected), must(got)); diff != "" {
				t.Errorf("unexpected result (-expected +got):\n%s", diff)
			}
		})
	}
}
