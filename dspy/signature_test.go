package dspy

import "testing"

func TestNewSignature(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	sig := NewSignature[Input, Output]("TestSignature", "Test description")

	if sig.Name != "TestSignature" {
		t.Errorf("Expected name 'TestSignature', got '%s'", sig.Name)
	}

	if sig.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got '%s'", sig.Description)
	}
}

func TestSignature(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	sig := Signature[Input, Output]{
		Name:        "Classification",
		Description: "Classify text",
	}

	if sig.Name != "Classification" {
		t.Errorf("Expected name 'Classification', got '%s'", sig.Name)
	}
}
