# Etapa 1: Build do binário
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache \
    libxml2-dev \
    libxml2-utils \
    gcc \
    musl-dev

WORKDIR /app

# Copiar arquivos go.mod e go.sum antes para baixar as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código (ajustar o caminho se necessário)
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./xsds ./xsds

# Construir o binário Go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/xml-validator 

# Etapa 2: Imagem final
FROM alpine:3.18

RUN apk add --no-cache libxml2-utils

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/xsds /app/xsds

EXPOSE 8080

CMD ["/app/main"]
