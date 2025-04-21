//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "mvdan.cc/gofumpt"
)

// This file imports tool dependencies.
// It ensures they're tracked in go.mod but not included in the build.
// Run 'go install' to install them:
// go install github.com/golangci/golangci-lint/cmd/golangci-lint
// go install mvdan.cc/gofumpt
