package jsonrpc

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type allowedToolGate struct {
	allowed map[string]any
}

// NewAllowedToolGate creates a new replacer that only allows the specified tools.
func NewAllowedToolGate(allowed []string) *allowedToolGate {
	allowedMap := make(map[string]any)
	for _, item := range allowed {
		allowedMap[item] = nil
	}
	return &allowedToolGate{
		allowed: allowedMap,
	}
}

// Replace implements the replacer interface.
// It filters the response to `tools/list` method, allowing only the tools specified in the allowed list.
//
// see also https://modelcontextprotocol.io/specification/2025-06-18/basic
func (a *allowedToolGate) Replace(input string) (string, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return "", fmt.Errorf("failed to unmarshal input: %w", err)
	}

	result, ok := raw["result"]
	if !ok {
		return input, nil
	}
	result, err := a.replaceResult(result)
	if err != nil {
		return "", fmt.Errorf("failed to replace result: %w", err)
	}
	raw["result"] = result

	bs, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to marshal raw: %w", err)
	}
	return string(bs), nil
}

func (a *allowedToolGate) replaceResult(result json.RawMessage) (json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	tools, ok := m["tools"]
	if !ok {
		return result, nil
	}
	slog.Info("Tools before filtering", "tools", len(tools))

	tools, err := a.replaceTools(tools)
	if err != nil {
		return nil, fmt.Errorf("failed to replace tools: %w", err)
	}
	m["tools"] = tools

	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result map: %w", err)
	}
	return bs, nil
}

func (a *allowedToolGate) replaceTools(tools json.RawMessage) (json.RawMessage, error) {
	var toolList []map[string]json.RawMessage
	if err := json.Unmarshal(tools, &toolList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tools: %w", err)
	}

	var filtered []map[string]json.RawMessage
	for _, tool := range toolList {
		nameRaw, ok := tool["name"]
		if !ok {
			return nil, fmt.Errorf("tool has no name field")
		}
		var name string
		if err := json.Unmarshal(nameRaw, &name); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tool name: %w", err)
		}
		if _, ok := a.allowed[name]; ok {
			filtered = append(filtered, tool)
			slog.Info("Allowed tool", "name", name)
			continue
		}
		slog.Info("Filtered tool", "name", name)
	}

	bs, err := json.Marshal(filtered)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filtered tools: %w", err)
	}
	return bs, nil
}
