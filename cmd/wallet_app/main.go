package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"wallet/internal/handler"
	"wallet/internal/repository"
	"wallet/internal/service"
)

func main() {
	_ = godotenv.Load("config.env")
	dbURL := os.Getenv("DATABASE_URL")
	db, err := repository.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	svc := service.NewWalletService(db)
	h := handler.NewWalletHandler(svc)

	r.Post("/api/v1/wallet", h.HandleTransaction)
	r.Get("/api/v1/wallets/{id}", h.GetBalance)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
