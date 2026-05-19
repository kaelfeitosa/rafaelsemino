const puppeteer = require('puppeteer');
const path = require('path');

(async () => {
    try {
        console.log('Launching browser...');
        const browser = await puppeteer.launch({
            args: ['--disable-gpu', '--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage'],
            headless: true
        });
        const page = await browser.newPage();

        // Disable timeout
        await page.setDefaultNavigationTimeout(0);
        await page.setDefaultTimeout(0);

        // Define screen viewport (HD)
        const width = 1920;
        await page.setViewport({ width, height: 1080 });

        // Path to index.html
        const filePath = path.resolve(__dirname, '../../frontend/index.html');
        const fileUrl = `file:///${filePath.replace(/\\/g, '/')}`;
        console.log(`Loading page: ${fileUrl}`);

        await page.goto(fileUrl, { waitUntil: 'networkidle0' });

        console.log('Page loaded. Applying PDF fixes...');

        // Emulate screen media to keep screen styles, not print styles
        await page.emulateMediaType('screen');

        await page.evaluate(() => {
            // ── 1. Hide original section labels completely ──
            document.querySelectorAll('.section-label').forEach(el => {
                el.style.display = 'none';
            });

            // ── 2. Inject new visible heading before each section's content ──
            document.querySelectorAll('.section-label').forEach(el => {
                const text = el.textContent.trim();
                const banner = document.createElement('h2');
                banner.textContent = text;
                banner.setAttribute('style', [
                    'font-family: "Archivo Black", sans-serif',
                    'font-size: 2.5rem',
                    'color: #000000',
                    'background: #f2f2f2',
                    'padding: 15px 40px',
                    'margin: 0 0 30px 0',
                    'border-bottom: 4px solid #e60000',
                    'text-transform: uppercase',
                    'display: block',
                    'width: 100%',
                    'line-height: 1.2',
                    'letter-spacing: 2px',
                ].join('; '));

                // Insert the banner right after the hidden label
                el.parentNode.insertBefore(banner, el.nextSibling);
            });

            // ── 3. Inject global CSS fixes ──
            const style = document.createElement('style');
            style.innerHTML = `
                /* Kill grain overlay */
                .noise { display: none !important; }

                /* Kill scroll animations */
                * {
                    transition: none !important;
                    animation: none !important;
                }

                /* Force all sections visible */
                section {
                    opacity: 1 !important;
                    transform: none !important;
                }

                /* Fix sticky nav */
                nav, header {
                    position: relative !important;
                    backdrop-filter: none !important;
                }

                /* Fix hero vh */
                .hero-section, #home {
                    min-height: 1080px !important;
                    height: auto !important;
                }

                /* Sections: collapse 2-col grid to single column */
                section {
                    grid-template-columns: 1fr !important;
                    gap: 0 !important;
                }

                /* Carousels: show all slides side by side */
                .carousel-container {
                    height: 400px !important;
                }
                .carousel-wrapper {
                    transform: none !important;
                    display: flex !important;
                    gap: 2px !important;
                }
                .carousel-slide {
                    flex: 0 0 calc((100% - 4px) / 3) !important;
                    height: 100% !important;
                }
                .carousel-btn, .carousel-indicators {
                    display: none !important;
                }

                /* Compact card carousels: single slide */
                .catalog-grid--compact .carousel-container {
                    height: 300px !important;
                }
                .catalog-grid--compact .carousel-slide {
                    flex: 0 0 100% !important;
                }

                /* Hide intro overlay */
                .intro-overlay { display: none !important; }
            `;
            document.head.appendChild(style);
        });

        // Debug screenshot
        console.log('Taking debug screenshot...');
        await page.screenshot({
            path: path.resolve(__dirname, '../../frontend/rafael_pdf_debug.png'),
            fullPage: true
        });
        console.log('Debug screenshot saved.');

        const bodyHeight = await page.evaluate(() => document.documentElement.scrollHeight);
        console.log(`Calculated height: ${bodyHeight}px`);

        await page.setViewport({ width, height: bodyHeight });
        await new Promise(r => setTimeout(r, 2000));

        const pdfPath = path.resolve(__dirname, '../../frontend/index.pdf');
        console.log(`Generating PDF at ${pdfPath}...`);

        await page.pdf({
            path: pdfPath,
            width: `${width}px`,
            height: `${bodyHeight}px`,
            printBackground: true,
            pageRanges: '1',
            timeout: 0
        });

        console.log('PDF generation complete:', pdfPath);
        await browser.close();
    } catch (e) {
        console.error('Error:', e);
        process.exit(1);
    }
})();
