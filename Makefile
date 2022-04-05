.PHONY: run
run:
	go run ./main.go

crawl:
	go run ./main/crawler/main.go

convert:
	go run ./main/converter/main.go