default: run

run:
	go run cmd/cli/main.go

up:
	docker compose up -d --build
