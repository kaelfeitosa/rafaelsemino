# Relatório da Primeira Auditoria com IA (Camada 0-6)

Este documento registra os resultados da primeira auditoria automatizada do acervo, realizada em conformidade com os princípios de "não apagar, mas classificar".

## CAMADA 1 — SANIDADE DO MODELO

### 1.1 — Tipagem correta
- **Status**: ✅ Nenhuma inconsistência de tipo encontrada entre o diretório e o metadado `type`.

### 1.2 — Relações proibidas
Identificadas violações onde `Record` está ligado diretamente a `Work`. O modelo exige que registros documentem a materialização da obra (Participação ou Evento), não a ideia (Obra).

- **Violação**: `record-record-rebordose-imdb` → `work-rebordose`
  - **Diagnóstico**: O link direto para a obra ignora a participação (ex: Atuação, Direção) ou o evento de lançamento.
  - **Ação Recomendada**: Criar `participation-rafael-rebordose` e relinkar o registro.

- **Violação**: `record-work-rastros-de-exu-003` → `work-rastros-de-exu`
  - **Diagnóstico**: Mesmo caso acima.
  - **Ação Recomendada**: Criar `participation-rafael-rastros-de-exu` e relinkar.

## CAMADA 2 — REDUNDÂNCIA E RUÍDO

### 2.1 — Participações duplicadas
- **Status**: ✅ Nenhuma duplicação óbvia (mesmo Agente + Evento + Papel) detectada automaticamente.

### 2.2 — Eventos de baixo valor documental
- **Status**: ✅ Nenhum evento com palavras-chave de baixo valor (reunião, ensaio interno) detectado.

## CAMADA 3 — VALOR DO ACERVO

### 3.1 — Obras centrais (Top Citadas)
As seguintes obras aparecem como hubs de conexões, indicando centralidade na carreira:
1.  **Exu Não Vem Hoje** (4 referências)
2.  **Vão** (3 referências)
3.  **Rebordose** (2 referências)
4.  **Habite-se** (1 referência)
5.  **Rastros de Exu** (1 referência)

### 3.2 — Participações de alto impacto
Participações com maior volume de conexões (provavelmente registros):
1.  *Farol Novo Temporada Porto Dragão 2023* (5 refs)
2.  *Rafael Vão CCBNB 2023* (2 refs)

## CAMADA 4 — CLIPPING E REGISTROS

### 4.1 — Clipping relevante
- **Alerta**: Nenhum registro possui a tag `clipping`. Recomenda-se revisar registros de imprensa e adicionar esta tag para facilitar a geração de press-kits.

### 4.2 — Lacunas de Evidência
Identificadas **19 participações** sem nenhum registro documental direto vinculado. Exemplos prioritários para busca de arquivo:
- *Atuação e Sonoplastia em Exu Não Vem Hoje*
- *Direção Musical e Performance em Vão*
- *Atuação e Criação em Vão*
- *Professor de Artes* (várias escolas)

## CAMADA 5 — COERÊNCIA DA CARREIRA

### 5.1 — Linha do Tempo
- **2015-2021**: Baixa densidade de eventos registrados (1-3 por ano). Pode indicar sub-notificação ou fase de formação menos documentada.
- **2022-2024**: Crescimento exponencial de registros e eventos (4 a 8 por ano).

### 5.2 — O que falta
- Participações explícitas para obras audiovisuais (*Rebordose*).
- Documentação da fase inicial (2015-2021).
- Classificação de clipping.

## CAMADA 6 — PLANO DE CORREÇÃO (Ações Priorizadas)

1.  **Correção Estrutural (Urgente)**:
    - Criar entidade `participation` para *Rebordose* (papel no filme).
    - Criar entidade `participation` para *Rastros de Exu*.
    - Mover os links dos Records correspondentes das Obras para essas Participações.

2.  **Enriquecimento (Médio Prazo)**:
    - Revisar as 19 participações sem registro e digitalizar/vincular pelo menos 1 evidência para as mais relevantes (Exu, Vão).
    - Adicionar tag `clipping` em registros de imprensa.

3.  **Curadoria (Longo Prazo)**:
    - Investigar a lacuna 2015-2021: adicionar eventos chave se existirem, ou assumir como período de hiato/formação.

---
*Relatório gerado automaticamente em auditoria nível 1.*
