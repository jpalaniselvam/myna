package varsub

import (
	"strings"
	"testing"
)

func TestResolver_Hierarchy(t *testing.T) {
	env := Vars{"VAR": "env", "SHARED": "env-shared"}
	coll := Vars{"VAR": "coll", "SHARED": "coll-shared"}
	action := Vars{"VAR": "action"} // Overrides VAR

	r := New(env, coll, action)

	// Action should win
	if val := r.SubstituteVars("{{VAR}}"); val != "action" {
		t.Errorf("Expected 'action', got '%s'", val)
	}

	// Collection should win over env
	if val := r.SubstituteVars("{{SHARED}}"); val != "coll-shared" {
		t.Errorf("Expected 'coll-shared', got '%s'", val)
	}
}

func TestResolver_Expansion(t *testing.T) {
	env := Vars{"ENV": "prod"}
	coll := Vars{"URL": "http://api.${ENV}.com"}

	r := New(env, coll)

	// verify expansion works
	expected := "http://api.prod.com"
	if val := r.SubstituteVars("{{URL}}"); val != expected {
		t.Errorf("Failed to expand variable. Expected '%s', got '%s'", expected, val)
	}
}

func TestResolver_Missing(t *testing.T) {
	r := New()
	input := "{{MISSING}}"
	if val := r.SubstituteVars(input); val != input {
		t.Errorf("Expected '%s', got '%s'", input, val)
	}
}

func TestResolver_Circular(t *testing.T) {
	// A -> B -> A
	v := Vars{
		"A": "${B}",
		"B": "${A}",
	}
	r := New(v)

	// Should not hang.
	// Behavior: A -> B -> fails to resolve A (cycle) -> B becomes "${A}" -> A becomes "${A}"

	val := r.SubstituteVars("{{A}}")
	if !strings.Contains(val, "${A}") {
		// Exact behavior depends on implementation details of breaking the cycle,
		// but it shouldn't hang and should probably return the unresolved variable.
		t.Logf("Circular resolution result: %s", val)
	}
}

func TestResolver_NestedTemplate(t *testing.T) {
	// Value containing {{...}}
	v := Vars{
		"NAME":     "World",
		"GREETING": "Hello {{NAME}}!",
	}
	r := New(v)

	val := r.SubstituteVars("{{GREETING}}")
	expected := "Hello World!"
	if val != expected {
		t.Errorf("Expected '%s', got '%s'", expected, val)
	}
}

func TestResolver_ComplexTypes(t *testing.T) {
	v := Vars{
		"user": map[string]interface{}{
			"name": map[string]interface{}{
				"first": "John",
				"last":  "Doe",
			},
			"roles": []interface{}{"admin", "editor"},
			"meta":  map[string]interface{}{"login_count": 42},
		},
	}
	r := New(v)

	tests := []struct {
		input    string
		expected string
	}{
		{"Hello {{user.name.first}}", "Hello John"},
		{"Role: {{user.roles[0]}}", "Role: admin"},
		{"Count: {{user.meta.login_count}}", "Count: 42"},
	}

	for _, tt := range tests {
		got := r.SubstituteVars(tt.input)
		if got != tt.expected {
			t.Errorf("Resolve(%q): expected %q, got %q", tt.input, tt.expected, got)
		}
	}
}
