/**
 * MetadataReader
 * Responsável por extrair o ponto de foco (XMP) de arquivos JPG e PNG.
 */
class MetadataReader {
    static async read(src) {
        try {
            const response = await fetch(src);
            if (!response.ok) return null;
            const buffer = await response.arrayBuffer();
            const view = new DataView(buffer);

            // Tenta extrair XMP baseado no tipo de arquivo
            let xmpString = "";
            if (this.isJPEG(view)) {
                xmpString = this.extractXMPFromJPEG(view);
            } else if (this.isPNG(view)) {
                xmpString = this.extractXMPFromPNG(view);
            }

            if (!xmpString) return null;
            return this.parseXMPFocus(xmpString);
        } catch (e) {
            console.warn("MetadataReader error:", e);
            return null;
        }
    }

    static isJPEG(view) {
        return view.getUint16(0) === 0xFFD8;
    }

    static isPNG(view) {
        return view.getUint32(0) === 0x89504E47;
    }

    static extractXMPFromJPEG(view) {
        let offset = 2;
        while (offset < view.byteLength) {
            const marker = view.getUint16(offset);
            const length = view.getUint16(offset + 2);

            // APP1 Marker (XMP)
            if (marker === 0xFFE1) {
                const identifier = new TextDecoder().decode(new Uint8Array(view.buffer, offset + 4, 29));
                if (identifier.startsWith("http://ns.adobe.com/xap/1.0/")) {
                    return new TextDecoder().decode(new Uint8Array(view.buffer, offset + 4 + 29, length - 4 - 29));
                }
            }
            offset += 2 + length;
        }
        return "";
    }

    static extractXMPFromPNG(view) {
        let offset = 8; // Skip PNG header
        while (offset < view.byteLength) {
            const length = view.getUint32(offset);
            const type = new TextDecoder().decode(new Uint8Array(view.buffer, offset + 4, 4));

            if (type === "iTXt") {
                const data = new Uint8Array(view.buffer, offset + 8, length);
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
        return ["src", "fit", "fallback", "debug"];
    }

    constructor() {
        super();
        this.attachShadow({ mode: "open" });
    }

    connectedCallback() {
        this.render();
    }

    attributeChangedCallback() {
        this.render();
    }

    async render() {
        const src = this.getAttribute("src");
        if (!src) return;

        const fit = this.getAttribute("fit") || "cover";
        const fallback = this.getAttribute("fallback") || "center center";
        const debug = this.hasAttribute("debug");

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
            <img src="${src}" style="object-fit: ${fit}; object-position: ${fallback};">
        `;

        // Busca metadados
        const focus = await MetadataReader.read(src);
        const img = this.shadowRoot.querySelector("img");

        let xPct, yPct;
        if (focus) {
            xPct = (focus.x * 100).toFixed(2) + "%";
            yPct = (focus.y * 100).toFixed(2) + "%";
            img.style.objectPosition = `${xPct} ${yPct}`;
        } else {
            img.style.objectPosition = fallback;
            // Best effort to parse fallback for debug marker
            xPct = "50%"; yPct = "50%";
            if (fallback.includes("%")) {
                const parts = fallback.split(" ");
                xPct = parts[0] || "50%";
                yPct = parts[1] || "50%";
            }
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
    }
}

// Registra o componente
customElements.define("focus-image", FocusImage);
