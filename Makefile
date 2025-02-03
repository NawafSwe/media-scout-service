export GO111MODULE=on

help: ## This help dialog.
help h:
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%-20s %s\n" "target" "help" ; \
	printf "%-20s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
		IFS=$$':' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf '\033[36m'; \
		printf "%-20s %s" $$help_command ; \
		printf '\033[0m'; \
		printf "%s\n" $$help_info; \
	done



#===================#
#== Env Variables ==#
#===================#
DOCKER_COMPOSE_FILE ?= docker-compose.yml
RUN_IN_DOCKER ?= docker-compose exec builder

#===============#
#== App Build ==#
#===============#

build: build-app

build-native: RUN_IN_DOCKER=
build-native: build

build-app: ## Build binaries for the domain app
build-app:
	@echo "Building app binary"
	${RUN_IN_DOCKER} go build -o ./bin/media-scout ./cmd/main.go

mock:
	@echo "=========================================="
	@echo "Generating mocks using gomock"
	@echo "=========================================="
	mockgen -source=${source} -destination=${destination} -package=${package}

#===============#
#=== App Run ===#
#===============#

http: ## Run http server
	${RUN_IN_DOCKER} sh -c 'bin/media-scout'

#===============#
#=== Apply Migrations ===#
#===============#
migrate: ## Run migrations against non test DB
	docker-compose up migrate

migrate-create: ## Create a DB migration. You need to pass the file name, e.g. `make migrate-create name=migration-name`
	docker-compose run --rm migrate create -ext sql -dir /migrations $(name)

#=======================#
#== ENVIRONMENT SETUP ==#
#=======================#

create-env-file:
ifeq (,$(wildcard .env))
	cp .env.sample .env
endif

clean: # clean executables
	rm -rf bin/*

docker-ready:
	@echo "Waiting until Docker is ready..."
	@until docker version --format 'Server version: {{.Server.Version}}' >/dev/null 2>&1; do \
		echo "Waiting for Docker..." ; \
		sleep 1; \
	done
	@echo "Docker is ready"

docker-start:
	@echo "Starting Docker Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} up -d --build --remove-orphans
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

docker-stop:
	@echo "Stopping Docker Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} stop
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

docker-clean: docker-stop
	@echo "Removing Docker Containers"
	docker-compose -f ${DOCKER_COMPOSE_FILE} rm -v -f

docker-restart: docker-stop docker-start

environment: ## The only command needed to start a working environment
environment: docker-ready create-env-file docker-restart build-app


clean-environment: ## The only command needed to clean the environment
clean-environment: docker-clean clean

tests-unit:
	@echo "=================="
	@echo "Running unit tests"
	@echo "=================="
	go test -tags unit -shuffle=on -coverprofile coverage.out ./...

tests-integration: environment
	@echo "======================================="
	@echo "Running integration tests with coverage"
	@echo "======================================="
	docker exec media-scout-app sh -c "cp ./test/.env.test ./test/test/.env"
	docker exec media-scout-app sh -c "go clean -cache"
	docker exec media-scout-app sh -c "go test -v ./... -tags=integration -coverprofile=coverage.out"

#====================#
#== QUALITY CHECKS ==#
#====================#

format:
	@echo "=========================================="
	@echo "Formatting your code"
	@echo "=========================================="
	gci write -s standard -s default . --skip-generated --skip-vendor  && gofumpt -l -w .