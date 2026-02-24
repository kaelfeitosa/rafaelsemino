# AGENTS.md

## Estrutura do Repositório

Este repositório contém o portfólio profissional de Rafael Semino.

### Diretórios Principais

- `acervo/`: Fonte da verdade dos dados (banco de dados editorial).
- `frontend/`: Código fonte do site estático (HTML/CSS/JS).
- `_materials/`: Documentos brutos (PDFs, DOCX) usados para extração de texto.

### Modelo de Dados (`acervo`)

O acervo segue um modelo editorial focado em **Ações** (Actions), **Obras** (Works) e **Agentes** (Agents).

1.  **Actions (`acervo/entities/actions`)**:
    -   Representam fatos concretos da trajetória (e.g., "Atuação em Vão", "Professor no Projeto Abarca").
    -   Substituem os antigos conceitos de `Participation` e `Event`.
    -   Possuem campos obrigatórios: `performed_by` (quem fez), `my_role` (papel editorial), `context` (onde/quando).

2.  **Works (`acervo/entities/works`)**:
    -   Representam as obras artísticas em si (e.g., espetáculo "Vão", livro "Contos de Exu").
    -   Não possuem autoria interna; a autoria é definida pelas Actions que criaram a obra.

3.  **Agents (`acervo/entities/agents`)**:
    -   Pessoas ou Coletivos.
    -   Campo `kind` define se é `person` ou `collective`.

### Ferramentas CLI (`acervo/cli`)

O repositório inclui uma ferramenta CLI em Go para gerenciar o acervo.

- **Reindexar Banco de Dados:**
  ```bash
  cd acervo/cli
  go run main.go reindex
  ```
  Isso lê os arquivos Markdown em `acervo/entities` e recria o arquivo `acervo/db.sqlite`.

- **Verificar Integridade:**
  ```bash
  cd acervo/cli
  go run main.go verify
  ```
  Isso executa validações sintáticas (campos obrigatórios) e semânticas (links quebrados).

- **Criar Nova Entidade:**
  ```bash
  cd acervo/cli
  go run main.go ingest create action [slug] title="Título" performed_by="[[agent-rafael-semino]]" ...
  ```

### Regras de Manutenção

1.  **Nunca edite o `db.sqlite` manualmente.** Ele é derivado dos arquivos Markdown.
2.  **Imagens:** Use o componente `<focus-image>` no frontend. O CLI possui o comando `set-focus` para ajustar o ponto de foco (XMP).
3.  **Migração 2024:** O sistema foi migrado de um modelo `Event/Participation` para `Action/Work`. Não crie pastas antigas (`events`, `participations`, `records`).
