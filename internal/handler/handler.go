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
// @Success 200 {file} png "QR code image"
// @Failure 400 {string} string "Error message"
// @Router /qr [post]
func (h *QRHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	// CORS headers (adjust origin in production)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions { // preflight
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed; use POST"))
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil { // 20MB for safety
		log.Printf("parse form error: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	// size handling
	reqSize := h.qr.Size()
	if sizeStr := r.FormValue("size"); sizeStr != "" {
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

	var iconBytes []byte
	file, _, err := r.FormFile("icon")
	if err == nil && file != nil {
		defer file.Close()
		iconBytes, err = io.ReadAll(file)
		if err != nil {
			log.Printf("read icon error: %v", err)
			http.Error(w, "invalid icon file", http.StatusBadRequest)
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
}
