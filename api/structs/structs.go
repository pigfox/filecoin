package structs

import (
	"time"
)

type WalletBalance struct {
	FIL  string `json:"fil"`
	IFIL string `json:"iFil"`
}

type Transfer struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type Transaction struct {
	TxHash    string    `json:"tx_hash"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Amount    string    `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
