# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`goparsec` is a parser combinator library written in Go. Module path: `github.com/ajiyoshi-vg/goparsec`.

## Development Workflow

### Branching

Always work on a feature branch. Never commit or push directly to `master`.
Every change must go through a pull request.

### TDD

Implement using t-wada style TDD: write a failing test first, then make it pass, then refactor. Do not write implementation code without a test driving it.

### Issue Tracking

When you notice something worth tracking (a design question, a potential improvement, a known limitation), save it to `./issues/YYYYMMDD/{title}.md` using today's date.

## Commands

```bash
# Run all tests
go test ./...

# Run a single package's tests
go test ./path/to/pkg

# Run a single test by name
go test ./... -run TestName

# Build
go build ./...

# Lint (requires golangci-lint)
golangci-lint run

# Format
gofmt -w .
```

## Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/) を使う。

```
<type>(<scope>): <subject>
```

`type`: `feat` / `fix` / `refactor` / `test` / `docs` / `chore`

## Toolchain

Managed by `mise` (`mise.toml`). Go version: `latest`.
