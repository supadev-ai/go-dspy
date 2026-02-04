# go-dspy

[![Tests](https://github.com/supadev-ai/go-dspy/workflows/Tests/badge.svg)](https://github.com/supadev-ai/go-dspy/actions)
[![CI](https://github.com/supadev-ai/go-dspy/workflows/CI/badge.svg)](https://github.com/supadev-ai/go-dspy/actions)

**go-dspy** is a Go-native implementation of DSPy (Declarative Self-Improving Programs), enabling structured LLM pipelines, prompt optimization, and self-improving reasoning graphs with static typing, performance, and production readiness.

> **Developed by [SupaDev](https://github.com/supadev-ai)** - An AI services company specializing in production-ready AI infrastructure and tooling.

## What is DSPy?

DSPy = Declarative Self-Improving Programs

DSPy is a framework for building LLM applications that automatically optimizes prompts based on examples and metrics, rather than requiring manual prompt engineering. Instead of hard-coding prompts, you declare the structure of your pipeline and let DSPy optimize it.

## Why Go?

**go-dspy** brings DSPy's powerful concepts to Go with several key advantages:

- **Static Typing**: Type-safe interfaces prevent runtime errors
- **Concurrency-First**: Leverage Go's goroutines for parallel evaluation and optimization
- **Production-Grade**: Built-in timeouts, retries, and observability
- **No Runtime Magic**: Explicit, predictable behavior without hidden complexity
- **Performance**: Native Go performance for high-throughput applications

## Key Differentiators vs Python DSPy

| Feature | Python DSPy | go-dspy |
|---------|-------------|---------|
| Type Safety | Runtime checks | Compile-time guarantees |
| Concurrency | Limited | Native goroutines |
| Production Features | Basic | Timeouts, retries, tracing |
| Runtime Magic | Yes | No - explicit interfaces |

## Installation

```bash
go get github.com/supadev-ai/go-dspy
```

## Quick Start

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/supadev-ai/go-dspy/dspy"
    "github.com/supadev-ai/go-dspy/llm"
)

type QAInput struct {
    Question string
}

type QAOutput struct {
    Answer string
}

func main() {
    // Create an LLM client
    client := llm.NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))
    
    // Define a signature
    sig := dspy.NewSignature[QAInput, QAOutput](
        "QuestionAnswering",
        "Answer questions based on your knowledge.",
    )
    
    // Create a predictor
    qa := dspy.NewPredictor(sig, client)
    
    // Use it
    ctx := context.Background()
    out, _ := qa.Forward(ctx, QAInput{
        Question: "What is DSPy?",
    })
    
    fmt.Println(out.Answer)
}
```

### With Optimization

```go
// Create training examples
examples := []dspy.Example[QAInput, QAOutput]{
    dspy.NewExample(
        QAInput{Question: "What is Go?"},
        QAOutput{Answer: "Go is a programming language."},
    ),
    // ... more examples
}

// Create optimizer
opt := optimizer.NewBootstrapOptimizer[QAInput, QAOutput]().
    WithMaxIterations(10).
    WithMinScore(0.8)

// Create metric
metric := optimizer.ExactMatch[QAInput, QAOutput]()

// Optimize
optimized, _ := opt.Optimize(ctx, qa, examples, metric)
```

## Core Concepts

### Signature

A `Signature` defines the input/output contract for a module:

```go
type Signature[I any, O any] struct {
    Name        string
    Description string
}
```

### Module

A `Module` is a composable unit that transforms input to output:

```go
type Module[I any, O any] interface {
    Forward(ctx context.Context, input I) (O, error)
}
```

### Predictor

A `Predictor` is an LLM-backed module:

```go
predictor := dspy.NewPredictor[QAInput, QAOutput](signature, llmClient)
output, err := predictor.Forward(ctx, input)
```

## LLM Providers

go-dspy supports multiple LLM providers:

- **OpenAI**: `llm.NewOpenAIClient(apiKey)`
- **Anthropic**: `llm.NewAnthropicClient(apiKey)`
- **Mock**: `llm.NewMockClient()` (for testing)

## Optimization

### Bootstrap Optimizer

The bootstrap optimizer runs your module on examples and tracks performance:

```go
opt := optimizer.NewBootstrapOptimizer[I, O]().
    WithMaxIterations(10).
    WithMinScore(0.8).
    WithTimeout(5 * time.Minute)
```

### Metrics

Built-in metrics for evaluation:

- `ExactMatch`: Exact equality comparison
- `StringContains`: Substring matching (case-insensitive)

## Examples

See the `cmd/examples/` directory for complete examples:

- **basic**: Simple question-answering
- **chain_of_thought**: Step-by-step reasoning
- **optimizer**: Using the bootstrap optimizer

Run an example:

```bash
cd cmd/examples/basic
go run main.go
```

## Testing

The project includes comprehensive test coverage. Run tests with:

**Run all tests:**
```bash
go test ./...
```

**Run tests with verbose output:**
```bash
go test ./... -v
```

**Run tests for a specific package:**
```bash
go test ./dspy -v
go test ./llm -v
go test ./optimizer -v
```

**Run a specific test:**
```bash
go test ./dspy -v -run TestNewPredictor
```

**Run tests with coverage:**
```bash
go test ./... -cover
```

**Generate detailed coverage report:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Current test coverage:
- `dspy`: 78.9% coverage
- `llm`: 17.3% coverage (HTTP clients require integration testing)
- `optimizer`: 86.5% coverage
- `memory`: 92.6% coverage
- `tracing`: 100% coverage

## CI/CD

This project uses GitHub Actions for continuous integration:

- **Tests**: Automatically runs on every push and pull request
  - Tests across Go 1.21 and 1.22
  - Race detector enabled
  - Coverage reports generated
  - Results displayed in GitHub UI

- **Linting**: Code quality checks
  - `gofmt` formatting validation
  - `go vet` static analysis
  - Module verification

- **Build**: Ensures all packages compile successfully

View test results and coverage in the [Actions tab](https://github.com/supadev-ai/go-dspy/actions).

## Roadmap

### v0.1 (Current)
- ✅ Signatures
- ✅ Predictor
- ✅ OpenAI client
- ✅ Bootstrap optimizer
- ✅ Examples

### v0.2 (Planned)
- [ ] Teleprompter optimizer
- [ ] Multi-module pipelines
- [ ] Memory abstraction
- [ ] Chain modules

### v0.3 (Future)
- [ ] DAG execution
- [ ] Tracing (OpenTelemetry)
- [ ] Prompt caching
- [ ] More LLM providers (Azure, Ollama)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache 2.0 - see [LICENSE](LICENSE) for details.

## About SupaDev

**go-dspy** is developed and maintained by [SupaDev](https://github.com/supadev-ai), an AI services company specializing in production-ready AI infrastructure and tooling. SupaDev helps organizations build, deploy, and scale AI applications with enterprise-grade reliability and performance.

## Acknowledgments

Inspired by [DSPy](https://github.com/stanfordnlp/dspy) from Stanford NLP.
