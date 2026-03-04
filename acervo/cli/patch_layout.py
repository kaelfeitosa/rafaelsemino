import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

# 1. FIX NAVIGATION ORDER
old_nav = """<a href="#producao">Produção</a>
                <a href="#ensino">Ensino</a>
                <a href="#coletivo">Coletivo</a>"""
new_nav = """<a href="#ensino">Ensino</a>
                <a href="#coletivo">Coletivo</a>
                <a href="#producao">Produção</a>"""
html = html.replace(old_nav, new_nav)

# 2. FIX EIXOS
# Insert Tags in Perfil
old_perfil = "atravessando espetáculos, livros, ações pedagógicas e projetos audiovisuais.</p>\n                        </div>"
new_perfil = """atravessando espetáculos, livros, ações pedagógicas e projetos audiovisuais.</p>
                            <div class="eixos-tags" style="margin-top: 30px; display: flex; gap: 10px; flex-wrap: wrap;">
                                <span class="card-tag card-tag--red" style="margin-bottom: 0;">Teatro e performance</span>
                                <span class="card-tag card-tag--red" style="margin-bottom: 0;">Direção e dramaturgia</span>
                                <span class="card-tag card-tag--red" style="margin-bottom: 0;">Pesquisa de teatro</span>
                                <span class="card-tag card-tag--red" style="margin-bottom: 0;">Ensino e mediação cultural</span>
                            </div>
                        </div>"""
if old_perfil in html:
    html = html.replace(old_perfil, new_perfil)
else:
    print("Could not find the target text to insert eixos tags.")
    
# Remove Eixos Section entirely
html = re.sub(r'<!-- SEÇÃO: EIXOS DE ATUAÇÃO -->\s*<section id="eixos" tabindex="-1">.*?</section>\s+', '', html, flags=re.DOTALL)


# 3. FIX ENSINO
ensino_code_new = """            <!-- SEÇÃO: ENSINO E EXPERIÊNCIA PEDAGÓGICA -->
            <section id="ensino" tabindex="-1">
                <div class="section-label stamped">ENSINO</div>
                <div class="flex-col-gap">

                    <div class="catalog-grid catalog-grid--single-col gap-60">
                        <!-- ANGOLA ENSINO -->
                        <div class="card woodcut-border">
                            <div class="carousel-container carousel-container--cols-2">
                                <div class="carousel-wrapper">
                                    <div class="carousel-slide"><focus-image
                                            src="images/optimized/work-rafael-ensino-angola-2018/work-rafael-ensino-angola-2018-001.webp"
                                            alt="Angola Ensino"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                            src="images/optimized/work-rafael-ensino-angola-2018/work-rafael-ensino-angola-2018-002.webp"
                                            alt="Angola Teatro"></focus-image>
                                    </div>
                                </div>
                                <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                                <button class="carousel-btn carousel-btn--next">&#10095;</button>
                                <div class="carousel-indicators">
                                </div>
                            </div>
                            <div class="card-content">
                                <span class="card-tag">Intercâmbio e Voluntariado</span>
                                <h3 class="project-title">Ensino de Teatro em Luanda (2018)</h3>
                                <div class="project-context mt-20">
                                    <p class="project-desc light">
                                        Atuação voluntária realizada durante o período de graduação. O projeto consistiu na facilitação de práticas em teatro e intercâmbio de experiências pedagógicas com alunos em Luanda, Angola, marcando profundamente sua reflexão como educador.
                                    </p>
                                </div>
                            </div>
                        </div>

                        <!-- ABARCA -->
                        <div class="card woodcut-border card--light">
                            <div class="img-stack img-stack--cols-2">
                                <focus-image
                                    src="images/optimized/work-prof-hugo-sadrack/work-prof-hugo-sadrack-001.webp"
                                    alt="Aula do Projeto ABARCA"></focus-image>
                                <focus-image
                                    src="images/optimized/work-prof-hugo-sadrack/work-prof-hugo-sadrack-002.webp"
                                    alt="Atividade prática com alunos"></focus-image>
                            </div>
                            <div class="card-content">
                                <span class="card-tag card-tag--red">Escola Porto Iracema das Artes 📍 Fortaleza – CE – Brasil</span>
                                <h3 class="course-institution">Projeto ABARCA</h3>
                                <div class="mt-20">
                                    <p class="course-details"><strong>Territórios de atuação:</strong> Vicente Pinzón, Genibaú, Itapipoca.</p>
                                    <p class="course-details mt-20"><strong>Disciplinas ministradas:</strong> Teatro, Dramaturgia, Construção de personagem, Jogos teatrais, Processos de criação.</p>
                                </div>
                            </div>
                        </div>

                        <!-- HISTÓRICO PEDAGÓGICO CONSOLIDADO -->
                        <div class="card woodcut-border card--light">
                            <div class="card-content card-content--medium-padding">
                                <span class="card-tag card-tag--red">Experiência Pedagógica</span>
                                <h2 class="section-header">Histórico Pedagógico</h2>

                                <div class="atuacao-consolidated-grid">
                                    <!-- EDUCAÇÃO BÁSICA -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Educação Básica</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor de Artes</strong> – Escola Paulo Petrola</li>
                                            <li><strong>Professor de Artes</strong> – Escola Hugo Sadrack do Vale</li>
                                        </ul>
                                    </div>

                                    <!-- MEDIAÇÃO E CURSOS LIVRES -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Mediação e Cursos Livres</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor (Percurso Básico de Teatro)</strong> – Porto Iracema das Artes</li>
                                            <li><strong>Professor</strong> – Programa de Aceleração do Idoso</li>
                                            <li><strong>Orientador de Montagem e Mediação</strong> – Projeto no Bairro Terceiro</li>
                                        </ul>
                                    </div>

                                    <!-- ESTÁGIO DOCENTE -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Estágio Docente</p>
                                        <ul class="detail-list">
                                            <li><strong>Noite de Alegria da Rua Trinta e Sete</strong> – Intérprete / Supervisor: Circe Damascena</li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>"""
