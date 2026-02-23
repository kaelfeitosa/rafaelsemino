# Auditoria de Entidades - Relatório de Inconsistências

Este relatório lista as entidades identificadas no acervo cujos tipos parecem inadequados, ambíguos ou que apresentam problemas estruturais, conforme as diretrizes da auditoria.

## Entidades de Teste / Lixo

### ID: work-dummy2
- **Tipo Atual:** Work
- **Problema:** Entidade de teste ("Dummy Work 2").
- **Sugestão:** Remover.
- **Impacto:** Limpeza do banco de dados e remoção de ruído.

### ID: agent-foo
- **Tipo Atual:** Agent
- **Problema:** Entidade de teste ("foo").
- **Sugestão:** Remover.
- **Impacto:** Limpeza do banco de dados e remoção de ruído.

## Entidades com Identificação ou Título Problemáticos

### ID: work-constelacao
- **Tipo Atual:** Work
- **Problema:** Título malformado ("ConstelaÃ§Ã£o"). Campos com valores de placeholder (`language: teatro | audiovisual ...`, `status: em-andamento | concluida ...`).
- **Sugestão:** Corrigir título para "Constelação", definir `language` e `status` corretos baseados na obra real.
- **Impacto:** Corrige erros de codificação e garante a integridade dos dados.

### ID: work-convite
- **Tipo Atual:** Work
- **Problema:** O título "Convite" é genérico e pode ser confundido com um tipo de registro (convite impresso). A descrição indica que é um "Espetáculo de teatro".
- **Sugestão:** Verificar se este é o título completo e correto. Se possível, adicionar subtítulo ou contexto para desambiguar (ex: "O Convite").
- **Impacto:** Evita confusão semântica entre a obra e um documento.

### ID: work-a-serpente
- **Tipo Atual:** Work
- **Problema:** "A Serpente" é uma obra clássica de Nelson Rodrigues. A entidade atual representa a *encenação* feita por Rafael Semino, mas o ID e título sugerem a obra original.
- **Sugestão:** Clarificar na descrição que se trata de uma encenação/montagem. Considerar renomear o ID para algo como `work-encenacao-a-serpente` se o sistema permitir distinção entre obra original e derivada.
- **Impacto:** Distingue a criação do diretor da obra do dramaturgo original.

## Entidades Ambíguas ou Mal Classificadas

### ID: participation-avaliador-junino
- **Tipo Atual:** Participation
- **Problema:** Descreve um papel/função ("Avaliador") ao longo de um período ("ano 2023") sem vínculo explícito a um Evento específico. A "Participação" deve descrever uma ação concreta em um evento.
- **Sugestão:** Vincular a um evento agregador (ex: `event-ciclo-junino-2023`) ou dividir em participações específicas se houver eventos distintos documentados.
- **Impacto:** Ancora a participação em um contexto temporal e evental concreto.

### ID: participation-avaliador-ciclo-carnavalesco
- **Tipo Atual:** Participation
- **Problema:** Similar ao anterior. Descreve a função de avaliador sem vínculo direto a um evento cadastrado (o evento "Ciclo Carnavalesco" está implícito no texto).
- **Sugestão:** Criar ou vincular ao evento `event-ciclo-carnavalesco-2020`.
- **Impacto:** Garante que a participação tenha um "lugar" (evento) onde ocorreu.

### ID: participation-projeto-angola-bie
- **Tipo Atual:** Participation
- **Problema:** O ID sugere que a participação *é* o projeto. "Projeto Angola Bié" soa como um Evento (a viagem/intercâmbio) ou Work (o projeto em si).
- **Sugestão:** Renomear para `participation-rafael-projeto-angola-bie` e garantir que exista o `event-projeto-angola-bie` ao qual esta participação se conecta.
- **Impacto:** Clarifica a distinção entre o evento (projeto) e a ação do agente (participação).

### ID: participation-jogos-teatrais
- **Tipo Atual:** Participation
- **Problema:** "Jogos Teatrais" é um conceito ou metodologia. Como ID de participação, é ambíguo. Se refere a uma oficina ministrada ou cursada, deve ser explícito. Atualmente é uma entidade "auto-patched" (gerada automaticamente).
- **Sugestão:** Identificar o evento específico (ex: Oficina de Jogos Teatrais) e renomear/vincular corretamente. Se for apenas um tópico, não deve ser uma participação isolada sem evento.
- **Impacto:** Resolve a ambiguidade entre conceito e ação concreta.

## Duplicatas e Entidades Automáticas

### ID: participation-formacao-ufba
- **Tipo Atual:** Participation
- **Problema:** Entidade gerada automaticamente ("auto-patched") que parece duplicar `participation-rafael-pos-ufba` (que já descreve a pós-graduação na UFBA).
- **Sugestão:** Fundir com `participation-rafael-pos-ufba` e remover a duplicata.
- **Impacto:** Elimina redundância e consolida as informações.

### ID: event-premio-amarracoes-esteticas
- **Tipo Atual:** Event
- **Problema:** Classificado como Evento, mas "Prêmio" pode ser interpretado como o objeto (Record) ou a conquista (Participation). Se representar a *cerimônia* de premiação, o tipo Event é aceitável, mas deve ser verificado se não duplica a informação do prêmio em si.
- **Sugestão:** Se for a cerimônia, manter como Event. Se for o prêmio (objeto/conquista), considerar modelar como Record vinculado a uma Participation de "Premiação".
- **Impacto:** Garante a modelagem correta de "conquistas" vs "eventos".
