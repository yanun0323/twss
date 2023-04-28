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

.PHONY: mongo.local.run
mongo.local.run:
	docker run -d -p 27017:27017 \
	-e MONGODB_CLIENT_EXTRA_FLAGS=--authenticationDatabase=admin \
	-e MONGO_INITDB_ROOT_USERNAME=local -e MONGO_INITDB_ROOT_PASSWORD=local \
	--name mongo mongo