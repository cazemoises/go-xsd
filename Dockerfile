FROM golang:1.22

RUN apt-get update && apt-get install -y \
    libxml2 \
    libxml2-dev \
    libxml2-utils

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "main.go"]
