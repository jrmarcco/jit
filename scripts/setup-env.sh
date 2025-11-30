#!/bin/sh

echo "🚀 install & update gofumpt ..."
go install mvdan.cc/gofumpt@latest
echo "✅ done"

echo "🚀 install & update golangci-lint ..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
echo "✅ done"

echo "🚀 install & update goimports ..."
go install golang.org/x/tools/cmd/goimports@latest
echo "✅ done"

echo "🚀 install & update mockgen ..."
go install github.com/golang/mock/mockgen@latest
echo "✅ done"

echo "🎉 setup tools complete"
