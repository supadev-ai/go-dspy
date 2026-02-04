# Contributing to go-dspy

Thank you for your interest in contributing to go-dspy! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Local Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/go-dspy.git
   cd go-dspy
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run tests:
   ```bash
   go test ./...
   ```

## Development Workflow

### Code Style

- Follow standard Go formatting: `go fmt ./...`
- Use `golangci-lint` for linting (if configured)
- Follow Go naming conventions and idioms
- Write clear, self-documenting code

### Testing

- Write tests for new features
- Ensure all tests pass: `go test ./...`
- Aim for good test coverage
- Use table-driven tests where appropriate

### Commit Messages

- Use clear, descriptive commit messages
- Reference issues when applicable: `Fix #123`
- Follow conventional commit format when possible

### Pull Requests

1. Create a feature branch from `main`
2. Make your changes
3. Ensure tests pass and code is formatted
4. Update documentation if needed
5. Submit a pull request with a clear description

## Project Structure

```
go-dspy/
├── cmd/examples/     # Example programs
├── dspy/             # Core DSPy types and interfaces
├── llm/              # LLM client implementations
├── optimizer/        # Optimization algorithms
├── memory/           # Memory/store abstractions
└── tracing/          # Observability and tracing
```

## Areas for Contribution

- **LLM Providers**: Add support for new providers (Azure, Ollama, etc.)
- **Optimizers**: Implement new optimization strategies
- **Examples**: Add more example programs
- **Documentation**: Improve docs, add tutorials
- **Testing**: Increase test coverage
- **Performance**: Optimize hot paths

## Questions?

Open an issue for questions, bug reports, or feature requests.

## Code of Conduct

Be respectful and inclusive. We're all here to build something great together.
