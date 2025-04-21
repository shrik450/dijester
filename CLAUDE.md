# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

- Install tools: `make tools`
- Build: `make build`
- Run: `make run`
- Test all: `make test`
- Test single: `go test ./path/to/package -run TestName`
- Lint: `make lint`
- Format: `make fmt`
- End-to-end tests: `make e2e-test`
- Format + Lint + Test: `make check`

Tools are automatically installed via Go modules.

## Code Quality Guidelines

- Always run `make check` after each major change
- First run unit tests, then end-to-end tests after significant changes
- Write meaningful tests focused on complex functionality, not simple types
- Use interfaces to make components easy to test

## Code Guidelines

- **Formatting**: Use `gofumpt` (stricter gofmt) with 100 character line length limit
- **Imports**: Group standard library, third-party, and local imports
- **Error Handling**: Always check errors, prefer explicit error returns over panics
- **Types**: Use strong typing, define custom types with descriptive names
- **Functions**: Keep functions small and focused on single responsibility
- **Naming**: Use camelCase for unexported and PascalCase for exported names
- **Comments**: Document all exported functions, types, and packages
- **Architecture**: Follow clean architecture principles with clear separation of concerns

## Process Guidelines

- After every change, make a commit.
- Follow a trunk-based development model around the `main` branch.
- If the change you are making differs in purpose from the last commit, switch
  back to main, create a new branch and then commit your changes.
- Commit messages should be clear and concise. A body is not required unless it
  substantially clarifies the change. The first line MUST be less than 72
  characters.
- Once you are done with a chunk of changes related to a single purpose, create
  a pull request. The pull request should be clear and concise. It should
  summarize the changes you made and why they are necessary. Use the `gh` CLI to
  post PRs.

