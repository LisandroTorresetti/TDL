SHELL := /bin/bash
PWD := $(shell pwd)

mandale-mecha:
	docker-compose -f docker-compose.yml up -d --build
.PHONY: mandale-mecha

stop-app:
	docker-compose -f docker-compose.yml stop -t 1
.PHONY: stop-app

delete-app: stop-app
	docker-compose -f docker-compose.yml down
	echo "TE ESTAS PORTANDO MAL SERAS CASTIGADO"
.PHONY: delete-app