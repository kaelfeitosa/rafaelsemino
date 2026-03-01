# Guia do Acervo e Modelo de Dados

## Axioma Fundamental
Este sistema descreve a trajetória profissional de Rafael Semino, a partir do seu ponto de vista, apresentando ações que realizou, individualmente ou por meio de coletivos, associadas a obras, em determinados contextos, com evidências que sustentam essa narrativa.

### Implicações Diretas
- O sistema **não é institucional**.
- O sistema **não é neutro**.
- O sistema **não descreve tudo**.
- O sistema **assume autoria e recorte**.

## Modelo de Dados (`acervo`)

### 1. Agent (Pessoa, Coletivo ou Instituição)
**Regras:**
- Rafael Semino é um Agent do tipo `person`.
- Coletivos ou instituições são Agent do tipo `collective`.
- Agents não têm histórico interno linear, eles apenas existem como entidades nominais para cruzamento de dados (`collaborators`).

### 2. Work (Obra Artística ou Projeto) - O Núcleo do Sistema
**Regras:**
- É a espinha dorsal narrativa do Portfólio.
- Representa um espetáculo, um filme, uma pesquisa acadêmica, uma oficina ministrada ou um evento cultural.
- Pode ser algo criado por Rafael ou um projeto autônomo do qual Rafael participou.
- Em projetos autônomos de terceiros, o papel de Rafael é sinalizado diretamente via campo `role`.

**Campos Principais:**
- **Medium** (`medium`): Formato da obra (teatro, audiovisual, pesquisa, cultura_popular, jogos_digitais, etc.). 
- **Collaborators** (`collaborators`): Lista simples de Agents relacionados àquela Obra.

### 3. Occurrence (Linha do Tempo Aninhada)
A antiga entidade solta `Action` foi abolida. Agora, os eventos temporais de Rafael vivem exclusivamente *dentro* das restrições de uma Obra (`Work`), chamados de `Occurrences`.

**Campos Principais:**
- **Type** (`type`): Tipo de ocorrência (apresentacao, residencia, lancamento, oficina, premio, exposicao).
- **Date** (`start_date` e `end_date`): Balizadores temporais lógicos para construção da interface gráfica (Timeline).
- **Context** (`context`): O nome do festival, edital ou temporada.

### 4. Attachments (Evidências e Clipping)
- **Attachment**: Ficam anexados aos Works. Possuem um rígido esquema de `category`:
  - `documentation`: Fotografia de Cena / Registro puro.
  - `poster`: Material de Divulgação (Design vertical/quadrado).
  - `clipping`: Matérias de Imprensa e prêmios.
  - `program`: Fichas Técnicas documentais em PDF/Imagem.
  - `technical`: Riders técnicos, mapas de palco, etc.

## Relações Válidas (Grafo)
- `Work.collaborators` → `Agent`
*Nenhuma outra relação cruzada arbitrária é permitida pela engine para evitar teias de aranha.*

## Regras Editoriais
1. **Se não vira card no Frontend, não entra.**
2. Coletivos podem figurar como colaboradores, mas a narrativa é pautada na lente de criação/atuação de Rafael.
3. Clareza narrativa > Modelos engessados de arquivologia tradicional.

## Ferramentas CLI (`acervo/cli`)
O repositório inclui uma ferramenta CLI em Go para gerenciar o acervo.

- **Sincronização:** `go run main.go reindex` (atualiza o `db.sqlite`).
- **Verificação:** `go run main.go verify` (valida sintaxe e links).
- **Ingestão:** `go run main.go ingest create action [slug] ...`

## Decisões Arquiteturais
- O sistema é focado em **atuação editorial**, não em registros históricos exaustivos.
- Entidades são propositais e reduzidas para garantir manutenibilidade e clareza narrativa no portfólio.
