import os

base_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities"

targets = {
    "events/event-aceleracao-idoso.md": "Programa de Aceleração do Idoso",
    "events/event-licenciatura-ifce.md": "Licenciatura em Teatro",
    "events/event-mestrado-ufc.md": "Mestrado em Artes",
    "events/event-laboratorio-pesquisa-ccbj.md": "Laboratório de Pesquisa CCBJ",
    "events/event-percurso-basico-teatro.md": "Percurso Básico de Teatro",
    "events/event-pos-teatro-oprimido.md": "Pós-graduação em Teatro do Oprimido",

    "agents/agent-ccbj.md": "CCBJ - Centro Cultural Bom Jardim",
    "agents/agent-escola-hugo-sadrack.md": "Escola Hugo Sadrack do Vale",
    "agents/agent-escola-paulo-petrola.md": "Escola Paulo Petrola",
    "agents/agent-ifce.md": "IFCE - Instituto Federal do Ceará",
    "agents/agent-ufba.md": "UFBA - Universidade Federal da Bahia",
    "agents/agent-ufc.md": "UFC - Universidade Federal do Ceará"
}

for rel_path, name_val in targets.items():
    full_path = os.path.join(base_dir, rel_path)
    if os.path.exists(full_path):
        with open(full_path, "r", encoding="utf-8") as f:
            content = f.read()
        
        if "name:" not in content:
            # We insert name: "..." right after title: "..."
            import re
            content = re.sub(r'(title: ".+?")', r'\1\nname: "' + name_val + '"', content)

            with open(full_path, "w", encoding="utf-8") as f:
                f.write(content)
            print(f"✅ Fixed {rel_path}")
        else:
            print(f"⚠️ {rel_path} already has name.")
    else:
        print(f"❌ {rel_path} not found.")

print("Patch complete.")
