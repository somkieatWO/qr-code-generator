package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/somkieatWO/qr-code-generator/internal/usecase"
)

// QRHandler handles QR code generation requests
type QRHandler struct {
	qr *usecase.QRGenerator // default generator
}

// NewQRHandler creates a new QRHandler with the given QRGenerator
func NewQRHandler(qr *usecase.QRGenerator) *QRHandler { return &QRHandler{qr: qr} }

// GenerateQR godoc
// @Summary Generate QR code or Barcode
// @Description Generates a QR code or Barcode PNG for the provided text with optional icon (QR only) and custom size
// @Tags qr
// @Accept multipart/form-data
// @Produce image/png
// @Param text formData string true "Text or URL to encode"
// @Param type formData string false "Type of code to generate: 'qr' (default) or 'barcode'"
// @Param icon formData file false "Center icon image (PNG/JPG/GIF) - QR code only"
// @Param size formData int false "Pixel size (64-2048)"
// @Success 200 {file} png "QR code or Barcode image"
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

	codeType := r.FormValue("type")
	if codeType == "" {
		codeType = "qr"
	}
	codeType = strings.ToLower(codeType)

	if codeType != "qr" && codeType != "barcode" {
		http.Error(w, "invalid type, must be 'qr' or 'barcode'", http.StatusBadRequest)
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

	var pngBytes []byte
	var err error

	if codeType == "barcode" {
		pngBytes, err = gen.GenerateBarcode(text)
	} else {
		var iconBytes []byte
		file, _, errFile := r.FormFile("icon")
		if errFile == nil && file != nil {
			defer file.Close()
			iconBytes, errFile = io.ReadAll(file)
			if errFile != nil {
				log.Printf("read icon error: %v", errFile)
				http.Error(w, "invalid icon file", http.StatusBadRequest)
				return
			}
		}
		pngBytes, err = gen.Generate(text, iconBytes)
	}

	if err != nil {
		log.Printf("generate error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if codeType == "barcode" {
		w.Header().Set("Content-Disposition", "inline; filename=barcode.png")
	} else {
		w.Header().Set("Content-Disposition", "inline; filename=qr.png")
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngBytes)
}
