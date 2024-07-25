test:
	go test -v ./...

up:
	docker compose up -d

down:
	docker compose down

run:
	go run cmd/skillBuilder/main.go

migrate:
	goose -dir migrations postgres "postgres://root:root@localhost:54322/app" up

port:
	lsof -i :8080