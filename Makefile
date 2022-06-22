# CONFIG
.PHONY: run
run:
	CONFIG_NAME=config go run ./main.go

.PHONY: update crawl convert check
update:
	CONFIG_NAME=config MODE=update go run ./main.go

crawl:
	CONFIG_NAME=config MODE=crawl go run ./main.go

convert:
	CONFIG_NAME=config MODE=convert go run ./main.go

check:
	CONFIG_NAME=config MODE=check go run ./main.go

# CONFIG-Production

.PHONY: run_production
run_production:
	CONFIG_NAME=config-production go run ./main.go

.PHONY: update_production check_production convert_production
update_production:
	CONFIG_NAME=config-production MODE=update go run ./main.go

check_production:
	CONFIG_NAME=config-production MODE=check go run ./main.go
	
convert_production:
	CONFIG_NAME=config-production MODE=convert go run ./main.go