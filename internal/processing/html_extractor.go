package processing

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
)

type HTMLExtractor struct{}

func NewHTMLExtractor() *HTMLExtractor {
	return &HTMLExtractor{}
}

func (e *HTMLExtractor) ExtractHTML(rawHTML string) (string, error) {
	if rawHTML == "" {
		return "", fmt.Errorf("empty rawHTML string")
	}

	// Step 1: Decode from base64
	compressedData, err := base64.StdEncoding.DecodeString(rawHTML)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Step 2: Decompress using zlib (this is what pako uses)
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decompress zlib data: %w", err)
	}

	return string(decompressed), nil
}
