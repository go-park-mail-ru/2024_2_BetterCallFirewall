test:
	go test -v ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

start:
	docker compose up --build

stop:
	docker compose stop

restart: stop start