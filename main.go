package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/somkieatWO/qr-code-generator/apidocs" // swagger docs
	"github.com/somkieatWO/qr-code-generator/internal/handler"
	"github.com/somkieatWO/qr-code-generator/internal/usecase"
)

// @title QR Code Generator API
// @version 1.0
// @description API for generating QR codes with optional icon overlay
// @BasePath /
func main() {
	uc := usecase.NewQRGenerator(256)
	h := handler.NewQRHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/qr", h.GenerateQR)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Render (and many PaaS) inject PORT env var; fall back to 8080 for local dev
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	server := &http.Server{Addr: addr, Handler: mux}
	log.Printf("QR code generator server listening on %s", addr)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
