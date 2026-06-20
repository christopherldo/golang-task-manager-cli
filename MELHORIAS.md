# Melhorias do Projeto — todo-cli

Roadmap de melhorias, em ordem de prioridade. Marque conforme for concluindo.

> Estado atual: bom para um primeiro projeto em Go — estrutura `cmd/` + `internal/` correta,
> separação repository/db/cli/api, `sync.RWMutex`, error wrapping com `%w`, DTOs e testes.

---

## 1. Bug de concorrência (prioridade máxima)

Em `internal/todo/repository/task.go`, o padrão abaixo se repete em `Append`, `Update`,
`MarkAsDone` e `Delete`:

```go
cacheMutex.Lock()
cachedTasks[idx] = taskToBeUpdated
tasksToWrite := cachedTasks   // só copia o header do slice, NÃO os dados
cacheMutex.Unlock()           // lock liberado aqui

err := db.WriteDatabase(tasksToWrite)  // lê o array enquanto outra goroutine pode mutá-lo
```

- `tasksToWrite := cachedTasks` copia só o cabeçalho do slice — o array por baixo é o mesmo.
- O lock é liberado **antes** do `WriteDatabase`, então no modo `api` (requisições HTTP
  concorrentes) outra goroutine pode mutar o array durante o `json.Marshal`. Isso é um *data race*.

**Como corrigir (escolher uma):**
- [ ] Simples: segurar o lock através do `WriteDatabase` (mover `defer cacheMutex.Unlock()` para o fim).
- [x] Melhor para concorrência: copiar de verdade antes de soltar o lock — `tasksToWrite := slices.Clone(cachedTasks)`.
- [ ] Escrever um teste que dispara N goroutines no repositório e rodar com `go test -race`
      para provar a corrida (os testes atuais são sequenciais e não a pegam).

---

## 2. Estado global → struct `Repository` (idiomático)

Hoje `cachedTasks`, `cachedLastTaskId` e `DbUrl` são variáveis de pacote. Encapsular num struct
com injeção de dependência é o padrão idiomático.

```go
type Repository struct {
    mu     sync.RWMutex
    tasks  []models.Task
    lastID int
    store  Store // interface: ReadDatabase/WriteDatabase
}

func New(store Store) *Repository { ... }
func (r *Repository) Add(desc string) (models.Task, error) { ... }
```

- [x] Encapsular estado no struct `Repository`.
- [x] Definir uma interface `Store` para o backend de persistência (permite mock e troca por SQLite depois).
- [x] Remover o `resetState()` global dos testes (cada teste cria sua própria instância).

---

## 3. Convenção de mensagens de erro

Em Go, strings de erro devem ser **minúsculas e sem pontuação final** (são concatenadas com `%w`).
`staticcheck` pega isso (o `go vet` não).

```go
// evite
fmt.Errorf("Error reading file: %w", err)
// prefira
fmt.Errorf("reading database file: %w", err)
```

- [ ] Padronizar todas as strings de erro (minúsculas, sem ponto final).
- [ ] Unificar o idioma das mensagens (hoje há mistura PT/EN, ex: "Task não encontrado" vs "Error reading file").

---

## 4. Geração de ID frágil

`GetLastTaskId()` assume que a última task do slice tem o maior ID, e `cachedLastTaskId` não é
atualizado em deletes — pode reutilizar IDs (adicionar 5, deletar, adicionar de novo → ID 5 reaparece).

- [ ] Usar um contador monotônico no struct (`lastID++` ao criar) ou calcular `max(IDs)` de verdade.
- [ ] Garantir que IDs nunca sejam reutilizados.

---

## 5. Hardening do servidor HTTP

`http.ListenAndServe(":8080", mux)` não tem timeouts (vulnerável a slowloris) nem graceful shutdown.

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
}
// + signal.NotifyContext(os.Interrupt) e srv.Shutdown(ctx)
```

- [ ] Adicionar `ReadTimeout` / `WriteTimeout` ao `http.Server`.
- [ ] Implementar graceful shutdown com `signal.NotifyContext` + `srv.Shutdown(ctx)`.
- [ ] `httpCreateTask` retornar a task criada em JSON (com o ID) no corpo, em vez de `201` vazio.

---

## 6. Cobertura de testes

Os testes de repositório estão bons. Faltam:

- [ ] Handlers HTTP com `net/http/httptest` (`httptest.NewRequest` + `ResponseRecorder`).
- [ ] Converter para **table-driven tests** (padrão idiomático para cobrir vários casos).
- [ ] Repensar `TestGetLastTaskId` (hoje chama `AppendTaskToDatabase(models.Task{})` com ID 0 e
      passa por coincidência da lógica atual) — revisar junto com o item 4.

---

## 7. Ferramentas e qualidade

- [ ] `Makefile` com alvos `build`, `test`, `lint`, `run`.
- [ ] Adicionar `staticcheck` ou `golangci-lint` (pegam muito mais que o `go vet`).
- [ ] GitHub Actions rodando `go test -race ./...` e o linter em cada push.
- [ ] Trocar `fmt.Println` por `log/slog` para mensagens internas/erros do servidor.

---

## 8. Detalhes menores

- [ ] O `help` documenta o comando `update`, mas o `switch` usa `edit` — corrigir a doc.
- [ ] `cliFuncEdit` não imprime mensagem de sucesso (todos os outros imprimem).
- [ ] Validar descrição vazia (CLI e API aceitam task sem texto).
- [ ] Extrair o código repetido entre `Completing/Editing/Deleting` em `interactive.go`
      numa função comum ("selecionar ID → executar → confirmar").
