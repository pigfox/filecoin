package main

import (
	"fmt"
	"glif/config"
	"glif/handlers"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	defer recoverFromPanic()
	http.HandleFunc("/wallet", handlers.BalanceHandler)
	http.HandleFunc("/transaction", handlers.SubmitTransactionHandler)
	http.HandleFunc("/transactions", handlers.GetTransactionsHandler)

	log.Println("Server is running on port", config.Settings.ServerPort)
	log.Fatal(http.ListenAndServe(config.Settings.ServerPort, nil))
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}
