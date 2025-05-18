# Makefile para rodar o projeto com Podman

APP_NAME=goexpert-rate-limiter
IMAGE_NAME=$(APP_NAME):latest
CONTAINER_NAME=rate-limiter
PORT=8080
REDIS_NAME=redis
REDIS_PORT=6379
NETWORK_NAME=rate-limit-net

# Build da imagem Go
build:
	podman build -t $(IMAGE_NAME) .

# Cria a rede se não existir
create-network:
	@if ! podman network exists $(NETWORK_NAME); then \
		echo "Criando rede $(NETWORK_NAME)..."; \
		podman network create $(NETWORK_NAME); \
	else \
		echo "Rede $(NETWORK_NAME) já existe."; \
	fi

# Roda o Redis e a aplicação (modo desenvolvimento)
up: create-network
	podman run -d --name $(REDIS_NAME) --network $(NETWORK_NAME) -p $(REDIS_PORT):6379 docker.io/library/redis:7
	podman run --rm -it \
		--name $(CONTAINER_NAME) \
		--env-file .env \
		--network $(NETWORK_NAME) \
		-p $(PORT):8080 \
		$(IMAGE_NAME)

# Somente a aplicação (assumindo Redis já rodando)
run: create-network
	podman run --rm -it \
		--name $(CONTAINER_NAME) \
		--env-file .env \
		--network $(NETWORK_NAME) \
		-p $(PORT):8080 \
		$(IMAGE_NAME)

# Para Redis e aplicação
stop:
	podman stop $(CONTAINER_NAME) || true
	podman rm $(CONTAINER_NAME) || true
	podman stop $(REDIS_NAME) || true
	podman rm $(REDIS_NAME) || true

# Rebuild completo da imagem e execução
rebuild: stop build up

# Logs da aplicação
logs:
	podman logs -f $(CONTAINER_NAME)

# Testes Go
test:
	go test ./... -v

# Verifica se a VM Podman está ativa
check:
	podman machine info