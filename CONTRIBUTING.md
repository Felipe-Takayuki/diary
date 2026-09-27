# Guia de Contribuição

Obrigado pelo interesse em contribuir com o Diary. Este projeto foi concebido em Go com foco em simplicidade, privacidade local e zero dependências externas em tempo de execução.

## Como começar

1. Faça um fork do repositório no GitHub.
2. Clone seu fork para seu computador:
   ```bash
   git clone https://github.com/SEU_USUARIO/diary.git
   cd diary
   ```
3. Crie uma branch para sua alteração:
   ```bash
   git checkout -b minha-melhoria
   ```

## Princípios do Projeto

- **Zero dependências externas:** Usamos exclusivamente a biblioteca padrão do Go no servidor e HTML, CSS e JavaScript nativos no navegador. Evite adicionar bibliotecas externas.
- **Clean Architecture:** Mantenha as regras de negócio puras em `internal/domain`, orquestração em `internal/usecase` e detalhes de entrega em `internal/adapter`.
- **Formato Markdown direto:** A persistência grava tarefas em arquivos Markdown puros no formato `metas/DD-MM-YYYY.md`.
- **Acessibilidade visual:** Qualquer mudança na interface precisa respeitar alto contraste (WCAG AAA), navegação por teclado e compatibilidade com modo claro e escuro.

## Verificações locais

Antes de abrir um Pull Request, execute a suíte completa pelo Makefile:

```bash
# Executa testes unitários
make test

# Testa condições de corrida (race detector)
make test-race

# Formata o código e roda a análise estática
make fmt
make vet

# Valida a compilação do binário
make build
```

## Enviando seu Pull Request

1. Use mensagens de commit claras no padrão Conventional Commits (por exemplo: `feat:`, `fix:`, `docs:`, `refactor:` ou `test:`).
2. Adicione testes para novas regras de domínio ou novos casos de uso.
3. Abra seu Pull Request descrevendo a mudança de forma objetiva.
