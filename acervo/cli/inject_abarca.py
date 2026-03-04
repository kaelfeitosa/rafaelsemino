import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

# 1. REMOVE ABARCA FROM TEXT LIST
target_list = """                                    <!-- MEDIAÇÃO E CURSOS LIVRES -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Mediação e Cursos Livres</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor (Projeto ABARCA e Percurso Básico de Teatro)</strong> – Porto Iracema
                                                das Artes</li>
                                            <li><strong>Professor</strong> – Programa de Aceleração do Idoso</li>
                                            <li><strong>Orientador de Montagem e Mediação</strong> – Projeto no Bairro
                                                Terceiro</li>
                                        </ul>
                                    </div>"""
                                    
replacement_list = """                                    <!-- MEDIAÇÃO E CURSOS LIVRES -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Mediação e Cursos Livres</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor</strong> – Programa de Aceleração do Idoso</li>
                                            <li><strong>Orientador de Montagem e Mediação</strong> – Projeto no Bairro
                                                Terceiro</li>
                                        </ul>
                                    </div>"""
html = html.replace(target_list, replacement_list)

# 2. INSERT ABARCA CARD ABOVE HISTORICO PEDAGOGICO
# Look for the start of the historico pedagogico card to insert before it.
target_insert = """                        <!-- HISTÓRICO PEDAGÓGICO CONSOLIDADO -->"""

abarca_card = """                        <!-- ABARCA E PERCURSO BÁSICO -->
                        <div class="card woodcut-border card--light">
                            <div class="carousel-container carousel-container--cols-3">
                                <div class="carousel-wrapper">
                                    <div class="carousel-slide"><focus-image
                                            src="images/optimized/work-prof-percurso-basico/work-prof-percurso-basico-001.webp"
                                            alt="Aula do Projeto ABARCA e Percurso Básico"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                            src="images/optimized/work-prof-percurso-basico/work-prof-percurso-basico-002.webp"
                                            alt="Atividade prática com alunos"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                            src="images/optimized/work-prof-percurso-basico/work-prof-percurso-basico-003.webp"
                                            alt="Formação artística"></focus-image>
                                    </div>
                                </div>
                                <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                                <button class="carousel-btn carousel-btn--next">&#10095;</button>
                                <div class="carousel-indicators"></div>
                            </div>
                            <div class="card-content">
                                <span class="card-tag card-tag--red">Escola Porto Iracema das Artes 📍 Fortaleza – CE – Brasil</span>
                                <h3 class="course-institution">Projeto Abarca e Percurso Básico de Teatro</h3>
                                <div class="mt-20">
                                    <p class="course-details"><strong>Atuação:</strong> Professor Formador (2023).</p>
                                    <p class="course-details mt-10">Desenvolvimento de processos pedagógicos voltados à formação artística em diferentes territórios de Fortaleza, incluindo Itapipoca, Vicente Pinzón e Genibaú, bem como na sede do Porto Iracema das Artes.</p>
                                </div>
                            </div>
                        </div>

                        <!-- HISTÓRICO PEDAGÓGICO CONSOLIDADO -->"""
                        
html = html.replace(target_insert, abarca_card)


with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)
print("ABARCA card injected successfully!")
