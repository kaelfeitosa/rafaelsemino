# Acervo Cultural: Rafael Semino

Bem-vindo(a) ao repositório do portfólio e acervo cultural de Rafael Semino, focado na preservação de sua trajetória de vida nas artes cênicas, pesquisa, docência e produção audiovisual.

O projeto opera como um **Motor de Dados Culturais e Semânticos**, garantindo que toda a linha do tempo e catalogação editorial permaneçam imutáveis, versionadas e prontas para consumo.

## 🚀 Arquitetura do Projeto

O repositório está estruturado em duas camadas principais:

### 1. Camada Web (`/frontend`)
Responsável pela exibição e interface com o usuário final. Aqui os dados do acervo são consumidos e apresentados em um portfólio moderno e interativo.

### 2. Camada de Dados (`/acervo`)
O "Cérebro" do projeto. Funciona como um banco de dados editorial em Markdown. Toda a história artística de Rafael Semino está mapeada semanticamente aqui.

- **`entities/`**: Arquivos Markdown puros que servem como "Nós" da rede:
    - **Works**: O núcleo do sistema. As obras artísticas, pesquisas, e projetos (espetáculos, livros, curtas, participações). Cada Work abriga **Occurrences** (eventos, apresentações, prêmios) aninhadas em seu interior, mapeando a linha do tempo.
    - **Agents**: Rafael Semino e os coletivos/instituições com quem colaborou.
- **Evidências (Attachments)**: Documentos e mídias (fotos de cena, pôsteres, press clippings, PDFs) são anexados diretamente aos Works, ganhando tags de `category` para fácil filtragem no portfólio.

### 🧠 Motor em Go (`/acervo/cli`)
Para garantir integridade e velocidade de busca, o projeto possui uma Command-Line Interface (CLI) em **Go**. Ela transforma os arquivos Markdown em um banco de dados relacional SQLite (`db.sqlite`).

#### Comandos Principais:
1. **Sincronizar Banco:**
   ```bash
   cd acervo/cli
   go run main.go reindex
   ```
2. **Validar Dados:**
   ```bash
   go run main.go validate
   ```
3. **Verificar Integridade:**
   ```bash
   go run main.go verify
   ```

---
Para detalhes técnicos sobre o modelo editorial e regras de manutenção, consulte o arquivo [AGENTS.md](AGENTS.md).
