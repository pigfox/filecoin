package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var Settings Values

type Values struct {
	ServerPort                string
	RequestTimeout            time.Duration
	TransactionTimeout        time.Duration
	TransactionReceiptTimeout time.Duration
	TransactionTickerInterval time.Duration
	RcpURL                    string
	PrivateKey                string
}

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Settings.RequestTimeout = 5 * time.Second
	Settings.TransactionTimeout = 30 * time.Second
	Settings.TransactionReceiptTimeout = 3600 * time.Second
	Settings.TransactionTickerInterval = 3 * time.Second
	Settings.ServerPort = ":8080"
	Settings.RcpURL = os.Getenv("RPC_URL")
	Settings.PrivateKey = os.Getenv("PRIVATE_KEY")
}
