package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	rabbitConn *amqp.Connection
}

func main() {
	rabbitConn, err := connectRabbitMQ()
	if err != nil {
		log.Panic(err)
	}

	defer rabbitConn.Close()

	app := Config{
		rabbitConn: rabbitConn,
	}

	log.Printf("Starting broker service service on port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	var counts int
	var backOff = 1 * time.Second

	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Println("Error connecting to RabbitMQ. Total retries:", counts)
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			log.Println("Error connecting to RabbitMQ. Max retries reached")
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Retrying in", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
