# Contributing to Taproot

Thank you for your interest in contributing to Taproot TUI Framework! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)

## Code of Conduct

Be respectful, inclusive, and collaborative. We aim to maintain a welcoming environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- A terminal that supports ANSI colors

### Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/taproot.git
   cd taproot
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/yourorg/taproot.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

### Building

```bash
# Build all packages
go build ./...

# Run an example
go run examples/demo/main.go
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Follow coding standards (see below)
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/tui/styles/...
```

### 4. Commit Changes

See [Commit Messages](#commit-messages) for guidelines.

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Coding Standards

### File Organization

```
internal/tui/
├── components/       # UI components
├── styles/          # Theme system
├── util/            # Utilities
├── page/            # Page management
├── app/             # Application framework
└── exp/             # Experimental components
```

### Naming Conventions

**Packages**: Lowercase, single word when possible
```go
package styles
package util
```

**Files**: `snake_case.go` for examples, `lowercase.go` for code
```go
internal/tui/styles/theme.go
internal/tui/styles/markdown.go
examples/demo/main.go
```

**Interfaces**: `-able` suffix for capabilities
```go
type Focusable interface { ... }
type Sizeable interface { ... }
type Positional interface { ... }
```

**Functions**: PascalCase for exported, camelCase for internal
```go
func GetMarkdownRenderer() *Renderer { ... }
func buildStyles() *Styles { ... }
```

**Constants**: UPPER_CASE for exported
```go
const MAX_WIDTH = 120
```

### Code Style

- Use `gofmt` for formatting
- Follow Go conventions from [Effective Go](https://go.dev/doc/effective_go)
- Keep functions focused and small
- Add comments for exported functions
- Use table-driven tests for multiple scenarios

### Example

```go
// Package styles provides theme management and styling for TUI components.
package styles

import (
    "github.com/charmbracelet/lipgloss"
)

// GetMarkdownRenderer returns a glamour TermRenderer configured with the current theme.
// The width parameter specifies the maximum line width for wrapped text.
func GetMarkdownRenderer(width int) *glamour.TermRenderer {
    t := CurrentTheme()
    r, _ := glamour.NewTermRenderer(
        glamour.WithStyles(t.S().Markdown),
        glamour.WithWordWrap(width),
    )
    return r
}
```

## Testing Guidelines

### Write Tests For

- New functions and methods
- Bug fixes
- Complex logic
- Public APIs

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
    }{
        {"case 1", input1, expected1},
        {"case 2", input2, expected2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            if result != tt.expected {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Goals

- Aim for >70% coverage for new code
- Test edge cases and error conditions
- Use mocks for external dependencies
- Keep tests fast and deterministic

## Commit Messages

Follow conventional commit format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process, tooling, dependencies

### Examples

```
feat(styles): add markdown rendering support

Implement glamour-based markdown rendering with theme integration.
Add GetMarkdownRenderer() and PlainMarkdownStyle() functions.

Closes #123
```

```
fix(dialogs): resolve stack overflow on rapid ESC presses

Limit dialog stack depth to prevent overflow when user rapidly
presses ESC key.

Fixes #456
```

## Pull Requests

### PR Title

Use the same format as commit messages:
```
feat(styles): add markdown rendering support
```

### PR Description

Include:
- **What**: Brief description of changes
- **Why**: Reason for the change
- **How**: Implementation approach
- **Testing**: How changes were tested
- **Breaking Changes**: Note if applicable

### PR Checklist

- [ ] Tests pass locally
- [ ] New tests added for new features
- [ ] Documentation updated
- [ ] Commits follow message format
- [ ] No merge conflicts

## Review Process

1. Automated checks must pass
2. At least one maintainer approval required
3. Address review comments
4. Squash commits if needed
5. Maintainer merges

## Getting Help

- Check existing documentation in `docs/`
- Review example code in `examples/`
- Open an issue for bugs or feature requests
- Ask questions in GitHub Discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
