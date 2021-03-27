all: build

build: ## Build the binary file
	@bash scripts/build.sh main.go $(TAG)

.PHONY: unused-package-check
unused-package-check:
	@echo "------------------"
	@echo "--> Check unused packages for the litmusctl"
	@echo "------------------"
	@tidy=$$(go mod tidy); \
	if [ -n "$${tidy}" ]; then \
		echo "go mod tidy checking failed!"; echo "$${tidy}"; echo; \
	fi