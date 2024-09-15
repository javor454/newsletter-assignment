.PHONY: up down

up: ## Start up application container
	#docker compose build --no-cache --progress=plain
	docker compose build
	docker compose up

down: ## Stop and remove the application containers
	docker compose down --volumes --remove-orphans
