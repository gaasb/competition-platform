FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -v -o ./bin/build ./...

FROM alpine:latest AS prod
CMD ["/app/bin/prod", "port"]
