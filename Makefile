.PHONY: deps bin_build bin_run bin_clean
deps:
	@go mod tidy
	@go mod download

bin_build:
	@go build -o bin/rybakov_test cmd/main.go

bin_run:
	@bin/rybakov_test

bin_clean:
	@rm -f bin/rybakov_test