html = re.sub(r'<!-- SEÇÃO: ENSINO E EXPERIÊNCIA PEDAGÓGICA -->\s*<section id="ensino" tabindex="-1">.*?</section>', ensino_code_new, html, flags=re.DOTALL)

# 4. FIX PRODUÇÃO
old_prod = """                    <div class="card woodcut-border card--light">
                        <div class="card-content">
                            <span class="card-tag">Coordenação de Formação (2018–2023)</span>
                            <h3>Black Heroes</h3>
                            <p>Atuação como Produtor Cultural e Coordenador de Formação.</p>
                        </div>
                    </div>
                    <div class="card woodcut-border">
                        <div class="card-content">
                            <span class="card-tag">Coordenação de Ensino (2018–2023)</span>
                            <h3>Azusa</h3>
                            <p>Coordenação de ensino, pesquisa e processos formativos.</p>
                        </div>
                    </div>
                    <div class="card woodcut-border">
                        <div class="card-content">
                            <span class="card-tag">Gestão Internacional (2024)</span>
                            <h3>Luanda – Angola</h3>
                            <p>Gestão de projetos culturais e intercâmbio artístico no exterior.</p>
                        </div>
                    </div>"""
new_prod = """                    <div class="card woodcut-border card--light">
                        <div class="card-content card-content--medium-padding">
                            <span class="card-tag card-tag--red">Produção e Projetos Especiais</span>
                            <h2 class="section-header">Gestão e Coordenação</h2>

                            <div class="atuacao-consolidated-grid">
                                <!-- COORDENAÇÃO DE FORMAÇÃO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Coordenação de Formação</p>
                                    <ul class="detail-list">
                                        <li><strong>Black Heroes (2018–2023)</strong> – Produtor Cultural e Coordenador de Formação</li>
                                        <li><strong>Azusa (2018–2023)</strong> – Coordenação de ensino, pesquisa e processos formativos</li>
                                    </ul>
                                </div>

                                <!-- GESTÃO INTERNACIONAL -->
                                <div class="atuacao-group">
                                    <p class="group-label">Gestão Internacional</p>
                                    <ul class="detail-list">
                                        <li><strong>Luanda, Angola (2024)</strong> – Gestão de projetos culturais e intercâmbio artístico no exterior</li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>"""
html = html.replace(old_prod, new_prod)

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)

print("HTML transformations applied successfully!")
