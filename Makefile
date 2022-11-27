.PHONY: run check job debug
run:
	MODE=server go run main.go
check:
	MODE=check go run main.go
job:
	MODE=job go run main.go
debug:
	MODE=debug go run main.go

.PHONY: test test.debug
test:
	go test --count=1 ./...
test.debug:
	go test -v --count=1 ./...