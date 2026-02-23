const { Plugin, Notice } = require('obsidian');
const fs = require('fs');
const path = require('path');

module.exports = class ImageFocusTaggerPlugin extends Plugin {
    onload() {
        console.log('Loading Image Focus Tagger Plugin');

        // Adiciona estilo para o cursor de mira nas imagens
        const styleId = 'image-focus-tagger-style';
        if (!document.getElementById(styleId)) {
            const style = document.createElement('style');
            style.id = styleId;
            style.innerHTML = `
                .workspace-leaf-content[data-type="image"] img { cursor: crosshair !important; }
            `;
            document.head.appendChild(style);
        }

        this.registerDomEvent(document, 'click', async (evt) => {
            if (evt.target.tagName !== 'IMG') return;
            const img = evt.target;

            // Tenta encontrar a view de imagem de forma mais robusta
            const leaf = app.workspace.getLeavesOfType('image').find(l => l.view.contentEl.contains(img));
            if (!leaf) return;

            const file = leaf.view.file;
            if (!file) return;

            // Caminho absoluto do arquivo
            const absolutePath = app.vault.adapter.getFullPath(file.path);

            // Calcula coordenadas
            const rect = img.getBoundingClientRect();
            const clickX = evt.clientX - rect.left;
            const clickY = evt.clientY - rect.top;

            const normX = Math.max(0, Math.min(1, clickX / rect.width));
            const normY = Math.max(0, Math.min(1, clickY / rect.height));

            const xStr = normX.toFixed(4);
            const yStr = normY.toFixed(4);

            // Bolinha verde
            const dot = document.createElement('div');
            dot.style.position = 'fixed';
            dot.style.width = '14px';
            dot.style.height = '14px';
            dot.style.backgroundColor = '#00ff00';
            dot.style.border = '2px solid white';
            dot.style.borderRadius = '50%';
            dot.style.left = (evt.clientX - 7) + 'px';
            dot.style.top = (evt.clientY - 7) + 'px';
            dot.style.pointerEvents = 'none';
            dot.style.zIndex = '99999';
            dot.style.boxShadow = '0 0 10px rgba(0,255,0,0.8)';
            document.body.appendChild(dot);
            setTimeout(() => {
                dot.style.transition = 'all 0.5s ease-out';
                dot.style.opacity = '0';
                dot.style.transform = 'scale(3)';
                setTimeout(() => dot.remove(), 500);
            }, 100);

            // Executa Go CLI
            const { spawn } = require('child_process');
            const cliPath = app.vault.adapter.getFullPath('cli');

            const exePath = path.join(cliPath, 'acervo.exe');
            const isExe = fs.existsSync(exePath);
            const cmd = isExe ? exePath : 'go';
            const args = isExe
                ? ['set-focus', `"${absolutePath}"`, xStr, yStr]
                : ['run', 'main.go', 'set-focus', `"${absolutePath}"`, xStr, yStr];

            console.log('Running CLI:', { cmd, args, cwd: cliPath });

            const proc = spawn(cmd, args, {
                cwd: cliPath,
                shell: true
            });

            let stdout = '';
            let stderr = '';
            proc.stdout.on('data', (data) => stdout += data.toString());
            proc.stderr.on('data', (data) => stderr += data.toString());

            proc.on('close', (code) => {
                if (code !== 0) {
                    const errorMsg = (stderr || stdout || 'Falha silenciosa').trim();
                    new Notice(`❌ Erro (${code}): ${errorMsg}`);
                    console.error('CLI Execution Failed:', { code, stdout, stderr });
                    return;
                }
                new Notice(`✅ Foco salvo em ${file.name}`);
                console.log('CLI Success:', stdout);

                const yamlResult = `XMP_RegionInfo:\n  - RegionType: Focus\n    X: ${xStr}\n    Y: ${yStr}\n    W: 0.0000\n    H: 0.0000`;
                navigator.clipboard.writeText(yamlResult);
            });

            proc.on('error', (err) => {
                new Notice(`❌ Erro ao iniciar processo: ${err.message}`);
                console.error('Spawn Error:', err);
            });
        });
    }

    onunload() {
        console.log('Unloading Image Focus Tagger Plugin');
    }
}
