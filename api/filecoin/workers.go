package filecoin

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"glif/config"
	"glif/constants"
	"glif/db"
	"log"
	"time"
)

func update(txHash string) {
	defer recoverFromPanic()
	ctx, cancel := context.WithTimeout(context.Background(), config.Settings.TransactionReceiptTimeout)
	defer cancel()
	ticker := time.NewTicker(config.Settings.TransactionTickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status, err := getTransactionStatus(txHash, ctx)
			if err != nil {
				log.Println("Error getting transaction status:", err)
				continue
			}

			err = db.UpdateTransaction(txHash, status)
			if err != nil {
				log.Println("Error updating transaction status:", err)
				continue
			}

			if status != constants.DEFAULT_STATUS {
				log.Println("Transaction complete:", txHash, "Status:", status)
				err = db.UpdateTransaction(txHash, status)
				if err != nil {
					log.Println("Error updating transaction status:", err)
				}
				return // Exit goroutine
			}
		case <-ctx.Done():
			err := db.UpdateTransaction(txHash, constants.TIME_OUT)
			if err != nil {
				log.Println("Error updating transaction status:", err)
			}
			log.Println("Update goroutine cancelled:", txHash)
			return
		}
	}
}

func getTransactionStatus(txHash string, ctx context.Context) (string, error) {
	client, err := ethclient.DialContext(ctx, config.Settings.RcpURL)
	if err != nil {
		return "", err
	}

	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return "", err
	}

	return constants.StatusMap[receipt.Status], nil
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}
