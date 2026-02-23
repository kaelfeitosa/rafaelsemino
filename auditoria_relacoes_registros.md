# Relatório de Auditoria de Relações e Registros

Este relatório detalha problemas estruturais (campos ausentes, relações proibidas) identificados pelo script de auditoria profunda e validados com os textos extraídos.

## 1. Campos Obrigatórios Ausentes

### Obras
- **`work-trapo-preto`**
    - **Problema:** Campo `created_by` ausente.
    - **Ação:** Adicionar `created_by: "[[agent-rafael-semino]]"` (assumindo autoria, ou marcar para verificação).

### Eventos (Falta de Localização)
Os seguintes eventos não têm o campo `location`. Devem ser preenchidos com a cidade/local onde ocorreram.
- **`event-laboratorio-teatro-porto-iracema`**: Fortaleza - CE (Porto Iracema das Artes).
- **`event-apresentacao-centro-cultural-nordeste`**: Fortaleza - CE (CCBNB).
- **`event-ocupacao-hub`**: Fortaleza - CE (Hub Cultural Porto Dragão).
- **`event-projeto-abarca-itapipoca`**: Itapipoca - CE.
- **`event-premio-amarracoes-esteticas`**: Fortaleza - CE (Porto Iracema das Artes).
- **`event-aniversario-farol-novo`**: Fortaleza - CE.

### Participações (Falta de Role ou Link)
- **`participation-pesquisador-ifce`**
    - **Problema:** Falta `event` ou `work`, e falta `role`.
    - **Ação:** Vincular a `event-licenciatura-ifce` (ou similar) e definir `role: Pesquisador`.
- **`participation-rafael-ensino-angola-2018`**
    - **Problema:** Falta `role`.
    - **Ação:** Definir `role: Professor e Produtor` (conforme descrição).
- **`participation-avaliador-junino`**
    - **Problema:** Falta `role`.
    - **Ação:** Definir `role: Avaliador`.
- **`participation-avaliador-ciclo-carnavalesco`**
    - **Problema:** Falta `role`.
    - **Ação:** Definir `role: Avaliador`.

## 2. Relações Suspeitas (Records ligados a Agents)

Registros devem provar uma *Participação* ou documentar um *Evento/Obra*, não apenas uma pessoa.

- **`record-record-rebordose-imdb`**
    - **Problema:** Link direto para `agent-rafael-semino`.
    - **Ação:** Deve linkar para `work-rebordose` (se existir) ou criar uma participação de Rafael em Rebordose. (Rebordose aparece no portfolio como obra).
- **`record-record-habite-se-livro`**
    - **Problema:** Link direto para `agent-rafael-semino`.
    - **Ação:** Se for um livro de autoria, deve haver um `work-habite-se`. Se for participação em antologia, criar Participation.
- **`record-record-mestres-do-mundo-catalogo`**
    - **Problema:** Link direto para `agent-rafael-semino`.
    - **Ação:** Deve linkar para a participação de Rafael no evento `event-mestres-do-mundo-*` ou para o evento em si.

## 3. Entidades "Trapo Preto"
- **Diagnóstico:** `work-trapo-preto` existe mas não foi encontrado nos textos.
- **Ação:** Verificar se é rascunho. Se não houver evidência, considerar exclusão ou marcar como "Em Progresso" se for recente.

## Ações Recomendadas
1.  **Preencher Locations:** Atualizar eventos com locais conhecidos (Fortaleza/Itapipoca).
2.  **Corrigir Roles:** Adicionar papeis explícitos nas participações.
3.  **Relinkar Records:** Mover links de Agent para Work/Participation/Event apropriados.
4.  **Corrigir `work-trapo-preto`:** Adicionar autor ou excluir.
