import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

# 1. Audiovisual - Rastros de Exu
# Find the start of the Audiovisual grid
av_grid_match = re.search(r'(<section id="audiovisual".*?<div class="catalog-grid[^>]*>)', html, re.DOTALL)
if av_grid_match:
    av_injection = r'''
                    <!-- RASTROS DE EXU -->
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
                            <div class="carousel-indicators">
                            </div>
                        </div>
                        <div class="card-content">
                            <span class="card-tag">Audiovisual Transmídia</span>
                            <h3 class="project-title">Rastros de Exu (2023)</h3>
                            <div class="project-context mt-20">
                                <p class="project-desc light">
                                    Série audiovisual e documentação poética derivada dos processos cênicos do espetáculo "Exu Não Vem Hoje". O projeto integrou as ações da Zona de Criação do Hub Cultural Porto Dragão no YouTube.
                                </p>
                            </div>
                        </div>
                    </div>
'''
    html = html.replace(av_grid_match.group(1), av_grid_match.group(1) + av_injection)

# 2. Atuação - Irreversível, Santo Bordel, Botas
# Find Atuação grid
at_grid_match = re.search(r'(<section id="atuacao".*?<div class="catalog-grid[^>]*>)', html, re.DOTALL)
if at_grid_match:
    at_injection = r'''
                    <!-- IRREVERSÍVEL -->
                    <div class="card woodcut-border">
                        <div class="carousel-container carousel-container--cols-2">
                            <div class="carousel-wrapper">
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-irreversivel/work-irreversivel-001.webp"
                                        alt="Irreversível"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-irreversivel/work-irreversivel-002.webp"
                                        alt="Irreversível Cena"></focus-image>
                                </div>
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>
                        <div class="card-content">
                            <span class="card-tag">Circulação Internacional</span>
                            <h3 class="project-title">Irreversível (2022)</h3>
                            <div class="project-context mt-20">
                                <p class="project-desc light">
                                    Atuação neste espetáculo integrando o percurso de práticas cênicas da Porto Iracema. A obra realizou apresentações em São Paulo e internacionalmente, em circulação por Angola.
                                </p>
                                <div class="project-details-mini mt-20">
                                    <strong>Ficha Técnica:</strong>
                                    <ul class="detail-list">
                                        <li><strong>Direção:</strong> Caique Melo</li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <!-- SANTO BORDEL E BOTAS -->
                    <div class="card woodcut-border">
                        <div class="carousel-container carousel-container--cols-3">
                            <div class="carousel-wrapper">
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-santo-bordel-de-tiatira/work-santo-bordel-de-tiatira-001.webp"
                                        alt="Santo Bordel"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-santo-bordel-de-tiatira/work-santo-bordel-de-tiatira-003.webp"
                                        alt="Santo Bordel"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-chega-de-falar-de-botas/work-chega-de-falar-de-botas-001.webp"
                                        alt="Chega de Falar de Botas"></focus-image>
                                </div>
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>
                        <div class="card-content">
                            <span class="card-tag">Atuação em Repertório</span>
                            <h3 class="project-title">Chega de Falar de Botas & Santo Bordel (2015-2017)</h3>
                            <div class="project-context mt-20">
                                <p class="project-desc light">
                                    Trabalhos que consolidaram os anos iniciais de atuação em palco sob direção de terceiros. <b>Santo Bordel de Tiatira</b> (Dir. Caique Melo, circulação em SP) e <b>Chega de Falar de Botas</b> (Dir. Andrei Bessa), nascida do laboratório de improvisação.
                                </p>
                            </div>
                        </div>
                    </div>
'''
    html = html.replace(at_grid_match.group(1), at_grid_match.group(1) + at_injection)

# 3. Produção - Trapo Preto
prod_grid_match = re.search(r'(<section id="producao".*?<div class="catalog-grid[^>]*>)', html, re.DOTALL)
if prod_grid_match:
    prod_injection = r'''
                    <!-- TRAPO PRETO -->
                    <div class="card woodcut-border">
                        <div class="carousel-container carousel-container--cols-3">
                            <div class="carousel-wrapper">
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-trapo-preto/work-trapo-preto-001.webp"
                                        alt="Trapo Preto"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-trapo-preto/work-trapo-preto-002.webp"
                                        alt="Trapo Preto Feira Negra"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-trapo-preto/work-trapo-preto-003.webp"
                                        alt="Trapo Preto Durag"></focus-image>
                                </div>
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>
                        <div class="card-content">
                            <span class="card-tag">Afroempreendedorismo e Cultura</span>
                            <h3 class="project-title">Trapo Preto (2018)</h3>
                            <div class="project-context mt-20">
                                <p class="project-desc light">
                                    Empreendedor e proprietário da marca <b>Trapo Preto</b>, voltada para a comercialização e valorização cultural das durags. Através da marca, Semino atua ativamente na Feira Negra e em ecossistemas de realizadores afro-cearenses.
                                </p>
                            </div>
                        </div>
                    </div>
'''
    html = html.replace(prod_grid_match.group(1), prod_grid_match.group(1) + prod_injection)

# 4. Ensino - Angola (Teatro Voluntário)
ensino_grid_match = re.search(r'(<section id="ensino".*?<div class="catalog-grid[^>]*>)', html, re.DOTALL)
if ensino_grid_match:
    ensino_injection = r'''
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
'''
    html = html.replace(ensino_grid_match.group(1), ensino_grid_match.group(1) + ensino_injection)

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)

print("Curatorial injections applied to HTML.")
