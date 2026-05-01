find . -name '*.go' -type f -exec gofmt -s -w {} \;
golangci-lint run --enable goimports --enable gofmt --timeout 10m