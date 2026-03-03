package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strconv"

	pngstructure "github.com/dsoprea/go-png-image-structure"
)

const (
	xmpNamespace = "http://ns.adobe.com/xap/1.0/\x00"
)

// SetFocus embeds XMP focus point metadata into an image file using pure Go.
func SetFocus(imagePath string, x, y float64) error {
	xmpBytes := generateXMP(x, y)

	ext := ""
	for i := len(imagePath) - 1; i >= 0 && imagePath[i] != '.'; i-- {
		ext = string(imagePath[i]) + ext
	}

	switch ext {
	case "jpg", "jpeg":
		return setFocusJPEG(imagePath, xmpBytes)
	case "png":
		return setFocusPNG(imagePath, xmpBytes)
	case "webp":
		return InjectXMPWebP(imagePath, xmpBytes)
	default:
		return fmt.Errorf("formato de imagem não suportado para metadados: %s", ext)
	}
}

// HasFocus checks if an image has a defined focal point in its XMP metadata.
func HasFocus(imagePath string) (bool, error) {
	xmp, err := ExtractXMP(imagePath)
	if err != nil {
		return false, err
	}
	if xmp == nil {
		return false, nil
	}
	// Focus point is defined by mwg-rs:RegionType>Focus
	return bytes.Contains(xmp, []byte("<mwg-rs:RegionType>Focus</mwg-rs:RegionType>")), nil
}

// ExtractXMP attempts to extract XMP metadata from JPEG or PNG files.
func ExtractXMP(imagePath string) ([]byte, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	ext := ""
	for i := len(imagePath) - 1; i >= 0 && imagePath[i] != '.'; i-- {
		ext = string(imagePath[i]) + ext
	}

	switch ext {
	case "jpg", "jpeg":
		return extractXMPJPEG(data)
	case "png":
		return extractXMPPNG(data)
	default:
		return nil, nil // Silently ignore unsupported formats for extraction
	}
}

func extractXMPJPEG(data []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
		return nil, fmt.Errorf("não é um JPEG válido")
	}
	offset := 2
	for offset < len(data) {
		for offset < len(data) && data[offset] != 0xFF {
			offset++
		}
		if offset >= len(data)-1 {
			break
		}
		marker := data[offset+1]
		if marker == 0xDA || marker == 0xD9 {
			break
		}
		if (marker >= 0xD0 && marker <= 0xD7) || marker == 0xD8 || marker == 0x01 {
			offset += 2
			continue
		}
		if offset+3 >= len(data) {
			break
		}
		length := int(binary.BigEndian.Uint16(data[offset+2 : offset+4]))
		if marker == 0xE1 && length >= 2+len(xmpNamespace)-1 {
			payload := data[offset+4 : offset+2+length]
			if bytes.HasPrefix(payload, []byte("http://ns.adobe.com/xap/1.0/")) {
				// Offset depends on whether there's a null terminator
				nsLen := 29
				if len(payload) > 29 && payload[29] == 0 {
					nsLen = 30
				}
				return payload[nsLen:], nil
			}
		}
		offset += 2 + length
	}
	return nil, nil
}

func extractXMPPNG(data []byte) ([]byte, error) {
	pmp := pngstructure.NewPngMediaParser()
	intfc, err := pmp.ParseBytes(data)
	if err != nil {
		return nil, err
	}
	cs := intfc.(*pngstructure.ChunkSlice)
	for _, chunk := range cs.Chunks() {
		if chunk.Type == "iTXt" && bytes.HasPrefix(chunk.Data, []byte("XML:com.adobe.xmp")) {
			// PNG iTXt: keyword (null) compressionflag (null) compressionmethod (null) lang (null) translatedkey (null) text
			parts := bytes.SplitN(chunk.Data, []byte{0}, 6)
			if len(parts) > 5 {
				return parts[5], nil
			}
		}
	}
	return nil, nil
}

