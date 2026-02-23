/**
 * MetadataReader
 * Responsável por extrair o ponto de foco (XMP) de arquivos JPG e PNG.
 * Optimized for partial content (Range requests).
 */
class MetadataReader {
    static async read(src, signal) {
        try {
            // Tenta buscar apenas os primeiros 64KB (cabeçalho + metadados geralmente estão no início)
            const headers = { "Range": "bytes=0-65535" };
            const response = await fetch(src, { headers, signal });

            if (!response.ok && response.status !== 206) return null;

            const buffer = await response.arrayBuffer();
            const view = new DataView(buffer);

            // Tenta extrair XMP baseado no tipo de arquivo
            let xmpString = "";
            if (MetadataReader.isJPEG(view)) {
                xmpString = MetadataReader.extractXMPFromJPEG(view);
            } else if (MetadataReader.isPNG(view)) {
                xmpString = MetadataReader.extractXMPFromPNG(view);
            }

            if (!xmpString) return null;
            return MetadataReader.parseXMPFocus(xmpString);
        } catch (e) {
            if (e.name !== 'AbortError') {
                console.warn("MetadataReader error:", e);
            }
            return null;
        }
    }

    static isJPEG(view) {
        // Safety check: ensure buffer has at least 2 bytes for SOI marker
        if (view.byteLength < 2) return false;
        return view.getUint16(0) === 0xFFD8;
    }

    static isPNG(view) {
        if (view.byteLength < 8) return false;
        return view.getUint32(0) === 0x89504E47;
    }

    static extractXMPFromJPEG(view) {
        let offset = 2;
        // Ensure enough bytes for marker (2) and length (2) to prevent out-of-bounds reads
        while (offset + 4 <= view.byteLength) {
            const marker = view.getUint16(offset);
            const length = view.getUint16(offset + 2);

            // Bounds check for segment data
            if (offset + 2 + length > view.byteLength) break;

            // APP1 Marker (XMP)
            if (marker === 0xFFE1) {
                // Check if we have enough bytes for the identifier "http://ns.adobe.com/xap/1.0/" (29 bytes)
                // Safety check: ensure segment length is sufficient to contain the identifier
                if (length >= 2 + 29) {
                    const identifier = new TextDecoder().decode(new Uint8Array(view.buffer, offset + 4, 29));
                    if (identifier.startsWith("http://ns.adobe.com/xap/1.0/")) {
                        // Extract XMP packet: length - 2 (length bytes) - 29 (identifier)
                        return new TextDecoder().decode(new Uint8Array(view.buffer, offset + 4 + 29, length - 2 - 29));
                    }
                }
            }
            offset += 2 + length;
        }
        return "";
    }

    static extractXMPFromPNG(view) {
        let offset = 8; // Skip PNG header
        // Ensure enough bytes for length (4) and type (4) to prevent out-of-bounds reads
        while (offset + 8 <= view.byteLength) {
            const length = view.getUint32(offset);
            const type = new TextDecoder().decode(new Uint8Array(view.buffer, offset + 4, 4));

            // Bounds check for chunk data + CRC (4 bytes)
            if (offset + 8 + length + 4 > view.byteLength) break;

            if (type === "iTXt") {
                const data = new Uint8Array(view.buffer, offset + 8, length);
                // iTXt structure: Keyword (null-terminated), compression flag, compression method, language tag (null-term), translated keyword (null-term), text
                // Simple search for XML packet in the data
                const text = new TextDecoder().decode(data);
                if (text.includes("XML:com.adobe.xmp")) {
                    const xmpStart = text.indexOf("<x:xmpmeta");
                    if (xmpStart !== -1) return text.substring(xmpStart);
                }
            }
            offset += 8 + length + 4; // Length + Type + Data + CRC
        }
        return "";
    }

    static parseXMPFocus(xmpString) {
        const parser = new DOMParser();
        const xmlDoc = parser.parseFromString(xmpString, "text/xml");

        // Busca stArea:x e stArea:y (pode vir com namespace ou não dependendo do parser)
        const getTagValue = (tagName) => {
            const el = xmlDoc.getElementsByTagName(tagName)[0] ||
                xmlDoc.getElementsByTagNameNS("*", tagName)[0];
            return el ? parseFloat(el.textContent) : null;
        };

        const x = getTagValue("x");
        const y = getTagValue("y");

        if (x !== null && y !== null) {
            return {
                x: Math.max(0, Math.min(1, x)),
                y: Math.max(0, Math.min(1, y))
            };
        }
        return null;
    }
}

