package varsub

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Vars map[string]interface{}

type Resolver struct {
	scopes []Vars
}

func New(scopes ...Vars) *Resolver {
	return &Resolver{scopes: scopes}
}

func (r *Resolver) AddScope(scope Vars) {
	r.scopes = append(r.scopes, scope)
}

// lookup finds the value for a key in the scopes (highest priority first).
// It supports dot notation for nested maps (e.g., "user.name") and array indexing (e.g., "users[0]").
func (r *Resolver) lookup(key string) (interface{}, bool) {
	// Identify the root variable name to find the correct scope
	root := getRootKey(key)

	for i := len(r.scopes) - 1; i >= 0; i-- {
		// Check if the root variable exists in this scope
		if _, ok := r.scopes[i][root]; ok {
			// Found the scope containing the root variable.
			// Now traverse the full path within this scope.
			return getValue(r.scopes[i], key)
		}
	}
	return nil, false
}

var templatePattern = regexp.MustCompile(`{{\s*([\w.\[\]-]+)\s*}}`)

// Resolve replaces variables in the input string.
// It handles {{VAR}} syntax in the input, and ${VAR} syntax recursively within values.
func (r *Resolver) Resolve(input []byte) []byte {
	result := r.SubstituteVars(string(input))
	return []byte(result)
}

func (r *Resolver) SubstituteVars(input string) string {
	// First, resolve {{...}} templates in the input string
	resolved := templatePattern.ReplaceAllStringFunc(string(input), func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		val, ok := r.recursiveLookup(key, make(map[string]bool))
		if ok {
			return fmt.Sprintf("%v", val)
		}
		// Log warning
		fmt.Fprintf(os.Stderr, "WARN: Variable '%s' not found in resolution hierarchy\n", key)
		return match // leave unresolved
	})

	// Apply usage of ${...} on the result, as per docs "String interpolation within TOML values"
	// effectively treating the input as a "value" too if it contains shell-style vars.
	return os.Expand(resolved, func(key string) string {
		if v, ok := r.recursiveLookup(key, make(map[string]bool)); ok {
			return fmt.Sprintf("%v", v)
		}
		fmt.Fprintf(os.Stderr, "WARN: Variable '%s' not found in resolution hierarchy\n", key)
		return ""
	})
}

// recursiveLookup resolves a key, handling recursion and expansion.
func (r *Resolver) recursiveLookup(key string, visited map[string]bool) (interface{}, bool) {
	if visited[key] {
		// Detect cycle
		return nil, false
	}
	visited[key] = true
	defer delete(visited, key) // Backtrack: allow other paths to visit this key

	rawVal, ok := r.lookup(key)
	if !ok {
		return nil, false
	}

	// Helper to warn
	warn := func(k string) {
		fmt.Fprintf(os.Stderr, "WARN: Variable '%s' not found in resolution hierarchy\n", k)
	}

	// Only perform string interpolation if the value is a string
	strVal, isString := rawVal.(string)
	if !isString {
		return rawVal, true
	}

	// 1. Expand ${VAR} within the value
	expanded := os.Expand(strVal, func(subKey string) string {
		val, found := r.recursiveLookup(subKey, visited)
		if found {
			return fmt.Sprintf("%v", val)
		}
		warn(subKey)
		return "${" + subKey + "}"
	})

	// 2. Expand {{VAR}} within the value
	expanded = templatePattern.ReplaceAllStringFunc(expanded, func(match string) string {
		subKey := strings.TrimSpace(match[2 : len(match)-2])
		val, found := r.recursiveLookup(subKey, visited)
		if found {
			return fmt.Sprintf("%v", val)
		}
		warn(subKey)
		return match
	})

	return expanded, true
}

// Helper functions for map traversal

func getRootKey(key string) string {
	idxDot := strings.Index(key, ".")
	idxBracket := strings.Index(key, "[")

	if idxDot == -1 && idxBracket == -1 {
		return key
	}
	if idxDot == -1 {
		return key[:idxBracket]
	}
	if idxBracket == -1 {
		return key[:idxDot]
	}
	// Both exist, take min
	if idxDot < idxBracket {
		return key[:idxDot]
	}
	return key[:idxBracket]
}

func getValue(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		// Check for array index: key[index]
		if idx := strings.Index(part, "["); idx != -1 && strings.HasSuffix(part, "]") {
			key := part[:idx]
			indexStr := part[idx+1 : len(part)-1]
			var index int
			if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
				return nil, false
			}

			// 1. Get the list from the map using the key
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil, false
			}

			listVal, exists := m[key]
			if !exists {
				return nil, false
			}

			// 2. Access the index
			list, ok := listVal.([]interface{})
			if !ok {
				return nil, false
			}

			if index < 0 || index >= len(list) {
				return nil, false
			}
			current = list[index]

		} else {
			// standard map access
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil, false
			}
			val, exists := m[part]
			if !exists {
				return nil, false
			}
			current = val
		}
	}
	return current, true
}
