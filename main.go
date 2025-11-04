package main

import (
	"log"
	"net/http"

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

	server := &http.Server{Addr: ":8080", Handler: mux}
	log.Println("QR code generator server listening on :8080")
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
