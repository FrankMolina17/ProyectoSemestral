package main

import (
	"log"
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

func main() {
	proformaStore := storage.NewProformaStorage()
	proformaHandler := handlers.NewProformaHandler(proformaStore)

	mux := http.NewServeMux()

	// Módulo 2 — Proformas y Cálculo
	mux.HandleFunc("POST /proformas", proformaHandler.CrearProforma)
	mux.HandleFunc("GET /proformas", proformaHandler.ListarProformas)
	mux.HandleFunc("GET /proformas/{id}", proformaHandler.ObtenerProforma)
	mux.HandleFunc("PUT /proformas/{id}", proformaHandler.ActualizarProforma)
	mux.HandleFunc("POST /proformas/{id}/calcular", proformaHandler.CalcularProforma)

	const addr = ":8080"
	log.Printf("API escuchando en %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
