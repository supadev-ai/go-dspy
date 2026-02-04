package optimizer

import (
	"testing"

	"github.com/supadev-ai/go-dspy/dspy"
)

func TestExactMatch(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	metric := ExactMatch[Input, Output]()

	predicted := Output{Label: "positive"}
	expected := Output{Label: "positive"}

	score := metric(predicted, expected)
	if score != 1.0 {
		t.Errorf("Expected score 1.0 for exact match, got %f", score)
	}

	predicted2 := Output{Label: "negative"}
	score2 := metric(predicted2, expected)
	if score2 != 0.0 {
		t.Errorf("Expected score 0.0 for mismatch, got %f", score2)
	}
}

func TestStringContains(t *testing.T) {
	type Input struct {
		Text string
	}

	metric := StringContains[Input]()

	predicted := "This is a positive review"
	expected := "positive"

	score := metric(predicted, expected)
	if score != 1.0 {
		t.Errorf("Expected score 1.0 for contains match, got %f", score)
	}

	predicted2 := "This is a negative review"
	score2 := metric(predicted2, expected)
	if score2 != 0.0 {
		t.Errorf("Expected score 0.0 for no match, got %f", score2)
	}
}

func TestStringContains_CaseInsensitive(t *testing.T) {
	type Input struct {
		Text string
	}

	metric := StringContains[Input]()

	predicted := "This is a POSITIVE review"
	expected := "positive"

	score := metric(predicted, expected)
	if score != 1.0 {
		t.Errorf("Expected score 1.0 for case-insensitive match, got %f", score)
	}
}

func TestStringContains_EmptyExpected(t *testing.T) {
	type Input struct {
		Text string
	}

	metric := StringContains[Input]()

	predicted := "any text"
	expected := ""

	score := metric(predicted, expected)
	if score != 1.0 {
		t.Errorf("Expected score 1.0 for empty expected string, got %f", score)
	}
}

func TestEvaluate(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	// Create a simple mock module
	mockModule := &mockModule[Input, Output]{
		responses: map[string]Output{
			"What is Go?": {Answer: "Go is a programming language."},
			"What is Python?": {Answer: "Python is a programming language."},
		},
	}

	examples := []dspy.Example[Input, Output]{
		dspy.NewExample(
			Input{Question: "What is Go?"},
			Output{Answer: "Go is a programming language."},
		),
		dspy.NewExample(
			Input{Question: "What is Python?"},
			Output{Answer: "Python is a programming language."},
		),
	}

	metric := ExactMatch[Input, Output]()

	score, err := Evaluate(mockModule, examples, metric)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if score != 1.0 {
		t.Errorf("Expected perfect score 1.0, got %f", score)
	}
}

func TestEvaluate_PartialMatch(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	mockModule := &mockModule[Input, Output]{
		responses: map[string]Output{
			"What is Go?": {Answer: "Go is a language."}, // Partial match
		},
	}

	examples := []dspy.Example[Input, Output]{
		dspy.NewExample(
			Input{Question: "What is Go?"},
			Output{Answer: "Go is a programming language."},
		),
	}

	metric := ExactMatch[Input, Output]()

	score, err := Evaluate(mockModule, examples, metric)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if score != 0.0 {
		t.Errorf("Expected score 0.0 for mismatch, got %f", score)
	}
}

func TestEvaluate_EmptyExamples(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	mockModule := &mockModule[Input, Output]{responses: make(map[string]Output)}
	examples := []dspy.Example[Input, Output]{}
	metric := ExactMatch[Input, Output]()

	score, err := Evaluate(mockModule, examples, metric)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if score != 0.0 {
		t.Errorf("Expected score 0.0 for empty examples, got %f", score)
	}
}

