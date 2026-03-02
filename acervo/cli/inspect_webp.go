package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: inspect_webp <file>")
		return
	}
	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer f.Close()

	header := make([]byte, 12)
	if _, err := io.ReadFull(f, header); err != nil {
		fmt.Printf("Error reading header: %v\n", err)
		return
	}

	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WEBP" {
		fmt.Println("Not a valid RIFF/WEBP file")
		return
	}

	totalSize := binary.LittleEndian.Uint32(header[4:8])
	fmt.Printf("File: %s\n", path)
	fmt.Printf("RIFF Size: %d\n", totalSize)

	offset := int64(12)
	for {
		chunkHeader := make([]byte, 8)
		n, err := io.ReadFull(f, chunkHeader)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading chunk header at offset %d: %v\n", offset, err)
			break
		}

		typeStr := string(chunkHeader[0:4])
		size := binary.LittleEndian.Uint32(chunkHeader[4:8])
		fmt.Printf("Chunk: [%s] Size: %d at Offset: %d\n", typeStr, size, offset)

		if typeStr == "VP8X" {
			payload := make([]byte, size)
			io.ReadFull(f, payload)
			fmt.Printf("  VP8X Flags: %08b\n", payload[0])
		} else {
			f.Seek(int64(size), io.SeekCurrent)
		}

		// Chunks are padded to even size
		if size%2 != 0 {
			f.Seek(1, io.SeekCurrent)
			offset += int64(size) + 9
		} else {
			offset += int64(size) + 8
		}

		if offset >= int64(totalSize)+8 {
			break
		}
	}
}
