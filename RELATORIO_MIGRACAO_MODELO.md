# Relatório Técnico: Migração de Modelo Ontológico (Acervo Pessoal)

**Data:** 24 de Maio de 2024
**Autor:** Information Architect (Simulado)
**Contexto:** Migração de modelo institucional (CIDOC/FRBR) para modelo editorial de portfólio.

---

## 1. Diagnóstico do Modelo Antigo

O modelo atual, inspirado em CIDOC CRM e FRBR, apresenta uma granularidade excessiva para um portfólio pessoal, gerando "ruído ontológico" onde a manutenção da estrutura supera o valor da informação armazenada.

### 1.1 Overmodeling e Infraestrutura Vazia
A análise dos arquivos existentes revelou que entidades do tipo `Event` e `Participation` frequentemente atuam apenas como "nós de ligação" sem conteúdo próprio relevante.

*   **Events como Placeholders:** Arquivos como `event-laboratorio-teatro-porto-iracema.md` contêm apenas `name`, `location` e `date`. Não há descrição, curadoria ou narrativa. Eles existem apenas para que uma `Participation` possa apontar para algo. Isso obriga a criação de dois arquivos (Participação + Evento) para registrar um único fato biográfico.
*   **Participations Redundantes:** Arquivos como `participation-rafael-mestrado-ufc.md` ligam o Agente ao Evento com um papel (`role`). O conteúdo textual é mínimo ou inexistente ("Detalhes específicos da participação"). A informação real ("Fiz mestrado na UFC") está diluída em três nós (Agent, Participation, Event).

### 1.2 O Problema dos Records
A entidade `Record` (ex: `record-work-vao-005.md`) funciona como um *wrapper* metadata para arquivos de mídia.
*   **Diagnóstico:** Cada foto ou documento exige um arquivo Markdown individual apenas para dizer "este arquivo existe e pertence a tal obra".
*   **Impacto:** Para um portfólio com 100 fotos, o modelo atual exige 100 arquivos de `Record`, poluindo o sistema de arquivos e o grafo de conhecimento com nós que não possuem valor semântico autônomo.

---

## 2. Análise de Migração de Entidades

A proposta de migração para o trio `Agent` - `Work` - `Action` (com `Context` e `Attachment` embedados) representa uma mudança de paradigma: de "preservação de dados" para "narrativa de trajetória".

| Entidade Antiga | Novo Modelo | Análise de Impacto |
| :--- | :--- | :--- |
| **Event** | **Context (String)** | **Perda:** Identidade única do evento (ID). Não é mais possível listar "todos os participantes do evento X" automaticamente se os nomes variarem.<br>**Ganho:** Eliminação de centenas de arquivos *stub*. Agilidade na criação de novos registros. |
| **Participation** | **Action** | **Ganho Estrutural:** A `Action` torna-se a unidade atômica da biografia. Ela carrega o "O quê" (Título), "Como" (Role) e "Onde" (Context) em um único objeto.<br>**Impacto Narrativo:** Foco na *agência* do artista ("Eu fiz") em vez da *ocorrência* ("Eu estava lá"). |
| **Record** | **Attachment** | **Impacto Semântico:** Deixa de ser um "objeto arquivístico" para ser "evidência". A imagem passa a ser subordinada à Ação ou Obra, o que reflete melhor a natureza de um portfólio (a foto serve para mostrar a obra). |
| **Agent** | **Agent** | **Continuidade:** Mantém-se a estrutura, expandindo o escopo para incluir Coletivos, o que é essencial para a prática teatral. |

---

## 3. Avaliação de Coletivos como Agents

A introdução de Coletivos (ex: grupos de teatro, instituições de ensino) como `Agent` resolve a ambiguidade da autoria compartilhada.

*   **Clareza:** Permite distinguir "Trabalho Autoral" de "Trabalho em Grupo".
*   **Modelagem de `Action`:**
    *   *Cenário A:* O Coletivo realiza uma Action (ex: "Publicação de Manifesto"). Eu tenho uma relação com o Coletivo.
    *   *Cenário B:* Eu realizo uma Action (ex: "Atuação") dentro do contexto do Coletivo.
*   **Risco Institucional:** Ao tratar instituições (ex: UFBA, Porto Iracema) como Agentes, corre-se o risco de querer modelar a história da instituição.
*   **Mitigação:** A regra "Toda Action representa algo que EU fiz" deve ser soberana. O Agente "Porto Iracema" só deve existir se for para ser referenciado como *contexto* ou *parceiro* em uma Ação minha.

