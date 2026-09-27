# Contributing Guide

Thank you for your interest in contributing to Diary! This project is crafted in Go with a strong commitment to simplicity, local-first privacy, clean architecture, and zero external runtime dependencies.

## Getting Started

1. Fork the repository on GitHub: [github.com/Felipe-Takayuki/diary](https://github.com/Felipe-Takayuki/diary)
2. Clone your fork to your machine:
   ```bash
   git clone https://github.com/Felipe-Takayuki/diary.git
   cd diary
   ```
3. Create a branch for your feature or bugfix:
   ```bash
   git checkout -b feature/my-improvement
   ```

## Project Principles

- **Zero runtime dependencies:** The core and server rely exclusively on the Go standard library, while web assets use vanilla HTML, CSS, and modern JavaScript. Avoid introducing heavy external frameworks.
- **Clean Architecture:** Maintain domain logic pure in `internal/domain`, application orchestration in `internal/usecase`, and delivery adapters in `internal/adapter`.
- **Transparent Markdown persistence:** Goals are saved locally as plain Markdown checklist files in `metas/DD-MM-YYYY.md` (e.g., `- [ ]` and `- [x]`).
- **Visual accessibility & Omarchy integration:** Interface components must follow high-contrast guidelines (WCAG AAA), complete keyboard navigation, and seamless light/dark theme adaptation.

## Local Verification

Before opening a Pull Request, ensure that all checks pass:

```bash
# Run all unit tests
make test

# Run tests with the Go race detector
make test-race

# Check test coverage
make test-cover

# Format code and run static analysis
make fmt
make vet

# Build the desktop binary
make build
```

## Submitting a Pull Request

1. Follow the **Conventional Commits** specification (e.g., `feat:`, `fix:`, `docs:`, `style:`, `refactor:`, or `test:`).
2. Include unit tests for any new business logic, domain models, or use cases.
3. Open a Pull Request with a clear description of the problem solved, design choices, and screenshots if UI changes were made.
