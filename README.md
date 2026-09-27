<p align="center">
  <img src="assets/logo.svg" alt="Diary Logo" width="540">
</p>

<p align="center">
  <strong>Gerenciador minimalista de metas diárias local-first em Go com estética e integração nativa ao Omarchy.</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-7fbbb3.svg" alt="License: MIT"></a>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.22+-83c092.svg" alt="Go 1.22+"></a>
  <a href="https://github.com/takayuki/diary/actions"><img src="https://img.shields.io/badge/CI-Passing-a7c080.svg" alt="CI Status"></a>
  <img src="https://img.shields.io/badge/Dependencies-Zero-d3c6aa.svg" alt="Zero Dependencies">
</p>

---

Seus dados ficam salvos no seu computador como arquivos Markdown (`./metas/DD-MM-YYYY.md`), sem bancos de dados externos e sem serviços na nuvem.

O projeto foi construído sobre Clean Architecture e roda com zero dependências externas: tanto o servidor quanto a interface gráfica funcionam exclusivamente com recursos nativos.

---

## Recursos principais

- **Persistência local em Markdown:** Cada dia gera um arquivo individual com checklists (`- [ ] ` e `- [x] `) compatíveis com qualquer editor de texto.
- **Zero dependências externas:** Backend puramente em Go nativo e frontend embutido no binário com `go:embed`.
- **Clean Architecture:** Camadas separadas de domínio, casos de uso, adaptadores HTTP e repositório de persistência.
- **Identidade Visual e Integração ao Omarchy:** Logotipo e interface criados sob as diretrizes geométricas e tokens do Omarchy, com detecção e sincronização dinâmica do tema ativo via `/api/theme`.
- **Interface nítida e acessível:** Tipografia JetBrains Mono, cantos nítidos, alto contraste (WCAG AAA), barra de progresso em tempo real e atalhos rápidos de teclado.
- **Navegação rápida por teclado:** Atalhos para focar no campo de texto (`/`), alternar entre datas (`←` / `→`) e retornar para hoje (`T`).
- **Operações seguras e concorrentes:** Acesso ao disco protegido por `sync.RWMutex` com tratamento contra path traversal.

---

## Estrutura do projeto

```
diary/
├── assets/                            # Identidade visual (logo horizontal e ícone SVG/PNG)
│   ├── logo.svg
│   ├── icon.svg
│   └── favicon.svg
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
│   │   ├── gui/                       # Interface Desktop Gráfica Nativa (Fyne + Omarchy Theme)
│   │   ├── handler/http/              # Servidor HTTP alternativo (modo headless/web)
│   │   ├── repository/markdown/       # Leitura e escrita concorrente de arquivos .md
│   │   └── theme/                     # Integração com Omarchy colors.toml
│   └── config/
│       └── config.go                  # Variáveis de ambiente (PORT, METAS_DIR)
├── web/
│   ├── embed.go                       # Arquivos estáticos embutidos via go:embed
│   ├── static/                        # Favicon SVG e ICO
│   └── template/
│       └── index.html                 # Interface web alternativa
├── metas/                             # Diretório onde os arquivos diários são salvos
├── .github/workflows/ci.yml           # Validação automatizada em CI
├── Makefile                           # Comandos de desenvolvimento, teste e instalação
├── CONTRIBUTING.md                    # Guia para novos contribuidores
├── LICENSE                            # Licença MIT
├── go.mod
└── main.go                            # Entrada conveniente na raiz
```

---

## Como executar e instalar

### Pré-requisitos
- Go 1.22 ou superior instalado.
- Ambiente Linux com suporte a Wayland ou X11.

### 1. Executar a Aplicação Desktop
Clone o repositório e execute diretamente:

```bash
git clone https://github.com/SEU_USUARIO/diary.git
cd diary
go run .
```

A janela nativa do Diary abrirá diretamente no seu ambiente gráfico (Hyprland / Omarchy), com navegação por mouse e teclado.

### 2. Instalar no Sistema (Lançador do Omarchy)
Para instalar o Diary como um aplicativo nativo no menu e lançadores do sistema (`rofi`, `walker` ou atalhos):

```bash
make install
```

O comando compila o binário para `~/.local/bin/diary`, copia o atalho `.desktop` para `~/.local/share/applications/` e registra os ícones em alta resolução.

### 3. Modo Web Alternativo (Opcional)
Se desejar iniciar a versão web/headless no navegador:

```bash
./bin/diary --web
# Ou acesse em http://localhost:8080
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

### `GET /api/theme`
Retorna as definições e cores do tema atualmente ativo (Omarchy ou fallback).

**Resposta de exemplo (`200 OK`):**
```json
{
  "source": "omarchy",
  "themeName": "Everforest",
  "mode": "dark",
  "colors": {
    "accent": "#7fbbb3",
    "background": "#2d353b",
    "foreground": "#d3c6aa",
    "green": "#a7c080",
    "red": "#e67e80"
  }
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
