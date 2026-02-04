package dspy

import "testing"

func TestNewExample(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	input := Input{Question: "What is Go?"}
	output := Output{Answer: "Go is a programming language."}

	ex := NewExample(input, output)

	if ex.Input.Question != input.Question {
		t.Errorf("Expected input question '%s', got '%s'", input.Question, ex.Input.Question)
	}

	if ex.Output.Answer != output.Answer {
		t.Errorf("Expected output answer '%s', got '%s'", output.Answer, ex.Output.Answer)
	}
}

func TestExample(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	ex := Example[Input, Output]{
		Input:  Input{Text: "I love this!"},
		Output: Output{Label: "positive"},
	}

	if ex.Input.Text != "I love this!" {
		t.Errorf("Expected input text 'I love this!', got '%s'", ex.Input.Text)
	}

	if ex.Output.Label != "positive" {
		t.Errorf("Expected output label 'positive', got '%s'", ex.Output.Label)
	}
}
