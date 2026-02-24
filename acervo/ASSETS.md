# Asset Optimization Pipeline

This project uses a dedicated pipeline to optimize images referenced in the frontend.

## Overview

Instead of committing optimized images manually, we commit **source masters** (high quality JPEGs/PNGs) in `acervo/media/images`. The build process then generates optimized WebP assets for the frontend.

## Prerequisites

The optimization tool requires `cwebp` (Google WebP tools) to be installed on your system.

*   **macOS:** `brew install webp`
*   **Linux:** `sudo apt-get install webp`
*   **Windows:** Download `libwebp` binaries from [Google Developers](https://developers.google.com/speed/webp/docs/precompiled), extract, and add the `bin` folder to your system PATH.

## Workflow

1.  **Add Master Image:**
    *   Place the original image in `acervo/media/images`.
    *   Naming convention: `category-slug-number.jpeg` (e.g., `work-exu-nao-vem-hoje-001.jpeg`).

2.  **Reference in HTML:**
    *   In `frontend/index.html`, reference the *future* optimized file in `images/optimized/`.
    *   Example: `<img src="images/optimized/work-exu-nao-vem-hoje-001.webp">`
    *   **Note:** The builder **always** handles filename normalization. It treats underscores (`_`) in the HTML `src` as hyphens (`-`) when looking up the source file. This means `work_exu_nao_vem_hoje_001.webp` in HTML will correctly map to `work-exu-nao-vem-hoje-001.jpeg` in the source directory.

3.  **Run Builder:**
    *   Run the CLI command from `acervo/cli`:
        ```bash
        go run main.go build-assets
        ```
    *   This command:
        *   Scans `frontend/index.html` for `src="images/optimized/..."`.
        *   Finds the corresponding master in `acervo/media/images`.
        *   Executes `cwebp` to resize (max width 1920px) and convert to WebP.
        *   **Preserves Metadata:** Instructs `cwebp` to copy XMP metadata (Focus Points) from the source.
        *   Saves to `frontend/images/optimized/`.

4.  **Commit:**
    *   Commit both the source image and the generated optimized asset (for simplicity in this repo).

## Technical Details

*   **Format:** WebP (Quality 80).
*   **Dimensions:** Max width 1920px (Proportional resize).
*   **Tooling:** Uses `cwebp` via `os/exec` for robust metadata handling and compliant WebP generation without complex CGO build dependencies.
