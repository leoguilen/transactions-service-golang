## Transactions Service (Go)

Este repositório implementa um serviço de transações bancárias simples seguindo a Arquitetura Hexagonal (Ports & Adapters / Clean Architecture). O objetivo é demonstrar separação de responsabilidades entre domínio, portas (interfaces), adaptadores de infraestrutura (Postgres) e adaptadores de entrega (HTTP).

### Tecnologias
- Linguagem: Go
- Banco de dados: PostgreSQL (init SQL em deploy/postgres/init.sql)
- Testes: go test (inclui testes unitários e de integração)
- Docker / docker-compose para ambiente local

### Visão geral da arquitetura (Hexagonal)
- internal/core: domínio (entities, erros) e portas (interfaces de serviço/repositório)
- internal/adapters: adaptadores de entrega (HTTP handlers)
- internal/adapters/db/postgres: adaptadores de persistência (implementação dos repositórios)
- cmd/app: bootstrap da aplicação (injeção de dependências e servidor HTTP)

### Estrutura principal
- cmd/app/main.go          -> ponto de entrada
- internal/core/...        -> domínio, portas e serviços
- internal/adapters/...    -> handlers HTTP e repositórios Postgres
- deploy/postgres/init.sql -> script de criação de tabelas e dados iniciais
- internal/tests/...       -> testes unitários e de integração

### Configuração / Variáveis de ambiente
- DB_CONNECTION_STRING: string de conexão com Postgres (obrigatório)
- HTTP_PORT: porta HTTP (ex.: 8080)

### Endpoints HTTP
#### 1. Criar conta
- POST /accounts
- Body: { "document_number": "<string>" }
- Responses:
  - 201 Created: { "id": int, "document_number": string, "created_at": string }
  - 400 Bad Request: entrada inválida
  - 409 Conflict: conta já existe
  - Location header: /accounts/{id}

#### 2. Buscar conta por ID
- GET /accounts/{id}
- Responses:
  - 200 OK: { "id": int, "document_number": string, "created_at": string }
  - 400 Bad Request: id inválido
  - 404 Not Found: conta não encontrada

#### 3. Criar transação
- POST /transactions
- Body: { "account_id": int, "operation_type_id": int, "amount": number }
- Responses:
  - 201 Created: { "id": int, "account_id": int, "operation_type_id": int, "amount": number, "event_date": string }
  - 400 Bad Request: dados inválidos (conta inexistente, operação inválida, amount inválido)

### Observações sobre Operation Types
- O arquivo deploy/postgres/init.sql insere tipos de operação iniciais (p.ex. Normal Purchase, Withdrawal, Credit Voucher). Use os IDs correspondentes ao criar transações.

### Como rodar localmente
1. Com Docker Compose (recomendado):
   ./scripts/run.sh
   -- ou --
   docker-compose up --build
   (o serviço Postgres montará deploy/postgres/init.sql)

2. Rodar local sem Docker (Postgres já disponível):
   export DB_CONNECTION_STRING="postgresql://user:pass@host:5432/dbname?sslmode=disable"
   export HTTP_PORT=8080
   go run ./cmd/app

### Exemplos com curl
- Criar conta:
```bash
  curl -v -X POST http://localhost:8080/accounts \
    -H "Content-Type: application/json" \
    -d '{"document_number":"12345678900"}'
```

- Criar transação:
```bash
  curl -v -X POST http://localhost:8080/transactions \
    -H "Content-Type: application/json" \
    -d '{"account_id":1,"operation_type_id":1,"amount":100.5}'
```

### Testes
- Executar toda a suíte de testes:
  go test ./...

- Executar testes de integração específicos (ex.: usam Testcontainers/Docker):
  go test ./internal/tests/integration/handlers/http -run TestCreateAccount -v

### Observações finais
- Código organizado para facilitar troca de adaptadores (ex.: trocar Postgres por outro DB ou adicionar cache).