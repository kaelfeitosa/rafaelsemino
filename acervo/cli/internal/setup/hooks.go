package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallHooks() error {
	repoRoot := "../../" // Since CLI is at acervo/cli
	hooksDir := filepath.Join(repoRoot, ".git", "hooks")
	preCommitHook := filepath.Join(hooksDir, "pre-commit")

	if _, err := os.Stat(filepath.Join(repoRoot, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("Diretório .git não encontrado.")
	}

	os.MkdirAll(hooksDir, 0755)

	hookContent := `#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
ACERVO_CLI_DIR="$REPO_ROOT/acervo/cli"

if [ -d "$ACERVO_CLI_DIR" ]; then
    echo "==> [Pre-commit] Validando integridade dos dados do Acervo..."
    cd "$ACERVO_CLI_DIR"
    
    if ! go run . validate; then
        echo "❌ ERRO DE COMMIT: Validação semântica falhou."
        exit 1
    fi
fi
`

	err := os.WriteFile(preCommitHook, []byte(hookContent), 0755)
	if err != nil {
		return err
	}

	fmt.Println("✅ Git pre-commit hook instalado com sucesso em", preCommitHook)
	return nil
}
