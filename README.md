# labs-auction-goexpert

API de leilão em Go com MongoDB.

## Pré-requisitos

- Go `1.20+`
- Docker e Docker Compose
- Git

## Variáveis de ambiente

O projeto usa o arquivo `cmd/auction/.env`.

Arquivos de ambiente dev disponíveis:

- `cmd/auction/.env.dev.docker` (para rodar com Docker Compose)
- `cmd/auction/.env.dev.local` (para rodar API localmente)
- `cmd/auction/.env.example` (modelo base)

Antes de rodar, copie o perfil desejado para `cmd/auction/.env`:

```bash
cp cmd/auction/.env.dev.docker cmd/auction/.env
```

ou

```bash
cp cmd/auction/.env.dev.local cmd/auction/.env
```

Exemplo atual:

```env
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s

MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

### Observações importantes

- `AUCTION_INTERVAL` define o tempo de duração do leilão (ex.: `20s`, `2m`, `5m`).
- `BATCH_INSERT_INTERVAL` e `MAX_BATCH_SIZE` controlam o processamento em lote de lances.
- A URL `mongodb://...@mongodb:27017` funciona no Docker Compose (hostname `mongodb`).
- Para rodar local sem Docker Compose, use `localhost` no `MONGODB_URL`.

Exemplo local:

```env
MONGODB_URL=mongodb://admin:admin@localhost:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

## Rodando em ambiente dev com Docker Compose

Na raiz do projeto:

```bash
cp cmd/auction/.env.dev.docker cmd/auction/.env
```

Depois execute:

```bash
docker compose up --build
```

A aplicação ficará disponível em:

- `http://localhost:8080`

Para encerrar:

```bash
docker compose down
```

Para remover volumes (resetar dados do Mongo):

```bash
docker compose down -v
```

## Rodando local (app em Go + Mongo separado)

### 1) Suba o MongoDB

Opção com Docker:

```bash
docker run --name auction-mongo \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=admin \
  -p 27017:27017 \
  -d mongo:latest
```

### 2) Ajuste o `.env`

Use o perfil local:

```bash
cp cmd/auction/.env.dev.local cmd/auction/.env
```

### 3) Execute a API

```bash
go run cmd/auction/main.go
```

API em:

- `http://localhost:8080`

## Testes

Rodar todos os testes:

```bash
go test ./...
```

Rodar apenas pacote de leilão:

```bash
go test ./internal/infra/database/auction -v
```

## Rotas principais

- `POST /auction`
- `GET /auction`
- `GET /auction/:auctionId`
- `GET /auction/winner/:auctionId`
- `POST /bid`
- `GET /bid/:auctionId`
- `GET /user/:userId`
