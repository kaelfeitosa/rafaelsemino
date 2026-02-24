# Relatório de Auditoria do Acervo (Modelo Editorial)

**Data:** 25/02/2025
**Escopo:** Estruturas, Ações, Evidências e Narrativa.
**Status Geral:** O modelo estrutural está sólido, mas há lacunas críticas em evidências (imagens faltantes) e ações duplicadas/vazias que precisam de correção manual.

---

## 1. Auditoria de Estruturas (Ontologia)
*Status: ✅ Aprovado*

*   **Tipagem:** Todos os arquivos analisados (Actions, Works, Agents) seguem a tipagem correta.
*   **IDs:** Não foram encontrados conflitos de ID globais.
*   **Campos Obrigatórios:** Todos os arquivos contêm os metadados essenciais (`id`, `title`, `performed_by`).
*   **Relações:** Os links para `performed_by` e `work_id` apontam para entidades existentes (exceto onde indicado em ações vazias).

---

## 2. Auditoria de Ações (Semântica)
*Status: ⚠️ Requer Atenção*

### 2.1 Duplicidade Crítica
Detectamos uma ação duplicada que descreve o mesmo fato com o mesmo título e período.

*   **Conflito:**
    *   `action-rafael-bolsa-ccbj` (Título: "Pesquisador em Laboratório de Pesquisa - CCBJ (Exu Não Vem Hoje)")
    *   `action-rafael-ccbj-exu` (Título: "Pesquisador em Laboratório de Pesquisa - CCBJ (Exu Não Vem Hoje)")
*   **Diagnóstico:** `action-rafael-bolsa-ccbj` está incompleta (sem descrição, sem link para obra), enquanto `action-rafael-ccbj-exu` está completa.
*   **Ação Recomendada:** **Remover** `action-rafael-bolsa-ccbj` e manter `action-rafael-ccbj-exu`.

### 2.2 Ações "Vazias" (Conteúdo Placeholder)
Várias ações possuem metadados corretos, mas o corpo do texto é genérico ou inexistente ("Detalhes específicos da participação."). Isso viola a regra de "clareza autoral".

*   **Afetados:**
    *   `action-colaboracao-curso-bece-2023` (Palestrante em Curso Protagonismo Negro)
    *   `action-prof-paulo-petrola` (Professor de Artes)
    *   `action-prof-hugo-sadrack`
    *   `action-prof-aceleracao`
*   **Ação Recomendada:** Preencher a descrição dessas ações com 1-2 parágrafos explicando o contexto e a atividade realizada, ou considerar se são relevantes para o portfólio.

---

## 3. Auditoria de Evidências (Attachments)
*Status: 🚨 Crítico*

### 3.1 Arquivos de Mídia Inexistentes (Links Quebrados)
As seguintes ações referenciam imagens que **não existem** no diretório `acervo/media/images`. Isso impede a geração correta do site.

*   **Action:** `action-farol-novo-temporada-porto-dragao-2023`
    *   Faltam: `record-temporada-hub-porto-dragao-2023-001.jpeg` (e outros 4 arquivos similares).
    *   *Nota:* Existe `event-ocupacao-hub-001.jpeg`, mas não corresponde aos nomes referenciados.
*   **Action:** `action-farol-novo-zona-de-criacao-2024`
    *   Falta: `record-zona-de-criacao-2024-001.jpeg`
*   **Work:** `work-rastros-de-exu`
    *   Falta: `record-work-rastros-de-exu-003.jpeg` (Imagens 001 e 002 existem).

**Ação Recomendada:**
1.  Verificar se os arquivos foram renomeados ou não commitados.
2.  Corrigir os nomes no arquivo Markdown ou adicionar os arquivos faltantes na pasta `media`.

### 3.2 Ausência de Evidências (Lacuna Documental)
Uma quantidade significativa de Works e Actions não possui nenhum anexo (imagem/vídeo). Embora não seja um erro técnico, é uma fraqueza narrativa ("Impacto documentado").

*   **Works sem imagem:** `work-constelacao`, `work-noite-de-alegria`.
*   **Actions sem imagem:** A maioria das ações de formação (`action-prof-*`) e pesquisa (`action-rafael-mestrado-ufc`, `action-rafael-bolsa-ccbj`).

**Ação Recomendada:** Priorizar a adição de pelo menos 1 imagem (mesmo que genérica/logo institucional) para Works e Actions principais.

---

## 4. Auditoria de Narrativa (Coerência)
*Status: ✅ Aprovado*

*   **Linha do Tempo:** Cobre de 2012 a 2024 sem hiatos inexplicáveis.
*   **Progressão:** A transição de papéis (Ator -> Criador/Pesquisador -> Diretor/Professor) é visível e bem suportada pelos metadados `my_role`.
*   **Atribuição:** O uso de `my_role` está consistente (ex: "Diretor", "Pesquisador", "Ator").

---

## Plano de Ação Imediato

1.  **Excluir** o arquivo `acervo/entities/actions/action-rafael-bolsa-ccbj.md`.
2.  **Corrigir links de imagem** em `action-farol-novo-temporada-porto-dragao-2023` e `work-rastros-de-exu`.
3.  **Escrever descrições** para as ações de docência listadas no item 2.2.
