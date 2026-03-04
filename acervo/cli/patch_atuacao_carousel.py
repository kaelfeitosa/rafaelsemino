import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

target = """                <div class="catalog-grid catalog-grid--single-col gap-60">
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
                                    Atuação neste espetáculo integrando o percurso de práticas cênicas da Porto Iracema.
                                    A obra realizou apresentações em São Paulo e internacionalmente, em circulação por
                                    Angola.
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
                                    Trabalhos que consolidaram os anos iniciais de atuação em palco sob direção de
                                    terceiros. <b>Santo Bordel de Tiatira</b> (Dir. Caique Melo, circulação em SP) e
                                    <b>Chega de Falar de Botas</b> (Dir. Andrei Bessa), nascida do laboratório de
                                    improvisação.
                                </p>
                            </div>
                        </div>
                    </div>

                    <div class="card woodcut-border card--light">

                        <div class="card-content card-content--medium-padding">
                            <span class="card-tag card-tag--red">Repertório e Performance</span>
                            <h2 class="section-header">Atuação e Criação</h2>

                            <div class="atuacao-consolidated-grid">
                                <!-- TEATRO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Teatro (Destaques)</p>
                                    <ul class="detail-list">
                                        <li><strong>Irreversível</strong> (Dir: Caique Melo)</li>
                                        <li><strong>Santo Bordel de Tiatira</strong> (Dir: Caique Melo)</li>
                                        <li><strong>De Louco, Todo Mundo Tem um Pouco</strong> (Dir: Hiroldo Serra)</li>
                                        <li><strong>A Serpente</strong> (Dir: Maria Vitória)</li>
                                        <li><strong>De Sucupira à Asa Branca</strong> (Dir: Fernando Lira)</li>
                                    </ul>
                                </div>

                                <!-- AUTORIA -->
                                <div class="atuacao-group">
                                    <p class="group-label">Autoria e Escrita Cênica</p>
                                    <ul class="detail-list">
                                        <li><strong>Cala-me os Olhos</strong> (2016)</li>
                                        <li><strong>Sociedade o Circo</strong> (2017)</li>
                                        <li><strong>O Caso da Casa</strong> (2019)</li>
                                    </ul>
                                </div>

                                <!-- AUDIOVISUAL -->
                                <div class="atuacao-group">
                                    <p class="group-label">Audiovisual</p>
                                    <ul class="detail-list">
                                        <li><strong>Astronauta</strong> (Produção Audiovisual – Angola)</li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>"""

replacement = """                <div class="catalog-grid catalog-grid--single-col gap-60">
                    <div class="card woodcut-border card--light">
                        <div class="carousel-container carousel-container--cols-2">
                            <div class="carousel-wrapper">
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-irreversivel/work-irreversivel-001.webp"
                                        alt="Irreversível"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-santo-bordel-de-tiatira/work-santo-bordel-de-tiatira-001.webp"
                                        alt="Santo Bordel de Tiatira"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-chega-de-falar-de-botas/work-chega-de-falar-de-botas-001.webp"
                                        alt="Chega de Falar de Botas"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-de-louco-todo-mundo-tem-um-pouco/work-de-louco-todo-mundo-tem-um-pouco-001.webp"
                                        alt="De Louco, Todo Mundo Tem Um Pouco"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-de-sucupira-a-asa-branca/work-de-sucupira-a-asa-branca-001.webp"
                                        alt="De Sucupira à Asa Branca"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-a-serpente/work-a-serpente-001.webp"
                                        alt="A Serpente"></focus-image>
                                </div>
                                <div class="carousel-slide"><focus-image
                                        src="images/optimized/work-sociedade-o-circo/work-sociedade-o-circo-001.webp"
                                        alt="Sociedade o Circo"></focus-image>
                                </div>
                            </div>
                            <button class="carousel-btn carousel-btn--prev">&#10094;</button>
                            <button class="carousel-btn carousel-btn--next">&#10095;</button>
                            <div class="carousel-indicators">
                            </div>
                        </div>
                        <div class="card-content card-content--medium-padding">
                            <span class="card-tag card-tag--red">Repertório e Performance</span>
                            <h2 class="section-header">Atuação e Criação Cênica</h2>

                            <div class="atuacao-consolidated-grid">
                                <!-- TEATRO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Teatro (Destaques)</p>
                                    <ul class="detail-list">
                                        <li><strong>Irreversível</strong> (2022) – Circulação Secult/Angola – Dir: Caique Melo</li>
                                        <li><strong>Santo Bordel de Tiatira</strong> (2017) – Dir: Caique Melo</li>
                                        <li><strong>De Louco, Todo Mundo Tem um Pouco</strong> (2017) – Dir: Hiroldo Serra</li>
                                        <li><strong>Chega de Falar de Botas</strong> (2015) – Dir: Andrei Bessa</li>
                                        <li><strong>De Sucupira à Asa Branca</strong> (2016) – Dir: Fernando Lira</li>
                                        <li><strong>A Serpente</strong> (2014) – Dir: Maria Vitória</li>
                                    </ul>
                                </div>

                                <!-- AUTORIA -->
                                <div class="atuacao-group">
                                    <p class="group-label">Autoria e Escrita Cênica</p>
                                    <ul class="detail-list">
                                        <li><strong>O Caso da Casa</strong> (2019)</li>
                                        <li><strong>Sociedade o Circo</strong> (2017)</li>
                                        <li><strong>Cala-me os Olhos</strong> (2016)</li>
                                    </ul>
                                </div>

                                <!-- AUDIOVISUAL -->
                                <div class="atuacao-group">
                                    <p class="group-label">Audiovisual Transversal</p>
                                    <ul class="detail-list">
                                        <li><strong>Astronauta</strong> (Produção Audiovisual – Angola)</li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>"""

if target in html:
    html = html.replace(target, replacement)
    with open(html_path, 'w', encoding='utf-8') as f:
        f.write(html)
    print("Super Carousel successfully injected.")
else:
    print("Could not find exact target block.")
