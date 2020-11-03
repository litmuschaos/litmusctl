PROJECT_NAME := "kuberactl"
PKG := "github.com/mayadata-io/$(PROJECT_NAME)"

all: build

dep: ## Get the dependencies
	@go mod download

build: dep ## Build the binary file
	@go build -i -o build/main $(PKG)