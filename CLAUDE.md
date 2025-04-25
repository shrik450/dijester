# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Build & Test Commands

- Install tools: `make tools`
- Build: `make build`
- Run: `make run`
- Test all: `make test`
- Test single: `go test ./path/to/package -run TestName`
- Lint: `make lint`
- Format: `make fmt`
- End-to-end tests: `make e2e-test`

## Code Quality Guidelines

- Always run `make check` after each major change.
- Do not add ANY comments. Doc comments are allowed, but NO OTHER COMMENTS.
- Your implementation should be as simple as possible. Do not over-engineer or
  over-complicate the code for the future. However, if additional complexity is
  required to make the code testable, it is acceptable. E.g. taking an interface
  as a parameter instead of a concrete type with the intention to mock it in
  tests.
- Follow the following order for any new functionality:
  1. Is there a standard library function that does this? If yes, use it. E.g.
     instead of combining strings, use `tmpl` or `strings.Join`.
  2. Is this functionality core to the working of this program, or simple
     enough that it can be implemented in a single function *without*
     compromising on edge cases or security? If yes, implement it yourself. E.g.
     fetching a URL, structuring the output, etc.
  3. Is there a high-quality, commonly used open-source library that does this?
     If yes, use it.
  4. If none of the above options are available, implement it yourself in a new
     package.
- Write unit, integration, and end-to-end tests for new features and bug fixes.
  Tests should be meaningful and cover a variety of scenarios. Prefer testing
  at higher complexity levels. Do not write tests that test getters, setters,
  trivial functions or simple types. If a test has to be changed in a refactor,
  it is probably a bad test.
- Run all tests after every change.
- Use interfaces to make components easy to test.

## Code Guidelines

- **Formatting**: Use `gofumpt` with 100 character line length limit
- **Imports**: Group standard library, third-party, and local imports
- **Error Handling**: Always check errors, prefer explicit error returns over
  panics
- **Types**: Use strong typing, define custom types with descriptive names
- **Functions**: Keep functions small and focused on single responsibility
- **Naming**: Use camelCase for unexported and PascalCase for exported names
- **Comments**: Document all exported functions, types, and packages
- **Architecture**: Follow clean architecture principles with clear separation
  of concerns

## Process Guidelines

- Follow a trunk-based development model around the `main` branch, with feature
  branches and pull requests. This repository follows a squash-and-merge
  strategy for a linear history.
- Before making any changes:

  1. Run all checks to ensure you are in a known good state.
  2. Run `git status` to check the git status. Make sure you are on a relevant
     branch and that there are no uncommitted changes. NEVER commit on main. If
     you are not on a relevant branch for this change, switch to the main branch
     and either create a new branch or switch to the existing relevant branch.
  3. If there are uncommitted changes, and you are going to switch branches,
     create a WIP commit and push it to avoid losing any work. If you are not
     going to switch branches, either build on top of the existing changes or
     reset them entirely.

- After every atomic change, make a commit. The commit MUST include the tests
  for the change if there are any.

  Commit message guidelines:

  1. Do NOT use conventional commits.
  2. The first line must be a short summary and less than 72 characters.
  3. A body is not always required. Do not make a list of all the changes in
     the body. Only include a body if the motivation for the change is not
     immediately clear from the commit message, and only use it to clarify that
     motivation.
  4. Always include attribution to claude code in commits that you have
     authored. If the changes are not made by you but committed by you, do NOT
     include that attribution.

- Once you are done with a chunk of changes related to a single purpose, create
  a GitHub pull request using the `gh` CLI, which should already be configured.
  Ideally, there is one commit per pull request. If you have multiple commits,
  squash them into one first unless they differ substantially in purpose. The
  content of the pull request should follow the same guidelines as commit
  messages.
