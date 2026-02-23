# Relatório de Auditoria de Entidades e Ações Recomendadas

Este relatório detalha as inconsistências identificadas no acervo e propõe ações corretivas específicas baseadas na análise dos documentos extraídos (`extracted_texts`).

---

## 1. Entidades de Teste e Lixo (Remoção Imediata)

### `work-dummy2`
- **Diagnóstico:** Entidade de teste sem valor documental ("Dummy Work 2").
- **Ação:** **Excluir** o arquivo `acervo/entities/works/work-dummy2.md`.

### `agent-foo`
- **Diagnóstico:** Entidade de teste ("foo").
- **Ação:** **Excluir** o arquivo `acervo/entities/agents/agent-foo.md`.

---

## 2. Correções de Títulos e Metadados

### `work-constelacao`
- **Diagnóstico:** Título corrompido ("ConstelaÃ§Ã£o") e metadados genéricos (`language: teatro | audiovisual...`).
- **Evidência:** O termo aparece em `portfolio_rafael.pdf.txt` (pág. 26) e há registros relacionados (`record-constelacao-2023-001.md`). Refere-se a um processo de residência ou criação em 2023.
- **Ação:**
    1.  Renomear `title` para **"Constelação"**.
    2.  Definir `language` como **"processo criativo"** ou **"performance"** (conforme a natureza exata, provavelmente performance/residência).
    3.  Atualizar `status` para **"concluída"** (2023).

### `work-convite`
- **Diagnóstico:** Título genérico ("Convite"), mas refere-se a uma peça teatral real de 2017.
- **Evidência:** `portfólio_semino_antigo.docx.txt`: "Peça 'Convite', 2017".
- **Ação:**
    1.  Manter a entidade.
    2.  Atualizar `description` para incluir: **"Espetáculo teatral encenado em 2017."**
    3.  Se possível, verificar se há subtítulo ou contexto de grupo (ex: Cia Del Artes?) para desambiguar no futuro.

### `work-a-serpente`
- **Diagnóstico:** ID e título sugerem a obra original de Nelson Rodrigues, não a montagem específica de Rafael Semino.
- **Evidência:** `portfólio_semino_antigo.docx.txt`: "A serpente, 2014".
- **Ação:**
    1.  Renomear ID para **`work-a-serpente-montagem-2014`** (se a política de IDs permitir mudança) ou manter ID e ajustar título.
    2.  Alterar `title` para **"A Serpente (Montagem 2014)"**.
    3.  Na `description`, explicitar: **"Montagem do texto de Nelson Rodrigues, realizada em 2014 com direção/atuação de Rafael Semino."**

---

## 3. Resolução de Ambiguidades de Participação (Vínculo com Eventos)

### `participation-avaliador-junino`
- **Diagnóstico:** Participação "flutuante" sem evento vinculado.
- **Evidência:** `atualizacao_potfolio.pdf.txt` (pág. 8): "Atuação como avaliador junino (2023)".
- **Ação:**
    1.  Criar a entidade **`event-ciclo-junino-ceara-2023`**.
        -   `type`: event
        -   `title`: **"Ciclo de Festivais Juninos do Ceará 2023"**
        -   `date_start`: **"2023-06-01"** (estimado, verificar datas exatas se possível ou usar mês).
        -   `location`: **"Ceará"**
    2.  Vincular `participation-avaliador-junino` a este novo evento (`event: [[event-ciclo-junino-ceara-2023]]`).

### `participation-avaliador-ciclo-carnavalesco`
- **Diagnóstico:** Participação sem evento vinculado.
- **Evidência:** `atualizacao_potfolio.pdf.txt` (pág. 10): "Atuação como avaliador do Ciclo Carnavalesco da Avenida Domingos Olímpio (2020)".
- **Ação:**
    1.  Criar a entidade **`event-ciclo-carnavalesco-2020`**.
        -   `type`: event
        -   `title`: **"Ciclo Carnavalesco 2020"**
        -   `location`: **"Avenida Domingos Olímpio, Fortaleza - CE"**
        -   `date_start`: **"2020-02-01"** (estimado carnaval).
    2.  Vincular `participation-avaliador-ciclo-carnavalesco` a este evento.

