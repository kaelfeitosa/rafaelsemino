# Acervo - Setup e DevOps Local

Este diretório contém a estrutura de armazenamento Event-Centric do portfólio.
Para garantir operações seguras e consistentes, implementamos práticas de DevOps local.

## 1. Setup Inicial

Abra o terminal **(recomendamos Git Bash ou WSL caso esteja no Windows)**, navegue para a pasta `acervo/` e execute:

```bash
cd cli && go run . setup
```

**O que este comando faz:**
1. **Instala Git Hooks**: Configura um gancho (hook) no seu repositório Git. Toda vez que você executar `git commit`, o sistema rotina automaticamente a validação dos schemas. Se qualquer arquivo estiver defeituoso, o commit é impedido.

## 2. Operações no Dia a Dia

Caso precise validar ou regenerar a base manualmente, utilize os comandos dentro da pasta `acervo/cli`:

- **`go run . validate`**: Valida a estrutura YAML e a presença de campos essenciais de todos os `.md`.
- **`go run . reindex`**: Destrói e reconstrói completamente a base `db.sqlite` iterando recursivamente pelos arquivos dentro de `entities/`.
- **`go run . audit`**: Detecta e avisa sobre lacunas nos dados (ex: eventos sem obra, participações sem registros anexados).
