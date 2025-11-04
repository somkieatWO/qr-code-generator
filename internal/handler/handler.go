package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/somkieatWO/qr-code-generator/internal/usecase"
)

// QRHandler handles QR code generation requests
type QRHandler struct {
	qr *usecase.QRGenerator // default generator
}

// NewQRHandler creates a new QRHandler with the given QRGenerator
func NewQRHandler(qr *usecase.QRGenerator) *QRHandler { return &QRHandler{qr: qr} }

// GenerateQR godoc
// @Summary Generate QR code
// @Description Generates a QR code PNG for the provided text with optional icon and custom size
// @Tags qr
// @Accept multipart/form-data
// @Produce image/png
// @Param text formData string true "Text or URL to encode"
// @Param icon formData file false "Center icon image (PNG/JPG/GIF)"
// @Param size formData int false "QR pixel size (64-2048)"
// @Param text query string true "Text or URL to encode"
// @Param iconUrl query string false "Center icon image URL (PNG/JPG/GIF)"
// @Param size query int false "QR pixel size (64-2048)"
// @Success 200 {file} png "QR code image"
// @Failure 400 {string} string "Error message"
// @Router /qr [get]
func (h *QRHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	// CORS headers (adjust origin in production)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions { // preflight
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Read from query string
	q := r.URL.Query()
	text := q.Get("text")
	if text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	// size handling
	reqSize := h.qr.Size()
	if sizeStr := q.Get("size"); sizeStr != "" {
		parsed, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "size must be an integer", http.StatusBadRequest)
			return
		}
		if parsed < 64 || parsed > 2048 {
			http.Error(w, "size out of range (64-2048)", http.StatusBadRequest)
			return
		}
		reqSize = parsed
	}
	// create a fresh generator if size differs
	gen := h.qr
	if reqSize != h.qr.Size() {
		gen = usecase.NewQRGenerator(reqSize)
	}

	// Optional icon via URL
	var iconBytes []byte
	if iconURL := q.Get("iconUrl"); iconURL != "" {
		// Simple fetch with size guard (max 5MB)
		resp, err := http.Get(iconURL)
		if err != nil {
			log.Printf("fetch icon url error: %v", err)
			http.Error(w, "failed to download iconUrl", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			http.Error(w, "iconUrl fetch returned non-2xx", http.StatusBadRequest)
			return
		}
		// limit to 5MB
		const maxIcon = 5 << 20
		limited := io.LimitReader(resp.Body, maxIcon+1)
		iconBytes, err = io.ReadAll(limited)
		if err != nil {
			log.Printf("read icon url error: %v", err)
			http.Error(w, "invalid iconUrl", http.StatusBadRequest)
			return
		}
		if len(iconBytes) > maxIcon {
			http.Error(w, "iconUrl too large (max 5MB)", http.StatusBadRequest)
			return
		}
	}

	pngBytes, err := gen.Generate(text, iconBytes)
	if err != nil {
		log.Printf("generate qr error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "inline; filename=qr.png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngBytes)
	return

}
