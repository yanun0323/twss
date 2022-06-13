# CONFIG
.PHONY: run
run:
	CONFIG_NAME=config go run ./main.go

.PHONY: update
update:
	CONFIG_NAME=config MODE=update go run ./main.go

.PHONY: crawl
crawl:
	CONFIG_NAME=config MODE=crawl go run ./main.go

.PHONY: convert
convert:
	CONFIG_NAME=config MODE=convert go run ./main.go

.PHONY: check
check:
	CONFIG_NAME=config MODE=check go run ./main.go

# CONFIG-PI

.PHONY: run_pi
run_pi:
	CONFIG_NAME=config-pi go run ./main.go

.PHONY: update_pi
update_pi:
	CONFIG_NAME=config-pi MODE=update go run ./main.go

.PHONY: check_pi
check_pi:
	CONFIG_NAME=config MODE=check go run ./main.go