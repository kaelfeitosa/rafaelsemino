# Relatório de Auditoria de Entidades - Rodada 2

Este relatório lista as inconsistências residuais identificadas após a primeira rodada de correções. O foco está em entidades com títulos genéricos, metadados incompletos ou nomes de participação ambíguos.

## 1. Obras com Metadados Incompletos ou Títulos Genéricos

### `work-astronauta`
- **Diagnóstico:** Metadados placeholder (`status: em-andamento | concluida...`, `year: null`). Título "Astronauta" é vago.
- **Justificativa:** Parece ser um template ou esboço não preenchido. Sem ano ou descrição, não agrega valor documental.
- **Sugestão:** Verificar se existe documentação de uma obra chamada "Astronauta". Caso contrário, considerar como rascunho e excluir ou marcar para preenchimento futuro.

### `work-grafo`
- **Diagnóstico:** Título "Grafo" é extremamente genérico. Descrição diz "Pesquisa acadêmica e reflexões sobre teatro".
- **Justificativa:** Se for uma tese ou artigo, deve ter o título real da publicação. "Grafo" pode ser um erro de importação ou referência à estrutura de dados.
- **Sugestão:** Renomear para o título real da pesquisa (ex: "Título da Tese/Artigo") ou excluir se for apenas um nó de teste de estrutura.

### `work-300-reais`
- **Diagnóstico:** Título curto. Descrição aponta para "Performance encenada no FIDA (2015)".
- **Justificativa:** O título pode ser real, mas verificar se é completo.
- **Sugestão:** Manter, mas adicionar contexto se possível (ex: "Performance 300 Reais").

## 2. Participações com Títulos Ambíguos ou "De Trabalho"

### `participation-gabriel-exu`
- **Diagnóstico:** Título "Gabriel Exu". Parece um apelido ou nome de arquivo, não um título de entidade de participação formal.
- **Justificativa:** Participações devem descrever a ação (ex: "Atuação em Exu Não Vem Hoje"). O título atual mistura o agente com a obra de forma informal.
- **Sugestão:** Renomear título para **"Atuação e Sonoplastia em Exu Não Vem Hoje"**.

### `participation-rafael-exu`
- **Diagnóstico:** Título "Rafael Exu". Similar ao anterior.
- **Justificativa:** Informalidade no título da entidade.
- **Sugestão:** Renomear título para **"Atuação em Exu Não Vem Hoje"**.

### `participation-zeis-vao` & `participation-rafael-vao`
- **Diagnóstico:** Títulos "Zeis Vao" e "Rafael Vao".
- **Justificativa:** Mesmo padrão informal.
- **Sugestão:** Renomear para **"Direção Musical e Performance em Vão"** e **"Atuação e Criação em Vão"**, respectivamente.

## 3. Participações Genéricas (Professores)

### `participation-prof-aceleracao`, `participation-prof-hugo-sadrack`, etc.
- **Diagnóstico:** Títulos "N/A" ou implícitos no ID.
- **Justificativa:** IDs como `participation-prof-aceleracao` são claros para devs, mas títulos descritivos ajudam na busca e visualização.
- **Sugestão:** Adicionar títulos explícitos como **"Docência no Projeto Aceleração"**, **"Professor de Artes na Escola Hugo Sadrack"**.

## Conclusão
A estrutura está mais limpa, mas a nomenclatura de participações e algumas obras ainda carece de formalidade e precisão documental. As ações acima visam polir esses detalhes para garantir um acervo profissional e semântico.
