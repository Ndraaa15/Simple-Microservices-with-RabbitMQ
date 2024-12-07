package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting authentication service at port %s", webPort)

	// TODO connect to DB
	conn := connectDB()
	if conn == nil {
		log.Panic("Error connecting to database")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	startTime := time.Now()

	for {
		if startTime.Add(1 * time.Minute).Before(time.Now()) {
			log.Println("Database connection timed out")
			return nil
		}

		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to database:", err)
			log.Println("Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		} else {
			log.Println("Connected to database")
			return connection
		}

	}
}
