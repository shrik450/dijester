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
- Comment sparingly and only when necessary. Comments should only be used to
  document or clarify when the code itself is not self-explanatory.

  To best follow this directive, you must:

  1. Write a first pass off the code as you would normally do.
  2. Strip all comments from it.
  3. Go over the code again and add documentation comments, along with comments
     ONLY where the code cannot be made self-explanatory.

  If the code you generate after step 2 has ANY comments, it will be
  automatically rejected.

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
