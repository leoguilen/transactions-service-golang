# Copilot instructions para este repositório — Transactions Service (Go)

Este arquivo orienta sessões futuras do Copilot/CLI a trabalhar eficientemente neste projeto.

1) Comandos úteis (build / test / run)
- Build (binário):
  make build
  -> produz bin/transactions (go build -o bin/transactions ./cmd/app)
- Rodar (local):
  make run
  -> executa `go run ./cmd/app` (requere DB_CONNECTION_STRING e HTTP_PORT)
- Rodar em dev (conexão local ao Postgres):
  make run-dev
  -> exemplo do Makefile: `DB_CONNECTION_STRING="postgres://postgres:postgres@localhost:5432/transactions_db?sslmode=disable" HTTP_PORT=5000 go run ./cmd/app`
- Rodar com Docker Compose (recomendado para ambiente local):
  make run-docker
  ./scripts/run.sh
  docker-compose up --build

- Testes (toda a suíte):
  make test
  -> executa `go test ./internal/tests/../... -v` (Makefile configura testes em internal/...)

- Executar um único teste (exemplos):
  go test ./internal/core/services -run TestCreateAccount -v
  go test ./internal/tests/integration/handlers/http -run TestCreateAccount -v
  go test ./internal/adapters/handlers -run TestName -v

Observação: os testes de integração usam testcontainers (Postgres) — ver internal/tests/integration/utils/testcontainers.go.

2) Visão geral da arquitetura (Hexagonal / Ports & Adapters)
- internal/core
  - domain: entidades e erros de domínio (Account, Transaction, validações)
  - ports: interfaces (AccountService, AccountRepository, TransactionService, TransactionRepository)
  - services: implementação dos casos de uso (NewAccountService, NewTransactionService) que dependem apenas das portas
- internal/adapters
  - handlers: adaptadores de entrega HTTP (handlers.HttpHandler). Rotas são registradas em cmd/app via ServeMux.
  - db/postgres: adaptadores de persistência que implementam as interfaces de repositório (NewAccountRepository, NewTransactionRepository)
- cmd/app
  - ponto de entrada: lê variáveis de ambiente (DB_CONNECTION_STRING obrigatório, HTTP_PORT), instancia repositórios e serviços e registra rotas
- deploy/postgres/init.sql
  - script de inicialização do banco (montado pelo docker-compose e usado pelos testes/integration container)

3) Convenções / padrões específicos do projeto
- Repositórios Postgres: construtores expostos `NewXRepository(connStr string)` abrem sql.DB e retornam um objeto que implementa a porta correspondente. Os construtores atualmente `panic` em erro de sql.Open — manter cuidado ao refatorar.
- Serviços: usar fábricas `NewAccountService` e `NewTransactionService` que retornam as interfaces em internal/core/ports; injeção de dependências por construtor é padrão.
- Erros de domínio: use os erros exportados em internal/core/domain (ex.: ErrAccountNotFound, ErrInvalidAccountID, ErrAccountAlreadyExists). Handlers fazem o mapeamento para códigos HTTP.
- Tratamento de unicidade: AccountService detecta erro de banco procurando a substring `"unique constraint"` no erro retornado para mapear para ErrAccountAlreadyExists. Alterações no driver/DB podem quebrar essa verificação.
- Registro de rotas HTTP: handlers.RegisterRoutes registra entradas no `http.ServeMux` usando strings no formato "<METHOD> <path>" (ex.: "GET /accounts/{id}"). Seguir este padrão ao adicionar rotas.
- Testes de integração: usar `internal/tests/integration/utils.SetupPostgresContainer` (testcontainers) que monta deploy/postgres/init.sql. Reaproveitar esse helper para isolar dependências.
- SQL init: mantenha deploy/postgres/init.sql consistente com queries em adaptadores (nomes de colunas, tipos e operation types iniciais).

4) Padrões de implementação comuns
- Nomeclatura: tipos de repositório usam sufixo Repository (AccountRepository, TransactionRepository).
- Serviços expõem apenas a interface definida em ports — retornar a interface facilita mocking nos testes.
- Ao abrir conexões DB em adaptadores, sempre use conn, defer conn.Close() em operações que usam `db.Conn(ctx)`.
- Ao adicionar validações de domínio, colocar regras em internal/core/domain para que serviços e adaptadores compartilhem.

5) Testes e ambiente local
- go.mod especifica `go 1.25.0` — usar essa versão para builds e CI.
- Testes de integração iniciam containers via testcontainers; tempo de startup pode variar. Ajuste timeouts em utils se necessário.
- Para executar testes de integração isolados, use o pacote de testes em internal/tests/integration/handlers/http e seus helpers.

6) Arquivos relevantes para automações e CI
- docker-compose.yaml: monta deploy/postgres/init.sql no container Postgres e define DB_CONNECTION_STRING de exemplo.
- scripts/run.sh: wrapper para docker-compose up --build.
- Makefile: targets build / test / run / run-dev / run-docker.

7) Sugestões ao criar PRs que alteram infra/domínio
- Mudanças em schemas SQL devem atualizar deploy/postgres/init.sql e testes de integração.
- Alterações em ports/interfaces exigem atualização coordenada de adaptadores e serviços e execução de `go test ./...`.

---
Resumo: este documento cobre os comandos de build/test/run, a arquitetura hexagonal, e convenções específicas (construtores NewXRepository, mapeamento de erros, formato de rotas, uso de testcontainers).