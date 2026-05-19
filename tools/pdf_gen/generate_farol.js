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
        
        // Path to farol-novo.html
        const filePath = path.resolve(__dirname, '../../frontend/farol-novo.html');
        const fileUrl = `file:///${filePath.replace(/\\/g, '/')}`;
        console.log(`Loading page: ${fileUrl}`);
        
        await page.goto(fileUrl, { waitUntil: 'networkidle0' });
        
        console.log('Page loaded. Adjusting styles for PDF...');
        
        // Emulate screen media
        await page.emulateMediaType('screen');
        
        await page.evaluate(() => {
            // PDF-specific fixes for farol-novo.html
            const style = document.createElement('style');
            style.innerHTML = `
                .noise { display: none !important; }

                /* Ensure carousels show their content correctly */
                .carousel-container {
                    height: 480px !important;
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
                .intro-overlay { display: none !important; }
            `;
            document.head.appendChild(style);
        });
        
        await new Promise(r => setTimeout(r, 2000));
        
        // Final debug screenshot to see what's happening
        await page.screenshot({ path: path.resolve(__dirname, '../../frontend/pdf_debug.png'), fullPage: true });
        
        const bodyHeight = await page.evaluate(() => {
            return document.documentElement.scrollHeight;
        });
        
        console.log(`Calculated height: ${bodyHeight}px`);
        
        await page.setViewport({ width, height: bodyHeight });
        await new Promise(r => setTimeout(r, 2000));
        
        const pdfPath = path.resolve(__dirname, '../../frontend/farol-novo.pdf');
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
