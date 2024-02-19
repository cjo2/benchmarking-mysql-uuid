package main

import (
	"database/sql"
	"log/slog"

	"github.com/go-sql-driver/mysql"
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
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	slog.Info("connected")

}
