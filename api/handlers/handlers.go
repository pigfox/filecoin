package handlers

import (
	"encoding/json"
	"glif/db"
	"glif/filecoin"
	"glif/structs"
	"net/http"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "missing address", http.StatusBadRequest)
		return
	}

	balance, err := filecoin.Balance(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		http.Error(w, "unmarshal error", http.StatusBadRequest)
		return
	}
}

func SubmitTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var transfer structs.Transfer

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	transaction, err := filecoin.SubmitTransaction(transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	transactions, err := db.GetTransactions(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
