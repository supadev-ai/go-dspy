package optimizer

import (
	"context"
	"testing"
	"time"

	"github.com/supadev-ai/go-dspy/dspy"
)

func TestNewBootstrapOptimizer(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	opt := NewBootstrapOptimizer[Input, Output]()
	if opt == nil {
		t.Fatal("Expected optimizer, got nil")
	}

	if opt.MaxIterations != 10 {
		t.Errorf("Expected max iterations 10, got %d", opt.MaxIterations)
	}

	if opt.MinScore != 0.8 {
		t.Errorf("Expected min score 0.8, got %f", opt.MinScore)
	}
}

func TestBootstrapOptimizer_WithMaxIterations(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	opt := NewBootstrapOptimizer[Input, Output]().WithMaxIterations(5)
	if opt.MaxIterations != 5 {
		t.Errorf("Expected max iterations 5, got %d", opt.MaxIterations)
	}
}

func TestBootstrapOptimizer_WithMinScore(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	opt := NewBootstrapOptimizer[Input, Output]().WithMinScore(0.9)
	if opt.MinScore != 0.9 {
		t.Errorf("Expected min score 0.9, got %f", opt.MinScore)
	}
}

func TestBootstrapOptimizer_WithTimeout(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	timeout := 10 * time.Minute
	opt := NewBootstrapOptimizer[Input, Output]().WithTimeout(timeout)
	if opt.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, opt.Timeout)
	}
}

func TestBootstrapOptimizer_Optimize_NoExamples(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	opt := NewBootstrapOptimizer[Input, Output]()
	module := &mockModule[Input, Output]{responses: make(map[string]Output)}
	examples := []dspy.Example[Input, Output]{}
	metric := ExactMatch[Input, Output]()

	ctx := context.Background()
	_, err := opt.Optimize(ctx, module, examples, metric)
	if err == nil {
		t.Error("Expected error for no examples, got nil")
	}
}

func TestBootstrapOptimizer_Optimize(t *testing.T) {
	type Input struct {
		Question string
	}
	type Output struct {
		Answer string
	}

	opt := NewBootstrapOptimizer[Input, Output]().
		WithMaxIterations(2).
		WithMinScore(0.0) // Low threshold to allow completion

	module := &mockModule[Input, Output]{
		responses: map[string]Output{
			"What is Go?": {Answer: "Go is a programming language."},
		},
	}

	examples := []dspy.Example[Input, Output]{
		dspy.NewExample(
			Input{Question: "What is Go?"},
			Output{Answer: "Go is a programming language."},
		),
	}

	metric := ExactMatch[Input, Output]()

	ctx := context.Background()
	optimized, err := opt.Optimize(ctx, module, examples, metric)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if optimized == nil {
		t.Fatal("Expected optimized module, got nil")
	}
}

func TestBootstrapOptimizer_Optimize_Timeout(t *testing.T) {
	type Input struct {
		Text string
	}
	type Output struct {
		Label string
	}

	opt := NewBootstrapOptimizer[Input, Output]().
		WithTimeout(1 * time.Nanosecond) // Very short timeout

	module := &mockModule[Input, Output]{responses: make(map[string]Output)}
	examples := []dspy.Example[Input, Output]{
		dspy.NewExample(
			Input{Text: "test"},
			Output{Label: "positive"},
		),
	}
	metric := ExactMatch[Input, Output]()

	ctx := context.Background()
	_, err := opt.Optimize(ctx, module, examples, metric)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
