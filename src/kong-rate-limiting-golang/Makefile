.DEFAULT: help

IMAGE_NAME ?= lampnick/kong-rate-limiting-plugin-golang
CENTOS_IMAGE_TAG ?= kong-2.1.3-centos-v1.2.1
ALPINE_IMAGE_TAG ?= kong-2.1.3-alpine-v1.2.1

KONG_NET_NAME := $(shell docker network ls|grep kong-net|awk '{print $$2}')
P="\\033[34m===>\\033[0m"

help:
	@echo
	@echo "  \033[34mbuild-centos\033[0m – build kong rate limiting plugin in centos docker"
	@echo "  \033[34mbuild-alpine\033[0m – build kong rate limiting plugin in alpine docker"
	@echo "  \033[34mtest-run-centos\033[0m – run centos builded kong"
	@echo "  \033[34mtest-run-alpine\033[0m – run alpine builded kong"
	@echo "  \033[34mtest-plugins\033[0m – test kong rate limiting plugins is running."
	@echo "  \033[34mrun-kong-konga-pg\033[0m – have a try with postgres,kong,konga in a command."
	@echo

build-centos:
	@docker build -t $(IMAGE_NAME):$(CENTOS_IMAGE_TAG) .

build-alpine:
	@docker build -f Dockerfile.alpine -t $(IMAGE_NAME):$(ALPINE_IMAGE_TAG) .

test-run-centos:
	@docker run --rm --name kong-rate-limiting-plugin-golang \
    -e "KONG_LOG_LEVEL=debug" \
    -e "KONG_NGINX_USER=root root" \
    -p 8000:8000 \
    -p 8443:8443 \
    -p 8001:8001 \
    -p 8444:8444 \
     $(IMAGE_NAME):$(CENTOS_IMAGE_TAG)

test-run-alpine:
	@docker run --rm --name kong-rate-limiting-plugin-golang \
    -e "KONG_LOG_LEVEL=info" \
    -e "KONG_NGINX_USER=root root" \
    -p 8000:8000 \
    -p 8443:8443 \
    -p 8001:8001 \
    -p 8444:8444 \
    $(IMAGE_NAME):$(ALPINE_IMAGE_TAG)

test-plugins:
	@curl -s http://localhost:8001/ |grep --color custom-rate-limiting

rm-kong-net:
ifeq ($(KONG_NET_NAME),kong-net)
	@echo $(P) "rm kong net";
	@docker network rm kong-net
endif

create-kong-net: rm-kong-net
	@echo $(P) "create kong net";
	@docker network create kong-net

run-postgres: create-kong-net
	@echo $(P) "run postgres";
	@docker run -d \
    --name kong-database \
    --network=kong-net \
    -p 5432:5432 \
    -e "POSTGRES_USER=kong" \
    -e "POSTGRES_DB=kong" \
    -e "POSTGRES_PASSWORD=kong" \
    postgres:9.6
	@echo $(P) "preparing postgres";

.ONESHELL:
sleep: run-postgres
	@echo $(P) "postgres ok"  $(shell sleep 10)

.ONESHELL:
run-migration: sleep
	@echo $(P) "run migration";
	@docker run --rm \
    --network=kong-net \
    -e "KONG_DATABASE=postgres" \
    -e "KONG_PG_HOST=kong-database" \
    -e "KONG_PG_USER=kong" \
    -e "KONG_PG_PASSWORD=kong" \
    -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
    $(IMAGE_NAME):$(CENTOS_IMAGE_TAG) kong migrations bootstrap

run-kong: run-migration
	@echo $(P) "run kong";
	@docker run -d --name kong-rate-limiting-plugin-golang \
    -e "KONG_LOG_LEVEL=info" \
    --network=kong-net \
    -e "KONG_NGINX_USER=root root" \
    -e "KONG_DATABASE=postgres" \
    -e "KONG_PG_HOST=kong-database" \
    -e "KONG_PG_USER=kong" \
    -e "KONG_PG_PASSWORD=kong" \
    -p 8000:8000 \
    -p 8443:8443 \
    -p 8001:8001 \
    -p 8444:8444 \
    $(IMAGE_NAME):$(CENTOS_IMAGE_TAG)

run-konga: run-kong
	@echo $(P) "run konga";
	@docker run -d \
    -p 1337:1337 \
    --network kong-net \
    -e "TOKEN_SECRET=kongtoken" \
    -e "DB_ADAPTER=postgres" \
    -e "DB_HOST=kong-database" \
    -e "DB_USER=kong"  \
    -e "DB_PASSWORD=kong" \
     --name konga \
    pantsel/konga

.ONESHELL:
run-kong-konga-pg: run-konga
	@echo $(P) "please visit: http://localhost:1337 to visit konga"