.PHONY: clean critic security lint test build run

APP_NAME = chatapp
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/internal/migration
DATABASE_URL=postgresql://myuser:pass@localhost:5432/postgres?sslmode=disable


clean:
	rm -rf ./server/build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

#run: swag build
#	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

docker.run: docker.postgres docker.redis docker.kafka

docker.fiber.build:
	docker build -t fiber .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name cgapp-fiber \
		-p 8080:8080 \
		fiber

docker.postgres:
	docker run --rm -d \
		--name cgapp-postgres\
		-e POSTGRES_USER=myuser\
		-e POSTGRES_PASSWORD=pass\
		-e POSTGRES_DB=postgres\
		-p 5432:5432\
		postgres

docker.redis:
	docker run --rm -d \
		--name cgapp-redis \
		-p 6379:6379 \
		redis

docker.kafka:
	docker run --rm -d \
		--name cgapp-kafka \
		-p 9092:9092 \
		-e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 \
		-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
		-e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
		bitnami/kafka:latest

docker.stop: docker.stop.postgres docker.stop.redis docker.stop.kafka migrate.down

docker.stop.fiber:
	docker stop cgapp-fiber

docker.stop.postgres:
	docker stop cgapp-postgres

docker.stop.redis:
	docker stop cgapp-redis

docker.stop.kafka:
	docker stop cgapp-kafka

#swag:
#	swag init
