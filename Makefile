.PHONY: run
run:
	go run main.go

.PHONY: test.all
test.all:
	go test --count=1 ./...