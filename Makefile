.PHONY: run build test docker-up docker-down clean build-worker run-worker

dev:
	go tool air  --build.cmd="go build -o tmp/pulse-server cmd/server/main.go" --build.bin="tmp/pulse-server"

run:
	go run cmd/server/main.go

build:
	mkdir -p build
	go build -o build/server cmd/server/main.go

dev-worker:
	go tool air  --build.cmd="go build -o tmp/pulse-worker cmd/worker/main.go" --build.bin="tmp/pulse-worker"

run-worker:
	go run cmd/worker/main.go

build-worker:
	mkdir -p build
	go build -o build/worker cmd/worker/main.go

test:
	go test ./...

docker-up:
	docker compose -f docker-compose.infrastructure.yml up -d

docker-down:
	docker compose -f docker-compose.infrastructure.yml down

clean:
	rm -rf build
	rm -rf tmp
