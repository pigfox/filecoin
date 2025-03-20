package main

import (
	"glif/config"
	"glif/handlers"
	"glif/recover"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	defer recover.RecoverPanic()
	http.HandleFunc("/wallet", handlers.BalanceHandler)
	http.HandleFunc("/transaction", handlers.SubmitTransactionHandler)
	http.HandleFunc("/transactions", handlers.GetTransactionsHandler)

	log.Println("Server is running on port", config.Settings.ServerPort)
	log.Fatal(http.ListenAndServe(config.Settings.ServerPort, nil))
}
