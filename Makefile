.PHONY: run
run:
	go run ./main.go

.PHONY: update
update:
	go run ./main/daily_update/main.go

.PHONY: crawl
crawl:
	go run ./main/crawler/main.go

.PHONY: convert
convert:
	go run ./main/converter/main.go