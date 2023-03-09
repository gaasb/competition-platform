.SILENT:
.DEFAULT_GOAL := run
build:
	go mod download && CGO_ENABLED=0 go build -o ./.bin/app ./cmd/api/api.go
run:
	go run ./cmd/api/api.go
test:
	go test -v ./...
migrate:
	migrate
migrate-create:
	migrate create -seq $(NAME)
migrate-revert:
deploy: