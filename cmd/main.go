package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	config := mysql.Config{
		User:                 "root",
		Passwd:               "admin",
		Addr:                 "localhost",
		DBName:               "TestDb",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error("error trying to close db: ", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	slog.Info("connected")

	iterations := 100000

	version := 1

	ids := make([]string, iterations)

	for i := 0; i < iterations; i++ {
		if version == 1 {
			val, err := uuid.NewUUID()
			if err != nil {
				slog.Error("Error creating UUID: ", err)
			}
			ids[i] = val.String()
		}

		if version == 4 {
			val, err := uuid.NewRandom()
			if err != nil {
				slog.Error("Error creating UUID: ", err)
			}
			ids[i] = val.String()
		}
	}

	start := time.Now()

	for i := 0; i < iterations; i++ {
		_, err = db.Exec("INSERT INTO TestTable (ID, Name) VALUES (?, ?)", ids[i], "")
		if err != nil {
			slog.Error("Error inserting row: ", err)
		}
	}

	since := time.Since(start)

	slog.Info(fmt.Sprintf("Time to insert rows: %s", since.String()))

}
