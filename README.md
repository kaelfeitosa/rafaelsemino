# Acervo Cultural: Rafael Semino

Bem-vindo(a) ao repositório do portfólio e acervo cultural de Rafael Semino, focado na preservação de sua trajetória de vida nas artes cênicas, pesquisa, docência e produção audiovisual.

O projeto evoluiu de uma simples página estática para um **Motor de Dados Culturais e Semânticos**, operando com uma clara Separação de Responsabilidades (Separation of Concerns).

## 🚀 Arquitetura do Projeto

O repositório está estruturado em duas camadas principais:

### 1. Camada Web (`/frontend`)
Responsável pela exibição e interface com o usuário final.
- **`index.html`**: Master do Portfólio.
- **`assets/`**: Estilizações CSS e lógicas JS estáticas.
- **`images/`**: Imagens isoladas da interface visual (icones, logos).

### 2. Camada de Dados (`/acervo`)
O "Cérebro" do projeto. Funciona como um CMS Nativo completamente rastreável via Git. Toda a história e catálogo artístico de Rafael Semino estão mapeados semanticamente aqui.

- **`entities/`**: Arquivos Markdown (`.md`) puros que servem como "Nós" de dados para **Agents** (Pessoas/Organizações), **Works** (Obras/Espetáculos), **Events** (Festivais) e **Participations** (Atuações Específicas).
- **`data/records/`**: Cada foto/mídia do acervo original possui um Record `.md` exclusivo. Ele mapeia a mídia à entidade correta, injetando *Contexto Textual* e *Geolocalização* exatos no momento em que a foto foi tirada.
- **`media/images/`**: Onde as fotos canônicas e renomeadas semanticamente estão abrigadas.

### 🧠 Motor em Go (`/acervo/cli`)
Para garantir escabilidade e buscas indexadas (Search/API), o projeto possui uma Command-Line Interface construída em **Go**. Ela lê todo o seu histórico em Markdown e tece um `db.sqlite` relacional e ultrarrápido.

#### Como utilizar os Comandos do Motor:
1. Abra o terminal e navegue: `cd acervo/cli`
2. Compile/Rode o atualizador de Índice:
   ```bash
   go run . reindex
   ```
   *Isso força o robô a reler todas as Entidades Markdown e atualizar o `db.sqlite` principal.*
3. Valide a Saúde do Grafo:
   ```bash
   go run . audit
   ```
   *Checa imediatamente se há quebra de links ou entidades sem descrição.*

---
Desenvolvido para garantir que toda a linha do tempo e catalogação etnográfica permaneça imutável, versionada e pronta para ser consumida por APIs no futuro.
