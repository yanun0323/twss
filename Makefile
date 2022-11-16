.PHONY: run run.server
run:
	go run main.go
run.server:
	MODE=server go run main.go

.PHONY: test.all
test.all:
	go test --count=1 ./...