include .env

DEV_DOCKER_COMPOSE_PATH = ./deployments/docker-compose.yml
PROD_DOCKER_COMPOSE_PATH = ./deployments/docker-compose.prod.yml

up:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) up --build -d

connect:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app dlv connect app:40000

stop:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) stop

down:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) down -v --remove-orphans

docker-rm-volume:
	docker volume rm telegram-bot_mysql-data

migrations-up:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app migrate -path migrations -database 'mysql://$(MYSQL_DSN)' up

migrations-down:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app migrate -path migrations -database 'mysql://$(MYSQL_DSN)' down

migrations-force-v:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app migrate -path migrations -database 'mysql://$(MYSQL_DSN)' force $(version)

migrations-create:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app migrate create -ext sql -dir migrations '$(name)'

mock:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app mockery

test:
	docker compose -f $(DEV_DOCKER_COMPOSE_PATH) run --rm app go test ./...

up-prod:
	docker compose -f $(PROD_DOCKER_COMPOSE_PATH) up --build -d

stop-prod:
	docker compose -f $(PROD_DOCKER_COMPOSE_PATH) stop

down-prod:
	docker compose -f $(PROD_DOCKER_COMPOSE_PATH) down -v --remove-orphans

migrations-up-prod:
	docker compose -f $(PROD_DOCKER_COMPOSE_PATH) run --rm migration migrate -path migrations -database 'mysql://$(MYSQL_DSN)' up

migrations-down-prod:
	docker compose -f $(PROD_DOCKER_COMPOSE_PATH) run --rm migration migrate -path migrations -database 'mysql://$(MYSQL_DSN)' down

migrations-force-v-prod:
	docker compose -f $(PROD_DOCKER_COMPOSE_PATH) run --rm migration migrate -path migrations -database 'mysql://$(MYSQL_DSN)' force $(version)