FROM golang:1.22-alpine

RUN apk add --no-cache \
    libxml2-dev \
    libxml2-utils \
    gcc \
    musl-dev

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "main.go"]