---

## 4. Simulação de Casos Reais

Abaixo, a comparação direta entre a representação atual (baseada nos arquivos analisados) e a proposta.

### Caso 1: Formação Acadêmica
*Fato:* Mestrado em Artes na UFC (2024).

**Modelo Antigo (3 Arquivos):**
1.  `event-mestrado-ufc.md` (Nome: Mestrado UFC, Data: 2024)
2.  `participation-rafael-mestrado-ufc.md` (Role: Mestrando, Link: event-mestrado-ufc)
3.  `agent-rafael-semino.md` (Linkado pela participation)

**Modelo Novo (1 Arquivo):**
```yaml
# action-mestrado-ufc.md
type: action
title: "Mestrado em Artes"
role: "Pesquisador"
context: "UFC - Universidade Federal do Ceará (2024)"
date: "2024-01-01"
agent: [[agent-rafael-semino]]
category: "formacao"
```
*Análise:* Redução de complexidade de 3:1. A integridade narrativa é mantida, pois o contexto "UFC" está preservado.

---

### Caso 2: Obra Teatral
*Fato:* Espetáculo "A Serpente" (2014).

**Modelo Antigo (4+ Arquivos):**
1.  `work-a-serpente.md` (A Obra conceitual)
2.  `event-temporada-2014.md` (A temporada específica)
3.  `participation-rafael-ator.md` (Minha atuação na temporada)
4.  `record-foto-cena.md` (Uma foto da peça)

**Modelo Novo (2 Arquivos):**
1.  `work-a-serpente.md` (A Obra conceitual - Mantida)
2.  `action-temporada-a-serpente-2014.md` (A realização concreta)

```yaml
# action-temporada-a-serpente-2014.md
type: action
title: "Temporada de Estreia: A Serpente"
role: "Ator e Diretor"
context: "Teatro Sesc (2014)"
related_work: [[work-a-serpente]]
gallery:
  - src: "/media/a-serpente-2014.jpg"
    caption: "Cena do ato final"
    type: "photo"
```
*Análise:* A distinção entre *Obra* (intelectual) e *Ação* (apresentação) é preservada, mas a infraestrutura de eventos e participações é colapsada. A foto (Record) vira um item na lista `gallery`.

---

### Caso 3: Registro de Processo (Vão)
*Fato:* Foto documental do processo criativo de "Vão".

**Modelo Antigo:**
Arquivo `record-work-vao-005.md` contendo apenas metadata e path.

**Modelo Novo:**
Inserção direta no arquivo `work-vao.md` (ou na `Action` de criação):

```yaml
# work-vao.md
...
gallery:
  - src: "/acervo/media/images/work-vao-005.jpeg"
    caption: "Registro de processo criativo"
    tags: ["documentacao", "processo"]
...
```
*Análise:* Eliminação total de arquivos de registro. A busca por imagens torna-se uma iteração dentro dos arquivos de Obras/Ações, o que é computacionalmente mais simples para gerar o site estático.

---

## 5. Avaliação Final

### O novo modelo é mais adequado?
**SIM.** Para um portfólio pessoal, a atomização da informação em `Events` e `Participations` cria uma fricção desnecessária na entrada de dados. O modelo baseado em `Action` reflete mentalmente como um artista organiza seu currículo: "Fiz X, no lugar Y, com função Z".

### O que foi removido (Complexidade)
*   Gestão de IDs para eventos triviais.
*   Necessidade de "joins" complexos para renderizar uma página simples (ex: buscar todas as participações, depois buscar seus eventos, depois buscar seus agentes).
*   Poluição do *filesystem* com arquivos de Records.

### O que foi introduzido (Complexidade/Riscos)
*   **Consistência de Strings (Risco Médio):** O campo `context` é texto livre. Se num lugar for escrito "Sesc Belenzinho" e no outro "SESC Belenzinho", a agregação automática falhará. Recomenda-se uso de linter ou autocomplete.
*   **Perda de Metadados de Mídia (Risco Baixo):** Ao transformar Records em lista YAML, perde-se a capacidade de ter campos customizados complexos por imagem (ex: ISO, Câmera, Autor da foto detalhado) *se* o schema da galeria não prever isso. Para um portfólio, `caption` e `credits` costumam bastar.

### Veredito
A migração é recomendada. O custo cognitivo de manutenção cairá drasticamente, e a estrutura de `Action` + `Context` é robusta o suficiente para gerar tanto o site (timeline, currículo) quanto visualizações de dados futuras, sem o peso do modelo institucional anterior.
