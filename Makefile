.PHONY: up down swagger

up:
	docker compose up -d --build --wait
	docker compose logs -f

down:
	docker compose down

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.3 init -g cmd/server/main.go -o docs
