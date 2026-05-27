# GoApi Backend

Backend em Go organizado para autenticação e criação de recursos com PostgreSQL.

## Endpoints

- `GET /healthz`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/items` protegido por JWT
- `GET /api/v1/items` protegido por JWT

## Rodando com Docker

```bash
cd backend
docker compose up --build
```

## Rodando localmente

Instale Go 1.22+ e PostgreSQL, copie `.env.example` para `.env`, ajuste as variáveis e execute:

```bash
go mod tidy
set -a
. ./.env
set +a
go run ./cmd/api
```

## Exemplo de uso

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Vitor","email":"vitor@example.com","password":"12345678"}'

TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"vitor@example.com","password":"12345678"}' | jq -r .token)

curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Primeiro item","body":"conteudo"}'
```

## Observações para 10k+ req/min

10k requisições/minuto é cerca de 167 req/s. O servidor usa `net/http`, pool de conexões PostgreSQL, timeouts e shutdown gracioso. Em produção, rode atrás de um proxy/load balancer, ajuste `DB_MAX_CONNS` conforme o PostgreSQL, use índices, monitore p95/p99 e escale horizontalmente quando necessário.
