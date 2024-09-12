package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"golangtest/internal/data"
	"log"
	"os"
)

type Config struct {
	Port      int
	SecretKey string
	DB        struct {
		DSN string
	}
}

type Application struct {
	logger *log.Logger
	config Config
	models data.Models
}

func main() {
	var config Config

	flag.IntVar(&config.Port, "port", 8080, "API server port")
	flag.StringVar(&config.DB.DSN, "db-dsn", "", "database connection string")
	flag.StringVar(&config.SecretKey, "secret-key", "foobar", "secret key")
	flag.Parse()

	// log
	logger := log.New(os.Stdout, "api: ", log.LstdFlags)

	// db
	db, err := openDB(config)
	if err != nil {
		logger.Fatal(err)
	}

	app := Application{
		config: config,
		logger: logger,
		models: data.NewModels(db),
	}

	logger.Fatal(app.serve())
}

func openDB(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DB.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
