package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func (app *application) connectDB() *gorm.DB {
	conn, err := sql.Open("postgres", app.config.db.dsn)
	if err != nil {
		app.logger.Println("Error while connecting to database: ", err)
		return nil
	}

	conn.SetMaxOpenConns(app.config.db.maxOpenConns)
	conn.SetMaxIdleConns(app.config.db.maxIdleConns)
	duration, err := time.ParseDuration(app.config.db.maxIdleTime)
	if err != nil {
		return nil
	}

	conn.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.PingContext(ctx)
	if err != nil {
		app.logger.Println(err)
		return nil
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})
	if err != nil {
		app.logger.Println("error occurred while initializing gorm with existing db connection")
	}

	return gormDb
}
