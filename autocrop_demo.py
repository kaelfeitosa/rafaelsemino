from PIL import Image, ImageOps, ImageFilter
import sys
import os
import glob

def autocrop(input_path, output_path, target_size=(600, 600)):
    try:
        img = Image.open(input_path)
        # Convert to grayscale and find edges to detect "detail" (potential focal points)
        detail = img.convert("L").filter(ImageFilter.FIND_EDGES)
        
        # Grid search for highest detail area
        w, h = detail.size
        grid_size = 20
        cell_w, cell_h = w // grid_size, h // grid_size
        
        max_detail = -1
        focal_point = (w // 2, h // 2)
        
        for i in range(grid_size):
            for j in range(grid_size):
                box = (i*cell_w, j*cell_h, (i+1)*cell_w, (j+1)*cell_h)
                cell = detail.crop(box)
                d = sum(cell.getdata())
                if d > max_detail:
                    max_detail = d
                    focal_point = (i*cell_w + cell_w//2, j*cell_h + cell_h//2)
        
        crop_size = min(w, h)
        left = max(0, min(focal_point[0] - crop_size // 2, w - crop_size))
        top = max(0, min(focal_point[1] - crop_size // 2, h - crop_size))
        
        cropped = img.crop((left, top, left + crop_size, top + crop_size))
        resized = cropped.resize(target_size, Image.Resampling.LANCZOS)
        
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        resized.save(output_path, quality=95)
        print(f"Processed: {os.path.basename(input_path)}")
    except Exception as e:
        print(f"Error processing {input_path}: {e}")

def batch_process(base_dir):
    # Search for all jpeg/png in subdirectories except 'cropped'
    extensions = ['*.jpeg', '*.jpg', '*.png']
    for ext in extensions:
        # Get files recursively but ignore 'cropped' folder
        files = glob.glob(os.path.join(base_dir, '**', ext), recursive=True)
        for f in files:
            if 'cropped' in f:
                continue
            
            # Map original path to cropped path
            # e.g. images/portfolio_rafael/img.jpg -> images/cropped/portfolio_rafael/img.jpg
            rel_path = os.path.relpath(f, base_dir)
            output_path = os.path.join(base_dir, 'cropped', rel_path)
            autocrop(f, output_path)

if __name__ == "__main__":
    images_root = os.path.join(os.getcwd(), 'images')
    if not os.path.exists(images_root):
        print(f"Images directory not found at {images_root}")
    else:
        print(f"Starting batch process in {images_root}...")
        batch_process(images_root)
        print("Done!")
