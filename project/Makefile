FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp

up:
	@echo "Starting Docker images"
	docker compose up -d
	echo "Docker images started"

up_build:
	@echo "Stopping the existing docker images"
	docker compose down
	@echo "Building the docker images"
	docker compose up --build -d
	@echo "Docker images built successfully"

down:
	@echo "Stoping the docker compose"
	docker-compose down
	@echo "Done!"

build_broker:
	@echo "Building the broker binary"
	cd ../broker && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

build_frontend: up
	@echo "Building the frontend binary"
	cd ../frontend/ && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

start: build_frontend
	@echo "Starting the frontend"
	cd ../frontend && ./${FRONT_END_BINARY} 