/**
 * <focus-image> Web Component
 */
class FocusImage extends HTMLElement {
    static get observedAttributes() {
        return ["src", "fit", "fallback", "debug", "alt", "loading"];
    }

    constructor() {
        super();
        this.attachShadow({ mode: "open" });
        this.abortController = null;
    }

    static parseFallbackPosition(pos) {
        if (!pos) return { x: "50%", y: "50%" };

        const parts = pos.trim().split(/\s+/);
        let x = null;
        let y = null;

        parts.forEach(p => {
            if (p === "left") x = "0%";
            else if (p === "right") x = "100%";
            else if (p === "top") y = "0%";
            else if (p === "bottom") y = "100%";
            else if (p === "center") {
                if (x === null) x = "50%";
                else y = "50%";
            } else {
                // Assume percentage or length value
                if (x === null) x = p;
                else y = p;
            }
        });

        return { x: x || "50%", y: y || "50%" };
    }

    connectedCallback() {
        this.render();
    }

    disconnectedCallback() {
        if (this.abortController) {
            this.abortController.abort();
        }
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (oldValue !== newValue) {
            this.render();
        }
    }

    async render() {
        // Cancel previous fetch if any to avoid race conditions and ensure only the latest request is processed
        if (this.abortController) {
            this.abortController.abort();
        }
        this.abortController = new AbortController();
        const signal = this.abortController.signal;

        const src = this.getAttribute("src");
        if (!src) return;

        const fit = this.getAttribute("fit") || "cover";
        const fallback = this.getAttribute("fallback") || "center center";
        const debug = this.hasAttribute("debug");
        const alt = this.getAttribute("alt") || "";
        const loading = this.getAttribute("loading") || "lazy";

        // UI Inicial
        this.shadowRoot.innerHTML = `
            <style>
                :host {
                    display: block;
                    width: 100%;
                    height: 100%;
                    overflow: hidden;
                    position: relative;
                }
                img {
                    width: 100%;
                    height: 100%;
                    display: block;
                    transition: object-position 0.3s ease;
                }
                .marker {
                    position: absolute;
                    width: 20px;
                    height: 20px;
                    border: 2px solid #e60000;
                    border-radius: 50%;
                    transform: translate(-50%, -50%);
                    pointer-events: none;
                    z-index: 10;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    color: #e60000;
                    font-weight: bold;
                    font-size: 14px;
                    background: rgba(255,255,255,0.5);
                }
                .marker::after {
                    content: '⊕';
                }
            </style>
            <img part="image" src="${src}" alt="${alt}" loading="${loading}" style="object-fit: ${fit}; object-position: ${fallback};">
        `;

        const img = this.shadowRoot.querySelector("img");

        try {
            // Busca metadados
            const focus = await MetadataReader.read(src, signal);
            if (signal.aborted) return;

            let xPct, yPct;
            if (focus) {
                xPct = (focus.x * 100).toFixed(2) + "%";
                yPct = (focus.y * 100).toFixed(2) + "%";
                img.style.objectPosition = `${xPct} ${yPct}`;
            } else {
                // If metadata fetch failed or no focus found, keep fallback
                // Calculate debug marker position based on fallback
                const parsed = FocusImage.parseFallbackPosition(fallback);
                xPct = parsed.x;
                yPct = parsed.y;
            }

            if (debug) {
                const marker = document.createElement("div");
                marker.className = "marker";
                marker.style.left = xPct;
                marker.style.top = yPct;
                if (!focus) {
                    // Dim the fallback marker differently
                    marker.style.borderColor = "#ff9900";
                    marker.style.color = "#ff9900";
                    marker.title = "Fallback Focus";
                }
                this.shadowRoot.appendChild(marker);
            }
        } catch (error) {
            // Ignore abort errors
        }
    }
}

// Registra o componente
customElements.define("focus-image", FocusImage);
