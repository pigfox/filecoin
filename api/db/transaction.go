package db

import (
	"context"
	"database/sql"
	"glif/structs"
	"log"
)

func InsertTransaction(transaction structs.Transaction) error {
	stmt, err := Get().Conn.Prepare("INSERT INTO transactions (tx_hash, sender, receiver, amount, status, timestamp) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(transaction.TxHash, transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Status, transaction.Timestamp)

	return err
}

func GetTransactions(address string) ([]structs.Transaction, error) {
	query := "SELECT tx_hash, sender, receiver, amount, status, timestamp FROM transactions"
	var rows *sql.Rows
	var err error

	if address != "" {
		query += " WHERE sender=$1 OR receiver=$1"
		rows, err = Get().Conn.QueryContext(context.Background(), query, address)
	} else {
		rows, err = Get().Conn.QueryContext(context.Background(), query)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var transactions []structs.Transaction
	for rows.Next() {
		var tx structs.Transaction
		if err := rows.Scan(&tx.TxHash, &tx.Sender, &tx.Receiver, &tx.Amount, &tx.Status, &tx.Timestamp); err != nil {
			return transactions, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, err
}

func UpdateTransaction(txHash, status string) error {
	stmt, err := Get().Conn.Prepare("UPDATE transactions SET status = $1 WHERE tx_hash = $2")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			
		}
	}(stmt)

	result, err := stmt.Exec(status, txHash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		log.Printf("Warning: UpdateTransaction affected %d rows (expected 1)\n", rowsAffected)
	}

	return nil
}
