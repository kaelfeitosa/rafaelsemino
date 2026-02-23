# Proposta de Reestruturação do `index.html`

## Visão Geral
A nova organização do `index.html` deve refletir a estrutura relacional do `acervo` (Agente -> Obra -> Participação -> Evento), transformando o portfólio estático em uma narrativa dinâmica da trajetória de Rafael Semino. O foco deixa de ser listas isoladas para se tornar uma história coesa de **Criação, Pesquisa e Ensino**.

---

## Lógica Narrativa (O "Porquê")
A estrutura proposta segue o arco da **Pesquisa-Criação**:
1.  **Quem sou (Identidade/Bio):** Estabelece o lugar de fala (Artista-Pesquisador).
2.  **O que crio (Obras):** A materialização da pesquisa em produtos artísticos.
3.  **Como ensino (Pedagogia):** A aplicação da pesquisa na formação de outros.
4.  **Onde atuo (Trajetória):** O reconhecimento e a inserção no campo profissional.

---

## Estrutura Proposta

### 1. Cabeçalho & Identidade (Hero)
*   **Conceito:** Manter o impacto visual atual, mas refinar a "tagline" para alinhar com o perfil de pesquisador.
*   **Título:** RAFAEL SEMINO
*   **Subtítulo:** Ator, Diretor, Dramaturgo & Pesquisador.
*   **Bio Curta (Floating):** "Investigação continuada sobre mito, oralidade e ancestralidade na cena contemporânea."

### 2. Sobre (Bio & Perfil)
*   **Fonte de Dados:** `acervo/entities/agents/agent-rafael-semino.md`
*   **Conteúdo:**
    *   Foto de perfil atualizada.
    *   Texto biográfico focado na tríade: **Criação Artística + Pesquisa Acadêmica (Mestrado UFC) + Ação Pedagógica**.
    *   Destaque para a formação: Mestre em Artes (UFC), Licenciado (IFCE), Pós-graduando (UFBA).

### 3. Obras em Destaque (Criação & Autoria)
*   **Fonte de Dados:** `acervo/entities/works/*.md`
*   **Formato:** Grid visual (Estilo Xilogravura).
*   **Categorias Sugeridas:**
    *   **Cênicas (Teatro & Performance):**
        *   *Exu Não Vem Hoje* (Destaque Principal - Coletivo Farol Novo).
        *   *Vão* (Teatro + Música).
        *   *Trapo Preto*.
    *   **Audiovisual:**
        *   *Rebordose* (Curta).
        *   *Mundo-Imagem* (Websérie).
        *   *Astronauta* (Produção Angola).
    *   **Escrita & Dramaturgia:**
        *   *Contos de Exu* (Livro).
        *   *Cala-me os Olhos* (Dramaturgia).

### 4. Trajetória Pedagógica & Pesquisa (O Diferencial)
*   **Conceito:** Esta seção é crucial para editais e para o perfil de educador.
*   **Subseções:**
    *   **Ensino Formal:** Escolas (Paulo Petrola, Hugo Sadrack).
    *   **Projetos & Cursos Livres:** Porto Iracema (Projeto Abarca, Percurso Básico).
    *   **Pesquisa & Tradição:** Grupo Miraira (IFCE), Reisado, Mestres do Mundo. *Aqui conecta a pesquisa acadêmica com a prática popular.*

### 5. Atuação & Performance (Colaborações)
*   **Fonte de Dados:** `acervo/entities/participations/participation-ator-*.md`
*   **Foco:** Trabalhos onde Rafael atua como intérprete em obras de terceiros.
*   **Lista Selecionada:**
    *   *Irreversível* (Dir. Caique Melo).
    *   *Santo Bordel de Tiatira*.
    *   *A Serpente*.

### 6. Gestão, Produção & Curadoria
*   **Fonte de Dados:** `acervo/entities/participations/participation-produtor-*.md` e `participation-avaliador-*.md`
*   **Conteúdo:**
    *   Coordenação pedagógica (Azusa).
    *   Produção (Black Heroes).
    *   Avaliação (Ciclo Junino, Carnaval).

### 7. Coletivo Farol Novo (Identidade Coletiva)
*   **Fonte de Dados:** `acervo/entities/agents/agent-coletivo-farol-novo.md`
*   **Conceito:** Apresentar o coletivo não apenas como um "grupo", mas como uma **plataforma de pesquisa** contínua.
*   **Destaque:** Link para *Exu Não Vem Hoje* e *Zona de Criação*.

### 8. Rodapé & Contato
*   **Informações:** E-mail, Redes Sociais (@coletivofarolnovo).
*   **Territórios:** Fortaleza, Luanda, Itapipoca (destacar a atuação internacional/estadual).

---

## Melhorias Técnicas & Visuais
1.  **Responsividade:** O grid atual precisa de ajustes para mobile (colunas empilhadas).
2.  **Imagens:** Substituir os placeholders `acervo_pendente` por imagens reais extraídas de `acervo/entities/records`.
3.  **Links:** Cada card de "Obra" deve, idealmente, abrir um modal ou levar a uma página de detalhe (se houver) ou link externo (YouTube/Drive) conforme o metadado `url` no `acervo`.
