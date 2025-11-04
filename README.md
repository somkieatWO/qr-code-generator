# QR Code Generator Service

Simple Go HTTP service for generating QR codes with optional center icon and custom size.

## Endpoints

GET /qr

Query parameters:
- text (string, required): Data to encode (URL or any text)
- iconUrl (string, optional): Public URL to an image (PNG/JPG/GIF) to draw centered (~20% of QR size)
- size (int, optional): Pixel size of the square QR image (64-2048). Defaults to 256.

POST /qr

Multipart form fields:
- text (string, required): Data to encode (URL or any text)
- icon (file, optional): Image file (PNG/JPG/GIF) drawn centered at ~20% of QR size
- size (int, optional): Pixel size of the square QR image (64-2048). Defaults to 256.

Response (both):
- 200 image/png (QR code image)
- 4xx on validation errors

## Run

```bash
go mod tidy
go run ./main.go
```

Server listens on :8080.

## Examples

GET (simple):
```bash
curl "http://localhost:8080/qr?text=https://example.com" --output qr.png
```

GET with size and iconUrl:
```bash
curl "http://localhost:8080/qr?text=Hello%20icon&size=512&iconUrl=https://example.com/icon.png" --output qr-with-icon.png
```

POST (supports file upload for icon):
```bash
curl -X POST -F 'text=https://example.com' http://localhost:8080/qr --output qr.png
```

POST with icon file and size:
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

## Swagger UI on GitHub Pages
GitHub Pages hosts only static content and blocks non-GET requests. That means calling the API path on `somkieatwo.github.io` (e.g., `POST https://somkieatwo.github.io/qr-code-generator/qr`) will return `405 Not Allowed`.

You can still use the hosted Swagger UI as a client by pointing it to your real API server using query params:

- Local dev server:
  - Run the server locally: `go run ./main.go`
  - Open: https://somkieatwo.github.io/qr-code-generator/docs/index.html?host=localhost:8080&schemes=http
- Public HTTPS API (example):
  - https://somkieatwo.github.io/qr-code-generator/docs/index.html?host=api.example.com&schemes=https

If your API is rooted under a subpath, also pass `basePath`, e.g. `&basePath=/v1`.

