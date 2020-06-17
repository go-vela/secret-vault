# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

build: binary-build

up: build compose-up

down: compose-down

seed: vault-seed

run: binary-build docker-build docker-run

test: binary-build docker-build docker-example


clean:
	#################################
	######      Go clean       ######
	#################################

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

binary-build:
	#################################
	######    Build Binary     ######
	#################################

	GOOS=linux CGO_ENABLED=0 go build -o release/secret-vault github.com/go-vela/secret-vault/cmd/secret-vault

docker-build:
	#################################
	######    Docker Build     ######
	#################################

	docker build --no-cache -t secret-vault:local .

compose-up:
	#################################
	###### Docker Build/Start  ######
	#################################

	@docker-compose -f docker-compose.yml.vault up -d # start a Vault app

compose-down:
	#################################
	###### Docker Tear Down    ######
	#################################

	@docker-compose -f docker-compose.yml.vault down	

vault-seed:
	#################################
	######  Vault Seed Data    ######
	#################################

	VAULT_ADDR=http://localhost:8200 \
		VAULT_TOKEN=vela \
		vault write secret/my-secret foo=bar

docker-run:
	#################################
	######     Docker Run      ######
	#################################

	docker run --rm \
		--network secret-vault_vault \
		-e PARAMETER_LOG_LEVEL \
		-e PARAMETER_ADDR \
		-e PARAMETER_TOKEN \
		-e PARAMETER_PATH \
		-e PARAMETER_KEYS \
		secret-vault:local	

docker-example:
	#################################
	######   Docker Example    ######
	#################################

	docker run \
		--network secret-vault_vault \
		-e PARAMETER_LOG_LEVEL=trace \
		-e PARAMETER_ADDR=http://vault:8200 \
		-e PARAMETER_TOKEN=vela \
		-e PARAMETER_PATH=secret/my-secret  \
		-e PARAMETER_KEYS=foo \
		secret-vault:local	