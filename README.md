# QR Code Generator Service

Simple Go HTTP service for generating QR codes with optional center icon and custom size.

## Endpoint

POST /qr

Multipart form fields:
- text (string, required): Data to encode (URL or any text)
- icon (file, optional): Image file (PNG/JPG/GIF) drawn centered at ~20% of QR size
- size (int, optional): Pixel size of the square QR image (64-2048). Defaults to 256.

Response:
- 200 image/png (QR code image)
- 4xx on validation errors

## Run

```bash
go mod tidy
go run ./main.go
```

Server listens on :8080.

## Examples

Generate a QR for a URL:
```bash
curl -X POST -F 'text=https://example.com' http://localhost:8080/qr --output qr.png
```

Generate with icon and custom size:
```bash
curl -X POST \
  -F 'text=Hello with icon' \
  -F 'size=512' \
  -F 'icon=@icon.png' \
  http://localhost:8080/qr --output qr-with-icon.png
```

## Tests

```bash
go test ./...
```

## Notes
- Icon is scaled preserving aspect ratio using nearest-neighbor.
- Error level is set to `qrcode.Medium` currently; adjust in `usecase/qr.go` if needed.
- For production, consider caching generated QR codes and adding request timeouts.

