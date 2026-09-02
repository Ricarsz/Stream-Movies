build:
	@go build -o bin/fs ./cmd

run: build
	@./bin/fs

test:
	@GO_TEST=1 go test ./tests/ -v -count=1