### `participation-projeto-angola-bie`
- **Diagnóstico:** Nome confunde "Projeto" (Evento/Obra) com "Participação" (Ação).
- **Evidência:** `curriculo_rafael_semino.pdf.txt` (pág. 2): "Cia Del Artes — Professor Voluntário... Angola – Província do Bié... Outubro de 2018 a Novembro de 2018".
- **Ação:**
    1.  Renomear ID para **`participation-rafael-ensino-angola-2018`**.
    2.  Alterar `title` para **"Ensino de Artes e Cinema em Angola"**.
    3.  Criar a entidade **`event-intercambio-angola-2018`**.
        -   `type`: event
        -   `title`: **"Intercâmbio Cultural e Educativo em Angola"**
        -   `location`: **"Província do Bié, Angola"**
        -   `date_start`: **"2018-10-01"**
        -   `date_end`: **"2018-11-30"**
    4.  Vincular a participação a este evento.

### `participation-jogos-teatrais`
- **Diagnóstico:** Termo genérico, entidade gerada automaticamente ("auto-patched").
- **Evidência:** `curriculo_rafael_semino.pdf.txt` (pág. 2) menciona a disciplina "Jogos e Africanidade" ministrada na Escola Mário Hugo Sadrak.
- **Ação:**
    1.  Verificar se `participation-prof-hugo-sadrack` já cobre esta atividade.
    2.  Se sim, **excluir** `participation-jogos-teatrais` e adicionar "Ministração da disciplina Jogos e Africanidade" à descrição/notas de `participation-prof-hugo-sadrack`.
    3.  Se não for redundante, renomear para algo específico (ex: `participation-oficina-jogos-teatrais-local-ano`) se houver evidência de outro contexto. Caso contrário, assumir redundância e excluir.

---

## 4. Deduplicação e Fusão

### `participation-formacao-ufba` vs `participation-rafael-pos-ufba`
- **Diagnóstico:** Duplicidade. A primeira é automática, a segunda é manual e correta.
- **Evidência:** `portfolio_rafael.pdf.txt` (pág. 2) confirma "Especialização em Estudos em Teatro do Oprimido / UFBA".
- **Ação:**
    1.  Transferir quaisquer dados úteis de `participation-formacao-ufba` (se houver) para `participation-rafael-pos-ufba`.
    2.  **Excluir** `participation-formacao-ufba`.
    3.  Garantir que `participation-rafael-pos-ufba` esteja vinculada ao evento `event-pos-teatro-oprimido` (verificar se existe ou criar).

---

## 5. Estruturação de Prêmios (Event vs Participation)

### `event-premio-amarracoes-esteticas`
- **Diagnóstico:** Ambiguidade entre a *Cerimônia* e a *Conquista*.
- **Evidência:** `portfolio_coletivo_farol_novo.pdf.txt` (pág. 3) diz: "...lançou o prêmio 'Amarrações Estéticas'. Esse prêmio incentivava artistas...".
- **Ação:**
    1.  Manter `event-premio-amarracoes-esteticas` representando o **Edital/Evento de Premiação** da Escola Porto Iracema (2023).
    2.  Garantir que exista uma **Participação** (ex: `participation-rafael-premio-amarracoes`) do tipo "Premiação" ou "Contemplado".
        -   `agent`: `[[agent-rafael-semino]]` (e `[[agent-zeis]]` se aplicável).
        -   `event`: `[[event-premio-amarracoes-esteticas]]`.
        -   `role`: **"Artista Premiado"** ou **"Contemplado"**.
        -   `related_to`: `[[work-vao]]` (pois o texto diz que o prêmio resultou no espetáculo "Vão").

---

## Resumo das Ações Críticas
1.  **Excluir:** `work-dummy2`, `agent-foo`, `participation-jogos-teatrais` (após merge), `participation-formacao-ufba`.
2.  **Criar Eventos:** `event-ciclo-junino-ceara-2023`, `event-ciclo-carnavalesco-2020`, `event-intercambio-angola-2018`.
3.  **Renomear/Corrigir:** `work-constelacao` (título), `work-a-serpente` (título/contexto), `participation-projeto-angola-bie` (ID e título).
4.  **Vincular:** Participações órfãs aos novos eventos criados.
