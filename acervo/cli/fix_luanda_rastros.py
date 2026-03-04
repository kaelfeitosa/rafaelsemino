import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

# 1. SWAP ABARCA AND LUANDA AND CONVERT LUANDA TO CAROUSEL
# Find the Ensino section body
target_ensino = """            <section id="ensino" tabindex="-1">
                <div class="section-label stamped">ENSINO</div>
                <div class="catalog-grid catalog-grid--single-col">

                    <div class="highlight-wrapper highlight-wrapper--compact highlight-wrapper--ensino">
                        <!-- LUANDA CARD -->
                        <div class="card woodcut-border card--light">
                            <div class="img-stack img-stack--cols-2">
                                <focus-image
                                    src="images/optimized/work-rafael-ensino-angola-2018/work-rafael-ensino-angola-2018-001.webp"
                                    alt="Angola Ensino"></focus-image>
                                <focus-image
                                    src="images/optimized/work-rafael-ensino-angola-2018/work-rafael-ensino-angola-2018-002.webp"
                                    alt="Angola Teatro"></focus-image>
                            </div>
                            <div class="card-content">
                                <span class="card-tag">Intercâmbio e Voluntariado</span>
                                <h3 class="project-title">Ensino de Teatro em Luanda (2018)</h3>
                                <div class="project-context mt-20">
                                    <p class="project-desc light">
                                        Atuação voluntária realizada durante o período de graduação. O projeto consistiu
                                        na facilitação de práticas em teatro e intercâmbio de experiências pedagógicas
                                        com alunos em Luanda, Angola, marcando profundamente sua reflexão como educador.
                                    </p>
                                </div>
                            </div>
                        </div>

                        <!-- ABARCA E PERCURSO BÁSICO -->
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
                                <span class="card-tag card-tag--red">Escola Porto Iracema das Artes 📍 Fortaleza – CE –
                                    Brasil</span>
                                <h3 class="course-institution">Projeto Abarca e Percurso Básico de Teatro</h3>
                                <div class="mt-20">
                                    <p class="course-details"><strong>Atuação:</strong> Professor Formador (2023).</p>
                                    <p class="course-details mt-10">Desenvolvimento de processos pedagógicos voltados à
                                        formação artística em diferentes territórios de Fortaleza, incluindo Itapipoca,
                                        Vicente Pinzón e Genibaú, bem como na sede do Porto Iracema das Artes.</p>
                                </div>
                            </div>
                        </div>"""

replacement_ensino = """            <section id="ensino" tabindex="-1">
                <div class="section-label stamped">ENSINO</div>
                <div class="catalog-grid catalog-grid--single-col">

                    <div class="highlight-wrapper highlight-wrapper--compact highlight-wrapper--ensino">
                        
                        <!-- ABARCA E PERCURSO BÁSICO -->
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
                                <span class="card-tag card-tag--red">Escola Porto Iracema das Artes 📍 Fortaleza – CE –
                                    Brasil</span>
                                <h3 class="course-institution">Projeto Abarca e Percurso Básico de Teatro</h3>
                                <div class="mt-20">
                                    <p class="course-details"><strong>Atuação:</strong> Professor Formador (2023).</p>
                                    <p class="course-details mt-10">Desenvolvimento de processos pedagógicos voltados à
                                        formação artística em diferentes territórios de Fortaleza, incluindo Itapipoca,
                                        Vicente Pinzón e Genibaú, bem como na sede do Porto Iracema das Artes.</p>
                                </div>
                            </div>
                        </div>
                        
                        <!-- LUANDA CARD -->
                        <div class="card woodcut-border card--light">
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
                                <div class="carousel-indicators"></div>
                            </div>
                            <div class="card-content">
                                <span class="card-tag">Intercâmbio e Voluntariado</span>
                                <h3 class="project-title">Ensino de Teatro em Luanda (2018)</h3>
                                <div class="project-context mt-20">
                                    <p class="project-desc light">
                                        Atuação voluntária realizada durante o período de graduação. O projeto consistiu
                                        na facilitação de práticas em teatro e intercâmbio de experiências pedagógicas
                                        com alunos em Luanda, Angola, marcando profundamente sua reflexão como educador.
                                    </p>
                                </div>
                            </div>
                        </div>"""

html = html.replace(target_ensino, replacement_ensino)

# 2. INJECT RASTROS DE EXU YOUTUBE LINK
target_rastros = """                        <!-- RASTROS DE EXU -->
                        <div class="card woodcut-border">
                            <div class="carousel-container carousel-container--cols-3">
                                <div class="carousel-wrapper">
                                    <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rastros-de-exu/work-rastros-de-exu-001.webp"
                                        alt="Rastros de Exu"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rastros-de-exu/work-rastros-de-exu-002.webp"
                                        alt="Rastros de Exu - Bastidores"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rastros-de-exu/work-rastros-de-exu-003.webp"
                                        alt="Rastros de Exu - Detalhe"></focus-image>
                                    </div>
                                </div>
                                <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                                <button class="carousel-btn carousel-btn--next">&#10095;</button>
                                <div class="carousel-indicators"></div>
                            </div>
                            <div class="card-content">
                                <span class="card-tag">Performance / Atuação</span>
                                <h3 class="project-title">Rastros de Exu (2023)</h3>
                                <div class="project-context mt-20">
                                    <p class="project-desc light">
                                        Ação performática derivada do encontro de laboratórios no CCBJ, trabalhando
                                        processos de transe e improvisação cênica.
                                    </p>
                                </div>
                            </div>
                        </div>"""
                        
replacement_rastros = """                        <!-- RASTROS DE EXU -->
                        <div class="card woodcut-border">
                            <div class="carousel-container carousel-container--cols-3">
                                <div class="carousel-wrapper">
                                    <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rastros-de-exu/work-rastros-de-exu-001.webp"
                                        alt="Rastros de Exu"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rastros-de-exu/work-rastros-de-exu-002.webp"
                                        alt="Rastros de Exu - Bastidores"></focus-image>
                                    </div>
                                    <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-rastros-de-exu/work-rastros-de-exu-003.webp"
                                        alt="Rastros de Exu - Detalhe"></focus-image>
                                    </div>
                                </div>
                                <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                                <button class="carousel-btn carousel-btn--next">&#10095;</button>
                                <div class="carousel-indicators"></div>
                            </div>
                            <div class="card-content">
                                <span class="card-tag">Performance / Atuação</span>
                                <h3 class="project-title">Rastros de Exu (2023)</h3>
                                <div class="project-context mt-20">
                                    <p class="project-desc light">
                                        Ação performática derivada do encontro de laboratórios no CCBJ, trabalhando
                                        processos de transe e improvisação cênica.
                                    </p>
                                    <p class="project-desc light" style="margin-top: 15px;">
                                        <a href="https://www.youtube.com/watch?v=bjwQwDCBGsI&t=3s" target="_blank" class="catalog-link">▶ Assista à Performance</a>
                                    </p>
                                </div>
                            </div>
                        </div>"""
                        
html = html.replace(target_rastros, replacement_rastros)

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)
print("Changes applied!")
