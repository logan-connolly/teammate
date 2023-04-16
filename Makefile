.DEFAULT_GOAL=help

test: #> Run test suite
	@go test -count=1 -race -shuffle=on ./...

coverage: #> Run test suite and show coverage
	@go test -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...
	@go tool cover -html coverage.txt

help: #> Show this help
	@echo
	@echo -e "\033[0;34m Teammate\033[0m 🏅"
	@echo
	@grep -E -h '\s#>\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?#> "}; {printf " \033[36m%-20s\033[0m %s\n", $$1, $$2}'
