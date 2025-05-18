# Etapa 1: builder com Go 1.23
FROM golang:1.23-alpine AS builder

# Instala ferramentas necessárias
RUN apk add --no-cache git

# Define o diretório de trabalho
WORKDIR /app

# Copia arquivos de dependência primeiro (para cache eficiente)
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila a aplicação para produção (binário estático)
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter ./cmd/server

# Etapa 2: imagem final enxuta
FROM alpine:latest

# Instala certificados raiz para conexões TLS
RUN apk --no-cache add ca-certificates

# Define diretório padrão
WORKDIR /root/

# Copia o binário da etapa anterior
COPY --from=builder /app/rate-limiter .

# Copia o arquivo .env (se desejar embutir)
COPY .env .

# Exposição da porta da API
EXPOSE 8080

# Comando padrão de execução
ENTRYPOINT ["./rate-limiter"]