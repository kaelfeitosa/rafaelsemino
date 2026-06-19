# Relatório de Análise Técnica: Migração de Ontologia (Acervo Pessoal)

## 1. Diagnóstico do Modelo Atual

O acervo encontra-se estruturado em torno das entidades `Agent`, `Work`, `Event`, `Participation` e `Record`. A análise quantitativa e qualitativa revela o seguinte cenário:

*   **Volume de Entidades:**
    *   **Records:** 159 (maior volume, predominantemente imagens).
    *   **Events:** 37.
    *   **Participations:** 30.
    *   **Works:** 26.
    *   **Agents:** 15.

*   **Uso Pleno vs. Overmodeling:**
    *   **Events:** Apenas **45.9%** dos eventos possuem descrição preenchida. Os demais 54.1% funcionam apenas como "contexto nominal" (um título e uma data), sem agregar valor narrativo ou informacional autônomo.
    *   **Participations:** Apenas **43.3%** possuem descrição própria. A maioria (56.7%) atua como tabela de ligação técnica (join table) entre Agente, Evento e Obra, sem adicionar camada semântica. Encontram-se resíduos de templates (`role: diretor | ator | ...`) em 20% das participações, indicando fadiga no preenchimento do modelo complexo.
    *   **Records:** A integridade ontológica (Record -> Participation -> Event) é respeitada em apenas **9.4%** dos casos (15 registros). A vasta maioria (100 registros) conecta-se diretamente a `Work` ou `Agent`, ignorando a mediação por `Participation`.

**Conclusão do Diagnóstico:** Há um claro *overmodeling*. O modelo CIDOC/FRBR está implementado estruturalmente, mas o uso real ("na prática") já opera em uma lógica simplificada, criando um passivo técnico de entidades vazias ou redundantes.

## 2. Análise de Dependência da Entidade Event

A entidade `Event` é o ponto mais crítico de redundância estrutural.

*   **Dependência Crítica:** Apenas **2 eventos** (`event-mestres-do-mundo-2016` e `event-hub-cultural-porto-dragao-2022`) possuem mais de uma participação vinculada.
*   **Redundância 1:1 e Órfãos:** Nos demais **35 eventos**, cerca de 25 mantêm uma relação estrita de 1:1 com uma `Participation`. Os outros 10 eventos não possuem participações vinculadas (são "órfãos" ou incompletos). Isso reforça o diagnóstico de que a entidade `Event` foi criada muitas vezes apenas para satisfazer a regra ontológica, sem utilidade prática real.
*   **Valor Estrutural:** A perda da entidade `Event` resultaria na perda de metadados centralizados (localização, data de fim) apenas para os 2 eventos agregadores que possuem múltiplas participações. Para a vasta maioria, a data e local podem ser atributos diretos da ação sem perda de informação ou duplicação excessiva.

## 3. Análise de Valor da Entidade Participation

A entidade `Participation` promete detalhar "como" um agente atuou, mas os dados mostram subutilização.

*   **Distinções de Papel:** Os papéis são muitas vezes genéricos ou não preenchidos corretamente (ex: `diretor | ator...`).
*   **Colapso em Action:**
    *   **Cenário Atual:** Entidade `Participation` (link para Evento X, Obra Y) + Entidade `Event` (Nome X, Data Z).
    *   **Cenário Simplificado (`Action`):** Entidade única com `context: "Nome do Evento X"`, `work: "Link para Obra Y"`, `role: "Papel"`, `date: "Data Z"`.
*   **Perda Semântica:** Nula para 95% do acervo. A distinção ontológica entre "O Evento" e "Minha Participação no Evento" é irrelevante quando o evento só existe no acervo porque o agente participou dele.

## 4. Análise da Entidade Record

A entidade `Record` já opera, na prática, como `Attachment`.

*   **Tipologia:** Predominância absoluta de imagens (`.jpeg`, `.png`), servindo como prova visual.
*   **Links Diretos:** O fato de a maioria dos registros apontar diretamente para `Work` ou `Agent` demonstra que a camada de "prova da participação em um evento específico" é frequentemente ignorada em favor de "imagem desta obra".
*   **Migração:** Converter `Record` em um atributo lista (`attachments` ou `gallery`) dentro de `Work` ou `Action` eliminaria 159 arquivos de metadados que hoje servem apenas para apontar para um arquivo de mídia.

## 5. Simulação de Migração

### Caso 1: Participação Genérica (Redundante)

**Modelo Atual:**
*   `participation-colaboracao-curso-bece-2023.md` (Role: genérico, Desc: null)
    *   -> Link para `event-curso-bece-2023.md` (Nome: "Curso BECE", Data: 2023)

**Modelo Simplificado (Action):**
```yaml
- type: action
  title: "Colaboração Curso BECE"
  role: "Colaborador"
  context: "Curso BECE 2023"
  date: "2023-11-22"
  # Sem necessidade de entidade separada para o evento
```
**Veredito:** Ganho de clareza, eliminação de 2 arquivos, nenhuma perda de informação.

### Caso 2: Participação Rica (Avaliador)

**Modelo Atual:**
*   `participation-avaliador-junino.md` (Role: Avaliador, Desc: "Atuação como avaliador...")
    *   -> Link para `event-ciclo-junino-ceara-2023.md`

**Modelo Simplificado (Action):**
```yaml
- type: action
  title: "Avaliador Junino"
  role: "Avaliador"
  context: "Ciclo Junino Ceará 2023"
  description: "Atuação como avaliador em festivais..."
  date: "2023"
```
**Veredito:** Integridade semântica mantida. A narrativa "Fui avaliador no Ciclo Junino" permanece intacta.

## 6. Avaliação Final

### O modelo simplificado é adequado?
**SIM.** O acervo pessoal analisado não possui complexidade ou volume que justifique a separação estrita entre Evento e Participação, nem a atomização de Registros. A estrutura atual gera fricção cognitiva (preenchimento de múltiplos arquivos para registrar um único fato) e inconsistência (campos vazios).

### Utilização do Modelo Atual
O modelo atual está **mal alinhado ao propósito**. Ele foi desenhado para preservação institucional complexa (onde múltiplos agentes participam do mesmo evento e geram múltiplos registros cruzados), mas está sendo usado para documentar a trajetória de um único agente principal.

### Riscos Reais da Simplificação
1.  **Normalização de Strings:** Ao eliminar a entidade `Event`, perde-se a "fonte única da verdade" para o nome do evento. Se o agente participou 3 vezes do "Festival de Curitiba", terá que digitar "Festival de Curitiba" 3 vezes (risco de typos). *Mitigação: Uso consistente de autocomplete ou linter.*
2.  **Perda de Metadados de Eventos Coletivos:** Se houver um evento com curadoria complexa, localizações múltiplas e organizadores detalhados, esses dados terão que ser resumidos no campo texto `context` ou perdidos.
3.  **Linkabilidade de Provas Específicas:** Se um registro prova *especificamente* a participação no dia X e não a obra em si, vinculá-lo genericamente à Obra pode perder essa especificidade temporal. *Mitigação: Legendas nos anexos.*

**Recomendação:** Migrar imediatamente. O custo de manutenção do modelo atual supera os benefícios de sua granularidade teórica.
