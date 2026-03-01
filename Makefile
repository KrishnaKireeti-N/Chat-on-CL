all: build listen

build:
	@go build -o bin/ ./cmd/chat/

listen:
	@./bin/chat -l --color Cyan
