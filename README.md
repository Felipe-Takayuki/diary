# Diary

Gerenciador minimalista de metas diárias em Go. Seus dados ficam salvos no seu computador como arquivos Markdown (`./metas/DD-MM-YYYY.md`), sem bancos de dados externos e sem serviços na nuvem.

O projeto foi construído sobre Clean Architecture e roda com zero dependências externas: tanto o servidor quanto a interface gráfica funcionam exclusivamente com recursos nativos.

---

## Recursos principais

- **Persistência local em Markdown:** Cada dia gera um arquivo individual com checklists (`- [ ] ` e `- [x] `) compatíveis com qualquer editor de texto.
- **Zero dependências externas:** Backend puramente em Go nativo e frontend embutido no binário com `go:embed`.
- **Clean Architecture:** Camadas separadas de domínio, casos de uso, adaptadores HTTP e repositório de persistência.
- **Interface nítida e acessível:** Alto contraste (WCAG AAA), alternador de tema claro e escuro, barra de progresso em tempo real e estados visuais imediatos.
- **Navegação rápida por teclado:** Atalhos para focar no campo de texto, alternar entre datas e retornar para o dia atual.
- **Operações seguras e concorrentes:** Acesso ao disco protegido por `sync.RWMutex` com tratamento contra path traversal.

---

## Estrutura do projeto

```
diary/
├── cmd/
│   └── diary/
│       └── main.go                    # Ponto de entrada padrão
├── internal/
│   ├── app/
│   │   └── app.go                     # Injeção de dependências e inicialização
│   ├── domain/                        # Entidades e regras de negócio puras
│   │   ├── date.go                    # Value Object Date (formato canônico DD-MM-YYYY)
│   │   ├── errors.go                  # Erros de domínio
│   │   ├── item.go                    # Entidade Item (validação e alternância)
│   │   ├── meta.go                    # Aggregate Root DailyGoal (progresso e métricas)
│   │   └── repository.go              # Interface GoalRepository
│   ├── usecase/                       # Casos de uso da aplicação
│   │   └── daily_goal.go              # DailyGoalUseCase (GetDailyGoals, SaveDailyGoals)
│   ├── adapter/                       # Adaptadores de entrada e saída
│   │   ├── handler/http/              # Handlers HTTP, DTOs e rotas
│   │   └── repository/markdown/       # Leitura e escrita concorrente de arquivos .md
│   └── config/
│       └── config.go                  # Variáveis de ambiente (PORT, METAS_DIR)
├── web/
│   ├── embed.go                       # Arquivos estáticos embutidos via go:embed
│   └── template/
│       └── index.html                 # Interface web (HTML, CSS e JavaScript nativo)
├── metas/                             # Diretório onde os arquivos diários são salvos
├── .github/workflows/ci.yml           # Validação automatizada em CI
├── Makefile                           # Comandos de desenvolvimento e teste
├── CONTRIBUTING.md                    # Guia para novos contribuidores
├── LICENSE                            # Licença MIT
├── go.mod
└── main.go                            # Entrada conveniente na raiz
```

---

## Como executar

### Pré-requisitos
- Go 1.22 ou superior instalado.

### 1. Iniciar o servidor
Clone o repositório e rode diretamente:

```bash
git clone https://github.com/SEU_USUARIO/diary.git
cd diary
go run .
```

O servidor iniciará no endereço [http://localhost:8080](http://localhost:8080).

### 2. Compilar um binário único
Para gerar um executável independente que contém a aplicação completa:

```bash
make build
./bin/diary
```

---

## Atalhos de teclado

| Tecla | Ação |
|---|---|
| `/` | Foca o campo de adicionar nova meta |
| `←` | Volta para o dia anterior |
| `→` | Avança para o próximo dia |
| `T` | Pula diretamente para o dia de hoje |
| `Enter` | Adiciona a meta digitada |
| `Esc` | Remove o foco do campo de texto |

---

## Variáveis de ambiente

Você pode configurar a porta e o diretório de armazenamento pelas variáveis:

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta onde o servidor HTTP escuta requisições |
| `METAS_DIR` | `./metas` | Caminho da pasta onde os arquivos `.md` são gravados |

Exemplo de execução personalizada:
```bash
PORT=3000 METAS_DIR=/caminho/meu-diario go run .
```

---

## Referência da API

### `GET /api/metas?date=DD-MM-YYYY`
Recupera a lista de metas para a data informada.

**Resposta de sucesso (`200 OK`):**
```json
{
  "date": "26-09-2026",
  "items": [
    {
      "text": "Finalizar revisão de design",
      "done": true
    },
    {
      "text": "Escrever documentação open source",
      "done": false
    }
  ]
}
```

### `POST /api/metas`
Grava as metas da data no arquivo Markdown correspondente.

**Payload:**
```json
{
  "date": "26-09-2026",
  "items": [
    {
      "text": "Finalizar revisão de design",
      "done": true
    }
  ]
}
```

---

## Testes e desenvolvimento

O repositório inclui testes unitários e de integração com verificação de concorrência.

```bash
# Executar todos os testes
make test

# Executar com detecção de concorrência (race detector)
make test-race

# Gerar relatório de cobertura
make test-cover

# Verificar formatação e análise estática
make fmt
make vet
```

---

## Licença

Distribuído sob a licença [MIT](file:///home/takayuki/Projects/diary/LICENSE). Consulte o arquivo `LICENSE` para mais detalhes.
