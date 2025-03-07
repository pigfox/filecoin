package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var PS Postgres

type Postgres struct {
	Conn *sql.DB
}

func Connect() (Postgres, error) { // Return Postgres and error
	dbConnString := "postgres://" + os.Getenv("DATABASE_USER") + ":" + os.Getenv("DATABASE_PASSWORD")
	dbConnString += "@" + os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT") + "/" + os.Getenv("DATABASE_NAME")
	dbConnString += "?sslmode=disable&connect_timeout=" + os.Getenv("DATABASE_TIMEOUT")
	conn, err := sql.Open("postgres", dbConnString)
	if err != nil {
		return Postgres{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close() //close connection on error.
		return Postgres{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return Postgres{Conn: conn}, nil
}

func Get() Postgres {
	if PS.Conn == nil {
		newPS, err := Connect()
		if err != nil {
			log.Println("Database connection error:", err)
			return Postgres{} // Return empty Postgres struct
		}
		PS = newPS
	}
	return PS
}