// InjectXMPWebP embeds XMP metadata into a WebP file.
func InjectXMPWebP(imagePath string, xmp []byte) error {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return fmt.Errorf("não é um arquivo WebP válido")
	}

	// We'll rebuild the WebP to ensure it has VP8X if it needs it.
	// But as a simpler approach for now, we'll just check if it has VP8X.
	hasVP8X := string(data[12:16]) == "VP8X"

	var buf bytes.Buffer
	buf.Write(data[0:12]) // RIFF header

	if !hasVP8X {
		// Insert VP8X chunk
		// We'll need dimensions. For now, since cwebp generated it,
		// we might need to parse dimensions from VP8/VP8L if we wanted thoroughness.
		// BUT if we just want to ADD XMP, we MUST have VP8X.
		// Let's implement a minimal VP8X insertion.
		// To get dimensions, we can parse them from VP8/VP8L.
		width, height := parseDimensions(data)

		buf.Write([]byte("VP8X"))
		binary.Write(&buf, binary.LittleEndian, uint32(10))
		flags := byte(0x10) // Set XMP bit (bit 4, 0x10)
		buf.WriteByte(flags)
		buf.Write([]byte{0, 0, 0}) // Reserved

		// Dimensions are stored as 24-bit (3 bytes), value-1
		w := uint32(width - 1)
		h := uint32(height - 1)
		buf.Write([]byte{byte(w), byte(w >> 8), byte(w >> 16)})
		buf.Write([]byte{byte(h), byte(h >> 8), byte(h >> 16)})

		buf.Write(data[12:]) // Original chunks
	} else {
		// Update VP8X if it exists
		flags := data[20]
		data[20] = flags | 0x10 // Set XMP bit
		buf.Write(data[12:])
	}

	// Append XMP chunk
	buf.Write([]byte("XMP "))
	binary.Write(&buf, binary.LittleEndian, uint32(len(xmp)))
	buf.Write(xmp)
	if len(xmp)%2 != 0 {
		buf.WriteByte(0) // Padding
	}

	finalData := buf.Bytes()
	// Update total RIFF size
	binary.LittleEndian.PutUint32(finalData[4:8], uint32(len(finalData)-8))

	return os.WriteFile(imagePath, finalData, 0644)
}

func parseDimensions(data []byte) (int, int) {
	for i := 12; i < len(data)-8; {
		chunkType := string(data[i : i+4])
		size := int(binary.LittleEndian.Uint32(data[i+4 : i+8]))
		if chunkType == "VP8 " {
			// VP8 bitstream header
			if i+18 <= len(data) {
				w := int(data[i+14]) | (int(data[i+15]&0x3F) << 8)
				h := int(data[i+16]) | (int(data[i+17]&0x3F) << 8)
				return w, h
			}
		} else if chunkType == "VP8L" {
			// VP8L bitstream header
			if i+12 <= len(data) {
				// 14 bits for w/h
				bits := uint32(data[i+9]) | (uint32(data[i+10]) << 8) | (uint32(data[i+11]) << 16) | (uint32(data[i+12]) << 24)
				w := int(bits&0x3FFF) + 1
				h := int((bits>>14)&0x3FFF) + 1
				return w, h
			}
		}
		i += 8 + size + (size % 2)
	}
	return 0, 0
}

func generateXMP(x, y float64) []byte {
	// Standard XMP MWG-RS (Regions) schema for a single point of focus
	xmp := fmt.Sprintf(`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Adobe XMP Core 5.1.0-jc003">
 <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about=""
    xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/"
    xmlns:stArea="http://www.metadataworkinggroup.com/schemas/regions/struct/Area#">
   <mwg-rs:Regions rdf:parseType="Resource">
    <mwg-rs:RegionList>
     <rdf:Bag>
      <rdf:li rdf:parseType="Resource">
       <mwg-rs:RegionType>Focus</mwg-rs:RegionType>
       <mwg-rs:Area rdf:parseType="Resource">
        <stArea:x>%.4f</stArea:x>
        <stArea:y>%.4f</stArea:y>
        <stArea:w>0.0000</stArea:w>
        <stArea:h>0.0000</stArea:h>
        <stArea:unit>normalized</stArea:unit>
       </mwg-rs:Area>
      </rdf:li>
     </rdf:Bag>
    </mwg-rs:RegionList>
   </mwg-rs:Regions>
  </rdf:Description>
 </rdf:RDF>
</x:xmpmeta>
<?xpacket end="w"?>`, x, y)
	return []byte(xmp)
}

