package main

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"net/http"
	"os"
	"time"
	"user-microservice/internal/data"
)

const appVersion = "1.0.0"
const appName = "workout-users-microservice"

type config struct {
	appName string
	port    int
	db      struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	env string
}

type application struct {
	config   config
	logger   *log.Logger
	models   data.Models
	producer *kafka.Producer
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4001, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production")
	flag.StringVar(
		&cfg.db.dsn,
		"dsn",
		os.Getenv("WORKOUT_USERS_DB_DSN"),
		"Postgres dsn URI")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	cfg.appName = appName

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config:   cfg,
		logger:   logger,
		producer: nil,
	}

	app.routes()

	// we are now using gorm, so conn contains a pointer to a gorm.DB struct
	conn := app.connectDB()

	app.models = data.NewModels(conn)

	logger.Printf("database connection pool established")

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting server %s on port %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
	}

}
