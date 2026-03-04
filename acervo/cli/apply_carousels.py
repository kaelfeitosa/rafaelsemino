import re

html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

# Helper to generate carousel HTML
def make_carousel(images):
    slides = ""
    for idx, (src, alt) in enumerate(images):
        # The first slide can be active, but our CSS usually handles that with JS.
        # Let's just generate standard structure.
        slides += f"""                                <div class="carousel-slide"><focus-image
                                        src="{src}"
                                        alt="{alt}"></focus-image>
                                </div>
"""
    
    return f"""<div class="carousel-container carousel-container--cols-2">
                            <div class="carousel-wrapper">
{slides.rstrip()}
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>"""

### 1. PERFIL
perfil_old = """<div class="img-stack img-stack--cols-2">
                            <focus-image src="images/optimized/agent-rafael-semino/agent-rafael-semino-008.webp"
                                alt="Retrato de Rafael Semino em performance"></focus-image>
                            <focus-image src="images/optimized/agent-rafael-semino/agent-rafael-semino-005.webp"
                                alt="Rafael Semino atuando no palco"></focus-image>
                        </div>"""
perfil_images = [
    ("images/optimized/agent-rafael-semino/agent-rafael-semino-008.webp", "Retrato de Rafael Semino em performance"),
    ("images/optimized/agent-rafael-semino/agent-rafael-semino-005.webp", "Rafael Semino atuando no palco"),
    ("images/optimized/agent-rafael-semino/agent-rafael-semino-002.webp", "Rafael Semino em movimento e expressão"),
    ("images/optimized/agent-rafael-semino/agent-rafael-semino-001.webp", "Performance teatral de Rafael Semino")
]
perfil_new = make_carousel(perfil_images)
html_content = html_content.replace(perfil_old, perfil_new)


### 2. EXU NÃO VEM HOJE
exu_old = """<div class="img-stack img-stack--cols-3">
                            <focus-image src="images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-008.webp"
                                alt="Rafael Semino em cena no espetáculo Exu Não Vem Hoje"></focus-image>
                            <focus-image src="images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-010.webp"
                                alt="Performance teatral de Exu Não Vem Hoje"></focus-image>
                            <focus-image src="images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-012.webp"
                                alt="Público e interação em Exu Não Vem Hoje"></focus-image>
                        </div>"""
exu_images = [
    ("images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-010.webp", "Foco na expressão do ator em Exu Não Vem Hoje"),
    ("images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-008.webp", "O corpo no espaço cênico"),
    ("images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-012.webp", "Relação e quebra da quarta parede com o público"),
    ("images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-001.webp", "Registro geral da atmosfera do espetáculo")
]
exu_new = make_carousel(exu_images)
html_content = html_content.replace(exu_old, exu_new)


### 3. VÃO
vao_old = """<div class="img-stack img-stack--cols-2">
                            <focus-image src="images/optimized/work-vao/work-vao-003.webp"
                                alt="Cena do espetáculo Vão"></focus-image>
                            <focus-image src="images/optimized/work-vao/work-vao-004.webp"
                                alt="Apresentação musical no espetáculo Vão"></focus-image>
                        </div>"""
vao_images = [
    ("images/optimized/work-vao/work-vao-003.webp", "Cena do espetáculo Vão: registro do ator em performance"),
    ("images/optimized/work-vao/work-vao-004.webp", "Vão: interação visual com elementos musicais"),
    ("images/optimized/work-vao/work-vao-005.webp", "Vão: contraste de luz e pesquisa de sombras"),
    ("images/optimized/work-vao/work-vao-006.webp", "Vão: perspectiva ampla do minimalismo cênico")
]
vao_new = make_carousel(vao_images)
html_content = html_content.replace(vao_old, vao_new)


### 4. REBORDOSE
rebordose_old = """<div class="carousel-container carousel-container--cols-2">
                            <div class="carousel-wrapper">
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rebordose/work-rebordose-001.webp"
                                        alt="Cena do curta-metragem Rebordose"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rebordose/work-rebordose-002.webp"
                                        alt="Ator atuando no curta-metragem Rebordose"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-012.webp"
                                        alt="Performance noturna na rua com espectadores do projeto Exu Não Vem Hoje"></focus-image>
                                </div>
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>"""
reb_images = [
    ("images/optimized/work-rebordose/work-rebordose-001.webp", "Cena inicial do curta Rebordose"),
    ("images/optimized/work-rebordose/work-rebordose-002.webp", "Interação entre os atores na quitinete"),
    ("images/optimized/work-rebordose/work-rebordose-003.webp", "Expressão dramática do protagonista"),
    ("images/optimized/work-rebordose/work-rebordose-004.webp", "Fotografia periférica do curta"),
    ("images/optimized/work-rebordose/work-rebordose-005.webp", "Clímax da narrativa de Rebordose")
]
rebordose_new = make_carousel(reb_images)
html_content = html_content.replace(rebordose_old, rebordose_new)


### 5. HÁ NÚMEROS QUE SONHAM
numeros_old = """<div class="carousel-container carousel-container--cols-2">
                            <div class="carousel-wrapper">
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-001.webp"
                                        alt="Poética visual do curta Há números que sonham"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-002.webp"
                                        alt="Cena conceitual do curta"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-003.webp"
                                        alt="Enquadramento do curta Há números que sonham"></focus-image>
                                </div>
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>"""
num_images = [
    ("images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-001.webp", "Poética visual do curta Há números que sonham"),
    ("images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-002.webp", "Cena conceitual do curta"),
    ("images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-003.webp", "Enquadramento do curta Há números que sonham"),
    ("images/optimized/work-ha-numeros-que-sonham/work-ha-numeros-que-sonham-004.webp", "Desfecho visual narrativo do curta")
]
numeros_new = make_carousel(num_images)
html_content = html_content.replace(numeros_old, numeros_new)


with open(html_file, 'w', encoding='utf-8') as f:
    f.write(html_content)

print("Carousel blocks successfully replaced in index.html.")
