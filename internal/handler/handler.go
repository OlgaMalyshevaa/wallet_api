package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"wallet/internal/model"
	"wallet/internal/service"

	"github.com/go-chi/chi/v5"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{service: svc}
}

func (h *WalletHandler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	go func() {
		if err := h.service.PerformTransaction(r.Context(), req); err != nil {
			log.Println("Transaction error:", err)
			return
		}
		log.Println("Transaction successful for wallet:", req.WalletID)
	}()

	w.WriteHeader(http.StatusAccepted)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	balance, err := h.service.GetBalance(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"balance": balance})
}
