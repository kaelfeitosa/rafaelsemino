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
        
        // Disable timeout for this page to avoid navigation/render timeouts
        await page.setDefaultNavigationTimeout(0);
        await page.setDefaultTimeout(0);
        
        // Define screen viewport (1080p)
        const width = 1920;
        await page.setViewport({ width, height: 1080 }); // Initial HD viewport
        
        // Use the absolute path to portfolio.html
        const filePath = path.resolve(__dirname, '../../frontend/index.html');
        // Construct a file URL
        const fileUrl = `file:///${filePath.replace(/\\/g, '/')}`;
        console.log(`Loading page: ${fileUrl}`);
        
        await page.goto(fileUrl, { waitUntil: 'networkidle0' });
        
        console.log('Page loaded. Adjusting viewport logic...');
        
        // Emulate screen media type to get exact screen styles (ignoring @media print)
        await page.emulateMediaType('screen');
        
        // The problem with page.pdf is that the paper height becomes the new layout viewport.
        // So 100vh becomes the entire document length instead of 1080px.
        // We fix this by hardcoding pixel heights for the 1080p elements BEFORE calculating total height.
        await page.evaluate(() => {
            // Fix hero section
            const hero = document.getElementById('home');
            if (hero) {
                hero.style.setProperty('min-height', '918px', 'important');
                hero.classList.remove('min-h-[85vh]');
            }
            
            // Fix root wrapper
            const wrapper = document.querySelector('.min-h-screen');
            if (wrapper) {
                wrapper.style.setProperty('min-height', '1080px', 'important');
                wrapper.classList.remove('min-h-screen');
            }
            
            // PDF-specific fixes
            const style = document.createElement('style');
            style.innerHTML = `
                /* Bypass IntersectionObserver opacity/transform */
                section { 
                    opacity: 1 !important; 
                    transform: none !important; 
                    transition: none !important;
                }
                
                /* Fix white/transparent menu */
                header { 
                    background: #ffffff !important; 
                    backdrop-filter: none !important;
                    position: relative !important; /* Avoid sticky bugs in PDF */
                    border-bottom: 3px solid #000 !important;
                }

                /* Disable grain overlay (saves size and avoids rendering artifacts) */
                .noise { display: none !important; }

                /* Fix section labels visibility */
                .section-label {
                    opacity: 1 !important;
                    color: #000000 !important;
                    -webkit-text-stroke: 0 !important;
                    writing-mode: horizontal-tb !important; /* Safer for PDF */
                    transform: none !important;
                    position: relative !important;
                    top: 0 !important;
                    margin: 20px 0 !important;
                    padding: 10px 40px !important;
                    font-size: 4rem !important;
                    background: #f2f2f2 !important;
                    width: 100% !important;
                    display: block !important;
                }

                .min-h-\\[85vh\\] { min-height: 918px !important; }
                .min-h-screen { min-height: 1080px !important; }
                .h-screen { height: 1080px !important; }
            `;
            document.head.appendChild(style);
        });
        
        // Let it settle any layout changes
        await new Promise(r => setTimeout(r, 2000));
        
        // Evaluate the full height of the body content
        const bodyHeight = await page.evaluate(() => {
            return document.documentElement.scrollHeight;
        });
        
        console.log(`Calculated height (scrollHeight): ${bodyHeight}px`);
        
        // Set viewport to the full height exactly so it renders completely
        await page.setViewport({ width, height: bodyHeight });
        
        // Wait another moment for any deferred renders
        await new Promise(r => setTimeout(r, 2000));
        
        // Generate PDF
        const pdfPath = path.resolve(__dirname, '../../frontend/index.pdf');
        console.log(`Generating PDF at ${pdfPath}...`);
        
        await page.pdf({
            path: pdfPath,
            width: `${width}px`,
            height: `${bodyHeight}px`,
            printBackground: true,
            pageRanges: '1', // Force a single page
            timeout: 0 // Disable PDF generation timeout
        });
        
        console.log('PDF generation complete. File saved to:', pdfPath);
        await browser.close();
    } catch (e) {
        console.error('Error generating PDF:', e);
        process.exit(1);
    }
})();