func setFocusJPEG(imagePath string, xmpBytes []byte) error {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
		return fmt.Errorf("não é um arquivo JPEG válido")
	}

	var buf bytes.Buffer
	// Write SOI
	buf.Write([]byte{0xFF, 0xD8})

	// Create and write new APP1 XMP Segment
	xmpData := append([]byte(xmpNamespace), xmpBytes...)
	segLen := uint16(len(xmpData) + 2)
	buf.Write([]byte{0xFF, 0xE1})
	binary.Write(&buf, binary.BigEndian, segLen)
	buf.Write(xmpData)

	offset := 2
	for offset < len(data) {
		// Find next marker (skip padding 0xFFs)
		for offset < len(data) && data[offset] != 0xFF {
			offset++
		}
		if offset >= len(data)-1 {
			break
		}

		marker := data[offset+1]

		// If padding or RST or TEM logic
		if marker == 0xFF {
			offset++
			continue
		}

		// SOS (Start of Scan) or EOI (End of Image)
		if marker == 0xDA || marker == 0xD9 {
			// From here on, just copy everything verbatim!
			// This preserves all scan data, headers, and trailing data perfectly.
			buf.Write(data[offset:])
			break
		}

		// Markers without length (RST0-7, SOI, TEM)
		if (marker >= 0xD0 && marker <= 0xD7) || marker == 0xD8 || marker == 0x01 {
			buf.Write(data[offset : offset+2])
			offset += 2
			continue
		}

		// Markers with length
		if offset+3 >= len(data) {
			break
		}
		length := binary.BigEndian.Uint16(data[offset+2 : offset+4])

		// Ensure length is valid to prevent panics on corrupted files
		if int(offset)+2+int(length) > len(data) {
			break
		}

		// Check if it's an existing XMP APP1 segment
		isXMP := false
		if marker == 0xE1 && length >= 2+uint16(len(xmpNamespace)) {
			payload := data[offset+4 : offset+2+int(length)]
			if bytes.HasPrefix(payload, []byte(xmpNamespace)) {
				isXMP = true
			}
		}

		// If not our XMP, write it to the new file
		if !isXMP {
			buf.Write(data[offset : offset+2+int(length)])
		}

		offset += 2 + int(length)
	}

	// Write safely to a temp file, then rename/overwrite
	tmpPath := imagePath + ".tmp"
	if err := os.WriteFile(tmpPath, buf.Bytes(), 0644); err != nil {
		return err
	}
	// Windows file replace requires target to be removed first or use os.Rename with overwrite logic depending on OS
	// Since os.Rename on Windows fails if target exists, we remove it.
	os.Remove(imagePath)
	return os.Rename(tmpPath, imagePath)
}

func setFocusPNG(imagePath string, xmpBytes []byte) error {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	pmp := pngstructure.NewPngMediaParser()
	intfc, err := pmp.ParseBytes(data)
	if err != nil {
		return err
	}

	cs := intfc.(*pngstructure.ChunkSlice)

	// Filter chunks
	chunks := cs.Chunks()
	newChunks := make([]*pngstructure.Chunk, 0)
	for _, chunk := range chunks {
		if chunk.Type == "iTXt" && bytes.Contains(chunk.Data, []byte("XML:com.adobe.xmp")) {
			continue
		}
		newChunks = append(newChunks, chunk)
	}

	// Create new iTXt chunk
	var buffer bytes.Buffer
	buffer.WriteString("XML:com.adobe.xmp")
	buffer.WriteByte(0) // Null separator 1
	buffer.WriteByte(0) // Compression flag
	buffer.WriteByte(0) // Compression method
	buffer.WriteByte(0) // Null separator 2 (after empty lang)
	buffer.WriteByte(0) // Null separator 3 (after empty translated)
	buffer.Write(xmpBytes)

	xmpChunk := &pngstructure.Chunk{
		Type: "iTXt",
		Data: buffer.Bytes(),
	}

	// Insert before IDAT (mandatory for some parsers)
	insertPos := 1
	for i, chunk := range newChunks {
		if chunk.Type == "IDAT" {
			insertPos = i
			break
		}
	}

	finalChunks := make([]*pngstructure.Chunk, 0)
	finalChunks = append(finalChunks, newChunks[:insertPos]...)
	finalChunks = append(finalChunks, xmpChunk)
	finalChunks = append(finalChunks, newChunks[insertPos:]...)

	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// PNG Signature
	if _, err := f.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}); err != nil {
		return err
	}

	for _, chunk := range finalChunks {
		if err := writePngChunk(f, chunk); err != nil {
			return err
		}
	}

	return nil
}

func writePngChunk(w io.Writer, chunk *pngstructure.Chunk) error {
	length := uint32(len(chunk.Data))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	if _, err := w.Write([]byte(chunk.Type)); err != nil {
		return err
	}
	if _, err := w.Write(chunk.Data); err != nil {
		return err
	}
	crc := crc32.NewIEEE()
	crc.Write([]byte(chunk.Type))
	crc.Write(chunk.Data)
	return binary.Write(w, binary.BigEndian, crc.Sum32())
}

func GetFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
