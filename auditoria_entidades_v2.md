# Relatório de Auditoria de Entidades - Rodada 2

Este relatório lista as inconsistências residuais identificadas após a primeira rodada de correções e propõe ações baseadas na busca por evidências nos textos extraídos.

## 1. Obras sem Evidência Documental (Provável Lixo/Teste)

### `work-astronauta`
- **Diagnóstico:** Título "Astronauta". Metadados vazios.
- **Evidência:** Nenhuma menção encontrada nos textos extraídos (`grep` retornou vazio).
- **Ação:** **Excluir** o arquivo `acervo/entities/works/work-astronauta.md`.
- **Justificativa:** Entidade órfã sem lastro documental.

### `work-grafo`
- **Diagnóstico:** Título "Grafo". Descrição genérica de pesquisa.
- **Evidência:** Nenhuma menção encontrada nos textos extraídos (`grep` retornou vazio).
- **Ação:** **Excluir** o arquivo `acervo/entities/works/work-grafo.md`.
- **Justificativa:** Entidade órfã sem lastro documental. Possivelmente um teste de estrutura de dados.

## 2. Obras com Títulos Confirmados

### `work-300-reais`
- **Diagnóstico:** Título curto "300 Reais".
- **Evidência:** Encontrado em `portfólio_semino_antigo.docx.txt`: "300 Reais, 2014".
- **Ação:** Manter. Atualizar descrição para **"Performance teatral apresentada em 2014."**

## 3. Correção de Títulos de Participação (Informais -> Formais)

### `participation-gabriel-exu`
- **Diagnóstico:** Título informal "Gabriel Exu".
- **Evidência:** Refere-se à atuação de Gabriel França no espetáculo "Exu Não Vem Hoje".
- **Ação:** Renomear `title` para **"Atuação e Sonoplastia em Exu Não Vem Hoje"**.

### `participation-rafael-exu`
- **Diagnóstico:** Título informal "Rafael Exu".
- **Evidência:** Refere-se à atuação de Rafael Semino no espetáculo "Exu Não Vem Hoje".
- **Ação:** Renomear `title` para **"Atuação e Co-fundação de Exu Não Vem Hoje"**.

### `participation-zeis-vao`
- **Diagnóstico:** Título informal "Zeis Vao".
- **Evidência:** Zeis compôs a trilha e atuou em "Vão".
- **Ação:** Renomear `title` para **"Direção Musical e Performance em Vão"**.

### `participation-rafael-vao`
- **Diagnóstico:** Título informal "Rafael Vao".
- **Evidência:** Rafael Semino atuou e criou "Vão".
- **Ação:** Renomear `title` para **"Atuação e Criação em Vão"**.

## 4. Formalização de Participações de Ensino

### `participation-prof-hugo-sadrack`
- **Diagnóstico:** Título ausente/implícito.
- **Evidência:** `curriculo_rafael_semino.pdf.txt`: "Escola Mário Hugo Sadrak do Vale — Professor de Artes".
- **Ação:** Definir `title` como **"Professor de Artes na Escola Mário Hugo Sadrak"**.

### `participation-prof-paulo-petrola`
- **Diagnóstico:** Título ausente/implícito.
- **Evidência:** `curriculo_rafael_semino.pdf.txt`: "Escola Paulo Petrola — Professor de Artes".
- **Ação:** Definir `title` como **"Professor de Artes na Escola Paulo Petrola"**.

### `participation-prof-aceleracao`
- **Diagnóstico:** Título ausente/implícito. "Aceleração" não encontrado nos textos.
- **Justificativa:** Pode se referir a "Projeto de Aceleração da Aprendizagem" (comum em escolas públicas).
- **Ação:** Definir `title` como **"Docência em Projeto de Aceleração da Aprendizagem"**. Manter nota de revisão se não houver certeza do vínculo institucional.

## Resumo das Ações
1.  **Excluir:** `work-astronauta`, `work-grafo`.
2.  **Renomear Títulos:** `participation-gabriel-exu`, `participation-rafael-exu`, `participation-zeis-vao`, `participation-rafael-vao`.
3.  **Definir Títulos:** `participation-prof-hugo-sadrack`, `participation-prof-paulo-petrola`, `participation-prof-aceleracao`.
