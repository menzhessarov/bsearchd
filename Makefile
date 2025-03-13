default: run

run:
	@go run cmd/main.go

test:
	@go test -v -race .

.PHONY: default run test