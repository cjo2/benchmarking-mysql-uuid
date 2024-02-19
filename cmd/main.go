package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/segmentio/conf"
)

func main() {
	var config = struct {
		Iterations             int    `conf:"iterations" help:"Number of iterations to run"`
		Version                int    `conf:"version" help:"Version of UUID to use"`
		DbAddress              string `conf:"dbaddress" help:"Address of database"`
		Db                     string `conf:"db" help:"Database to connect to"`
		DbUsername             string `conf:"dbusername" help:"Username to connect to database"`
		DbPassword             string `conf:"dbpassword" help:"Password to connect to database"`
		DbAllowNativePasswords bool   `conf:"dballownativepasswords" help:"Allow native passwords"`
	}{
		// defaults should match docker-compose
		DbAddress:              "localhost",
		Db:                     "TestDb",
		DbUsername:             "root",
		DbPassword:             "admin",
		DbAllowNativePasswords: true,
	}

	conf.Load(&config)

	dbConfig := mysql.Config{
		User:                 config.DbUsername,
		Passwd:               config.DbPassword,
		Addr:                 config.DbAddress,
		DBName:               config.Db,
		AllowNativePasswords: config.DbAllowNativePasswords,
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
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

	iterations := config.Iterations

	version := config.Version

	ids := make([]string, iterations)

	for i := 0; i < iterations; i++ {

		switch version {
		case 1:
			val, err := uuid.NewUUID()
			if err != nil {
				slog.Error("Error creating UUID: ", err)
			}
			ids[i] = val.String()
		case 4:
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
