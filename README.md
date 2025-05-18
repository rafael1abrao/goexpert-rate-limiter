# 📈 GoExpert Rate Limiter

Rate limiter em Go baseado em IP e Token com Redis como backend e Podman para conteinerização.

---

## 🚀 Visão Geral

Este projeto é um rate limiter que funciona como middleware HTTP e permite:

- Limitação por IP (ex: 5 req/s por IP)
- Limitação por Token (ex: 10 req/s por token no header `API_KEY`)
- Prioriza o limite por token quando presente
- Backend persistente via Redis
- Configuração via `.env`
- Conteinerização com Podman e controle via `Makefile`

---

## 🧠 Como funciona

1. Toda requisição HTTP passa pelo middleware `RateLimiterMiddleware`
2. Se o cabeçalho `API_KEY` estiver presente, o limiter usa o token como chave
3. Caso contrário, o IP da requisição é usado
4. A strategy (por padrão Redis) incrementa o contador e retorna se a requisição é permitida
5. Se o limite for excedido, retorna HTTP 429

---

## ⚙️ Configuração via `.env`

```env
REDIS_HOST=redis:6379
REDIS_PASSWORD=
RATE_LIMIT_IP=5
RATE_LIMIT_TOKEN_DEFAULT=10
BLOCK_TIME_SECONDS=300
```

- `REDIS_HOST`: host do Redis acessível pelo container
- `RATE_LIMIT_IP`: limite de requisições por segundo por IP
- `RATE_LIMIT_TOKEN_DEFAULT`: limite por token (se `API_KEY` presente)
- `BLOCK_TIME_SECONDS`: tempo de bloqueio após exceder o limite

---

## 🛠 Execução

### Build e subir o ambiente:

```bash
make rebuild
```

### Parar:

```bash
make stop
```

### Testar:

```bash
curl http://localhost:8080/
```

### Testar rate limit:

```bash
for i in {1..6}; do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/; done
```

Com `API_KEY`:

```bash
for i in {1..11}; do curl -H "API_KEY: abc123" -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/; done
```

---

## ✅ Testes

Testes unitários com mocks da strategy:

```bash
go test ./internal/limiter -v
```

Inclui cenários de:
- Permissão por IP
- Bloqueio por Token
- Falha de backend (ex: Redis offline)

---

## 📦 Estrutura de Projeto

```
internal/
├── limiter/             # Core do rate limiter
│   ├── limiter.go       # Middleware principal
│   ├── service.go       # Lógica do rate limiter
│   ├── strategy.go      # Interface de strategy
│   ├── redis/           # Implementação Redis
│   └── limiter_service_test.go
```

---

## 📌 Dependências

- [Go 1.23+](https://go.dev/dl/)
- Podman Desktop
- Redis 7+
- `github.com/go-redis/redis/v8`
- `github.com/gofiber/fiber/v2`
- `github.com/stretchr/testify`

---

## 🧪 Melhorias Futuras

- Suporte a limites configuráveis por token no Redis
- Exportação de métricas Prometheus
- Endpoint `/health` com check de Redis
- Logs estruturados com `zerolog`
