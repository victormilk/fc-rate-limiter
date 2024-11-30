# Rate Limiter

## Visão Geral

Este projeto implementa um rate limiter utilizando Redis para armazenar o estado das requisições. O rate limiter é usado para limitar o número de requisições que um cliente pode fazer a um servidor em um determinado período de tempo.

## Funcionamento

O rate limiter funciona verificando se um cliente (identificado por IP ou token) já atingiu o limite de requisições permitidas. Se o limite for atingido, o cliente é bloqueado por um período de tempo. Caso contrário, a requisição é permitida e o contador de requisições é incrementado.

### Principais Componentes

- **RedisLimiter**: Implementa a lógica do rate limiter utilizando Redis.
- **Middleware**: Middleware HTTP que aplica o rate limiter às requisições.
- **Configuração**: Variáveis de ambiente para configurar o rate limiter.

## Configuração

### Variáveis de Ambiente

As seguintes variáveis de ambiente podem ser configuradas no arquivo `.env`:

- `REDIS_HOST`: Endereço do servidor Redis (ex: `localhost:6379`).
- `IP_RATE`: Limite de requisições por IP.
- `TOKEN_RATE`: Limite de requisições por token.
- `BLOCK_DURATION`: Duração do bloqueio em segundos.

Exemplo de arquivo `.env`:

```
REDIS_HOST=localhost:6379
IP_RATE=10
TOKEN_RATE=10
BLOCK_DURATION=20
```

### Dockerfile

O projeto inclui um `Dockerfile` para construir a aplicação:

```dockerfile
FROM golang:1.23.2 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/rate-limiter/main.go

FROM scratch
COPY --from=builder /app/main /app/
WORKDIR /app
CMD ["./main"]
```

## Uso

### Inicialização

Para iniciar o servidor, execute o comando:

```sh
go run cmd/rate-limiter/main.go
```

### Exemplo de Requisição

Faça uma requisição HTTP para o servidor:

```sh
curl http://localhost:8080
```

Se o limite de requisições for atingido, o servidor retornará um status `429 Too Many Requests`.

## Testes

Os testes podem ser executados utilizando o comando:

```sh
go test ./...
```
### Rodando com Docker Compose

Para iniciar os serviços, execute o comando:

```sh
docker-compose up
```

Isso irá iniciar o servidor Redis e a aplicação de rate limiter. A aplicação estará disponível em `http://localhost:8080`.