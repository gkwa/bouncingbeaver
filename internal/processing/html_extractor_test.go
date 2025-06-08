package processing

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"testing"
)

func TestHTMLExtractor_ExtractHTML(t *testing.T) {
	extractor := NewHTMLExtractor()

	// Test with empty string - should return error as per implementation
	_, err := extractor.ExtractHTML("")
	if err == nil {
		t.Error("Expected error for empty string")
	}

	// Test with valid zlib + base64 encoded data (matching pako format)
	originalHTML := "<html><body>Hello World</body></html>"

	// Compress with zlib (this matches what pako does)
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	zw.Write([]byte(originalHTML))
	zw.Close()

	// Encode with base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Test extraction
	result, err := extractor.ExtractHTML(encoded)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != originalHTML {
		t.Errorf("Expected %s, got %s", originalHTML, result)
	}
}

func TestHTMLExtractor_ExtractHTML_ActualData(t *testing.T) {
	extractor := NewHTMLExtractor()

	// Test with the actual base64 string from your working data (second item)
	rawHTML := "eNrNVttu4zYQ/RWCzUMWiO5yJDmxgWbRdgu06D5tHwtKpCTalKiQtCz76zuU7PVt0+YxNiiRMyLnzIWHfKa8R4UgWi8wc4JoQ/VO4+UziJfPdXRSGZ7MVuGkQURx4giSM7HAX5Wkm8JgpKRgC1wpuenwaV6pW5lymEcOH+QbY2SLUa1YucC1MZ2eex5lgvdM7dyuKBqi1sxot5CNp41UzAOhA6Nm03Kzcw56r5ssay9KAvg7OZgEVKpizlZKamfkzpZTAC+lqZ2St8euVfGWGC5bJ3As+tElyjXJBaMLXBKh2ZkbjPXbJD+4/12ar9p0W19Lg4xT2dxIG59uHsm1mBVBlVspbyqkVaGZOUVlu926vNWGFESZMR68IeCeZgqC5QVZMkDzSi4MU3oOb3H/6/h7MGrDPs1LqRpi7ldd9cmjoWjVrKlJku+KlVsIuaGlkq1xW2aOwXRGA3ZF5o2h/If5aVams8JJy5Q5cU6Jk4Xs0YnifDaLQEeT0u3a6uGdoMPscYD2MUCjwJ0ND+id0KMsHqB9EOjh+4HPsmCA9kGAR+8HnqTpAO2DAI8HjCgxxDFMGw4swQ1rHEBOp0WBRQTsXnxOBSyiuoTt7Y2Mevs8Pk5z/GuKaB6TfrO7lq59P33V11KdBFVppboj7VGsC8VY6yhGKFOObMUOI212IxtLBbI58rvhCT7n3RwpVph7GD+gi8enJ1QzXtVmjgL7MRCrqQ99YOSKt3PkjCMJmSuF3M5RzSll7RPqCKW8rQ5mOqm5Zd45IrmWYmPYE15+3ijFWoM6xQs2R3epm2XPnvXi0heIKe+iIK6OPo7EPRlaYFsTZ+Hvkngm8PLuYqH/mhC8rnuTMrxM3z3lYOOE9vi6zWwgm0zWtyeDUvHra4J/POe7ckTy0dL2l+KwBBHoq83b0ffuhL9vs1LM8HUSyb7LJHh1l53l2eve2CGneTTqs73Ct19Ml4saShwQH45zwXp7Q4nPE9y1ZWr6mxTEyVqEBV6+gAB9theHPyw3jN3jReFvuFS8jWy/4uvksK7hxqIJECNFfWadtLEc4LCfFG9ygkcgXJBOMGrL6di7tdl3UBv7G2fyfRMlFst017q4rv1MKQqQJS70f76eh82fhf3syEEXHAg18jus9jKa+mWARNoMgJnby1LR0LWthL6aynCBQ8jNVJ5Tv+ds+yKHBfaRj8IYWZnl/AX+aWJ9jIZGtHq6Ih3Oj23kSlV5oe/7HqwN3Mb3bFqvkEIqIMCdBoi/KbLTBRHM9/EPtjSULTE1Apf+DHw3TVEQuUH4LfTr0A3j3nkE2ZfQ7x07rMfhN3iPylEX97Zv0+HZlWzq+uq68JWRVb7CSwjQJVMcnlPK3j4vvDq6EPwLr8Dpqw=="

	result, err := extractor.ExtractHTML(rawHTML)
	if err != nil {
		t.Fatalf("Failed to extract HTML from actual data: %v", err)
	}

	// The result should contain HTML content
	if len(result) == 0 {
		t.Error("Expected non-empty HTML result")
	}

	// Should contain some expected HTML elements
	if !containsAny(result, []string{"<div", "<a", "<span"}) {
		t.Errorf("Expected HTML elements in result, got first 200 chars: %s", result[:min(len(result), 200)])
	}

	t.Logf("Successfully extracted HTML (%d characters)", len(result))
	t.Logf("First 200 chars: %s", result[:min(len(result), 200)])
}

func TestHTMLExtractor_ExtractHTML_InvalidBase64(t *testing.T) {
	extractor := NewHTMLExtractor()

	_, err := extractor.ExtractHTML("invalid-base64!")
	if err == nil {
		t.Error("Expected error for invalid base64")
	}
}

func TestHTMLExtractor_ExtractHTML_InvalidZlib(t *testing.T) {
	extractor := NewHTMLExtractor()

	// Valid base64 but invalid zlib data
	invalidZlib := base64.StdEncoding.EncodeToString([]byte("not zlib data"))
	_, err := extractor.ExtractHTML(invalidZlib)
	if err == nil {
		t.Error("Expected error for invalid zlib data")
	}
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
