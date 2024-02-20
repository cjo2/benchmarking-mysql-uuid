package main

import (
	"benchmarking-mysql-uuid/internal"
	"database/sql"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/segmentio/conf"
)

func main() {
	var config = struct {
		Iterations int `conf:"iterations" help:"Number of iterations to run"`
		Version    int `conf:"version" help:"Version of UUID to use (1,4,6,7)"`

		DbAddress              string `conf:"address" help:"Address of database"`
		Db                     string `conf:"database" help:"Database to connect to"`
		DbUsername             string `conf:"username" help:"Username to connect to database"`
		DbPassword             string `conf:"password" help:"Password to connect to database"`
		DbAllowNativePasswords bool   `conf:"allow-native-passwords" help:"Allow native passwords"`
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
		closeErr := db.Close()
		if closeErr != nil {
			slog.Error("error trying to close db: ", closeErr)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	slog.Info("connected")

	iterations := config.Iterations
	version := config.Version

	stats := internal.Stats{}
	start := time.Now()

	// precomputed to remove uuid generation from the benchmark
	ids := make([]string, iterations)

	for i := 0; i < iterations; i++ {
		var val uuid.UUID
		var uuidErr error

		switch version {
		case 1:
			val, uuidErr = uuid.NewUUID()
		case 4:
			val, uuidErr = uuid.NewRandom()
		case 6:
			val, uuidErr = uuid.NewV6()
		case 7:
			val, uuidErr = uuid.NewV7()
		}

		if uuidErr != nil {
			slog.Error("Error creating UUID: ", err)
		}

		ids[i] = val.String()
	}

	go func() {
		lastSuccessfulInserts := 0
		lastTime := time.Now()

		for {
			select {
			case <-time.Tick(3 * time.Second):
				successfulInserts := stats.GetSuccessfulInserts()
				now := time.Now()

				// solve for how much time it took for the last batch of inserted rows on average
				rowsPerSecond := float64(successfulInserts-lastSuccessfulInserts) / (now.Sub(lastTime).Seconds())

				lastSuccessfulInserts = successfulInserts
				lastTime = now

				slog.With("success", successfulInserts, "failure", stats.GetFailedInserts(), "rowsPerSecond", rowsPerSecond).Info("stats", "timeElapsed", time.Since(start).String())
			}
		}
	}()

	for _, id := range ids {
		if _, err = db.Exec("INSERT INTO TestTable (ID, Name) VALUES (?, ?)", id, ""); err != nil {
			slog.Error("Error inserting row: ", err)
			stats.IncrementFailedInserts()
			continue
		}

		stats.IncrementSuccessfulInserts()
	}

	slog.With("success", stats.GetSuccessfulInserts(), "failure", stats.GetFailedInserts()).Info("stats", "timeElapsed", time.Since(start).String())
}
