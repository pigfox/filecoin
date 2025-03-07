package filecoin

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"glif/config"
	"glif/constants"
	"glif/db"
	"glif/structs"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

func Balance(address string) (structs.WalletBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Settings.RequestTimeout)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.Settings.RcpURL)
	if err != nil {
		return structs.WalletBalance{}, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	account := common.HexToAddress(address)

	balanceWei, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return structs.WalletBalance{}, fmt.Errorf("failed to retrieve balance: %w", err)
	}

	balanceEther := new(big.Float).Quo(new(big.Float).SetInt(balanceWei), big.NewFloat(math.Pow10(18)))

	return structs.WalletBalance{
		FIL:  balanceEther.Text('f', 18), // Format to 18 decimal places
		IFIL: strconv.Itoa(0),
	}, nil
}

func SubmitTransaction(transfer structs.Transfer) (structs.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Settings.TransactionTimeout)
	defer cancel()

	client, err := ethclient.Dial(config.Settings.RcpURL)
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(transfer.Amount, 10) // Convert from string to big.Int
	if !ok {
		return structs.Transaction{}, fmt.Errorf("invalid amount format")
	}

	fromAddress := common.HexToAddress(transfer.Sender)
	toAddress := common.HexToAddress(transfer.Receiver)

	privateKeyHex := config.Settings.PrivateKey
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to parse private key: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to get gas price: %v", err)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: amount,
	})
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to estimate gas limit: %v", err)
	}

	// Increase the gas limit slightly to ensure transaction success
	gasLimit = gasLimit * 12 / 10 // 20% buffer

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    amount,
		Data:     nil, // No data for a simple ether transfer
	})

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to get network ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		return structs.Transaction{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	txHash := signedTx.Hash().Hex()

	transaction := structs.Transaction{
		TxHash:    txHash,
		Sender:    transfer.Sender,
		Receiver:  transfer.Receiver,
		Amount:    transfer.Amount,
		Status:    constants.DEFAULT_STATUS,
		Timestamp: time.Now(),
	}

	err = db.InsertTransaction(transaction)
	if err != nil {
		log.Println("DB insert failed", txHash, err.Error())
	}

	go update(txHash)

	return transaction, nil
}
