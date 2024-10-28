test:
	go test -v ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

start:
	docker compose up --build

migrate:
	migrate -source file://DB/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) up 1
