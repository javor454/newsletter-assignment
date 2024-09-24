.PHONY: up down swag test

up: ## Start up application container
	#docker compose build --no-cache --progress=plain
	docker compose build
	docker compose up

down: ## Stop and remove the application containers
	docker compose down --volumes --remove-orphans

# TODO: dockerize
swag: ## Format and build swagger docs
	swag fmt
	swag init

test: ## Run all tests
	./script/test.sh