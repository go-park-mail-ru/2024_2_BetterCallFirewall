test:
	go test -v ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

start:
	docker compose up --build

stop:
	docker compose stop

restart: stop start

gen-proto:
	protoc \
     	--go_out=internal/api/grpc --go_opt=paths=import --go_opt=module=github.com/2024_2_BetterCallFirewall/internal/api/grpc \
     	--go-grpc_out=internal/api/grpc --go-grpc_opt=paths=import --go-grpc_opt=module=github.com/2024_2_BetterCallFirewall/internal/api/grpc \
      	proto/*.proto

lint:
	golangci-lint run

gen-easy-json:
	easyjson -all internal/models/*.go