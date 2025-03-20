FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY . .

EXPOSE ${PORT}

CMD if [ "$APP_ENV" = "development" ] ; then \
        air -c .air.toml; \
    else \
        go build -o /app/main ./cmd/api/main.go && /app/main; \
    fi

