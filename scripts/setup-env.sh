#!/usr/bin/env bash
set -euo pipefail

echo "🚀 install & update gofumpt ..."
go install mvdan.cc/gofumpt@latest
echo "✅ done"

echo "🚀 install & update golangci-lint ..."
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
echo "✅ done"

echo "🚀 install & update goimports ..."
go install golang.org/x/tools/cmd/goimports@latest
echo "✅ done"

echo "🚀 install & update mockgen ..."
go install go.uber.org/mock/mockgen@latest
echo "✅ done"

echo "🎉 setup tools complete"
