# HTML Compression Troubleshooting Guide

## Key Technical Insights

### JavaScript pako vs Go compression libraries

**Critical insight: JavaScript's `pako.deflate()` produces zlib format, not raw deflate.**

When JavaScript code uses:
```javascript
const compressed = pako.deflate(data, { level: 9 })
```

The output is zlib format (deflate + headers), which means in Go you must use:
```go
// Correct ✅
reader, err := zlib.NewReader(bytes.NewReader(compressedData))

// Wrong ❌ - will fail with "corrupt input before offset X"
reader := flate.NewReader(bytes.NewReader(compressedData))
```

### JSON HTML Escaping

Go's default JSON marshaling escapes HTML characters:
- `<` becomes `\u003c` 
- `>` becomes `\u003e`

**Before fix:**
```json
"RawHTMLExtracted": "\u003cdiv class=\"e-13udsys\"\u003e"
```

**After fix:**
```json
"RawHTMLExtracted": "<div class=\"e-13udsys\">"
```

**Solution:**
```go
encoder := json.NewEncoder(&buf)
encoder.SetEscapeHTML(false)  // This is the key line
```

## Debugging Compression Issues

### CLI Testing Commands

Test base64 + zlib decompression:
```bash
# Extract rawHtml from DynamoDB JSON and test decompression
cat data.json | jq -r '.Items[0].rawHtml.S' | python3 -c "
import base64, zlib, sys
data = base64.b64decode(sys.stdin.read().strip())
try:
    result = zlib.decompress(data).decode('utf-8')
    print('SUCCESS! First 200 chars:', result[:200])
except Exception as e:
    print('FAILED:', e)
"
```

Check compression format:
```bash
echo 'eNrNVltv4ygU...' | base64 -d | file -
```

Inspect raw bytes:
```bash
echo 'eNrNVltv4ygU...' | base64 -d | xxd | head -5
```

### Common Error Messages

| Error | Cause | Solution |
|-------|--------|----------|
| `flate: corrupt input before offset X` | Using `flate.NewReader()` on zlib data | Use `zlib.NewReader()` instead |
| `failed to decode base64` | Invalid base64 string | Check base64 string is complete |
| `\u003c` in JSON output | HTML escaping enabled | Use `SetEscapeHTML(false)` |

## Format Detection

The first few bytes of the decompressed data can help identify the format:
- Zlib: Usually starts with `0x78` (hex)
- Raw deflate: No standard header
- Check with: `echo 'base64...' | base64 -d | xxd | head -1`
```
