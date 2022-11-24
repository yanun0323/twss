.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test --count=1 ./...

.PHONY: test.debug
test.debug:
	go test -v --count=1 ./...