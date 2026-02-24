# Visão Geral do Repositório e Ferramentas (ADR)

## Axioma Fundamental
Este sistema descreve a trajetória profissional de Rafael Semino, a partir do seu ponto de vista, apresentando ações que realizou, individualmente ou por meio de coletivos, associadas a obras, em determinados contextos, com evidências que sustentam essa narrativa.

### Implicações Diretas
- O sistema **não é institucional**.
- O sistema **não é neutro**.
- O sistema **não descreve tudo**.
- O sistema **assume autoria e recorte**.

## Modelo de Dados (`acervo`)

### 1. Agent (Pessoa ou Coletivo)
**Regras:**
- Rafael Semino é um Agent do tipo `person`.
- Coletivos são Agent do tipo `collective`.
- Agents não têm histórico interno nem conhecem Actions diretamente.
- Não há hierarquia entre Agents.

### 2. Work (Obra Artística)
**Regras:**
- Work não age e não tem autoria interna.
- Autoria e atuação acontecem exclusivamente na **Action**.

### 3. Action (Núcleo do Sistema)
Representa o que "Eu fiz". Toda Action vira um item visível no portfólio.

**Regras Absolutas:**
- Toda Action representa algo que **Eu fiz**.
- Não existe Action sem `my_role`.
- Não existe Action "do coletivo" sem Rafael envolvido.
- Action descreve a atuação de Rafael.

### 4. Context & Attachments (Valores Embedados)
- **Context**: Existe apenas dentro da Action (festival, mostra, curso, etc).
- **Attachment**: Evidência visual ou documental (imagem, vídeo, pdf). Não possui ID próprio e não existe fora de Action ou Work.

## Relações Válidas
- `Action.performed_by` → `Agent`
- `Action.work_id` → `Work`
*Nenhuma outra relação é permitida.*

## Regras Editoriais
1. **Se não vira card, não entra.**
2. O sistema responde "o que eu fiz", não "o que aconteceu".
3. Coletivos podem agir, mas sempre com papel explícito de Rafael.
4. Clareza narrativa > normalização de dados.

## Ferramentas CLI (`acervo/cli`)
O repositório inclui uma ferramenta CLI em Go para gerenciar o acervo.

- **Sincronização:** `go run main.go reindex` (atualiza o `db.sqlite`).
- **Verificação:** `go run main.go verify` (valida sintaxe e links).
- **Ingestão:** `go run main.go ingest create action [slug] ...`

## Decisões Arquiteturais
- O sistema é focado em **atuação editorial**, não em registros históricos exaustivos.
- Entidades são propositais e reduzidas para garantir manutenibilidade e clareza narrativa no portfólio.
