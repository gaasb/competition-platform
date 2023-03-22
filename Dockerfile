FROM golang:1.20.0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mode verify
COPY . .
RUN go build -v -o app/bin/app ./...

CMD ["/app/bin/app", "port"]
# RUN go mod tidy
