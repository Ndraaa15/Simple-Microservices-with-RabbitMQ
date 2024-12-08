package main

import (
	"listener-service/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connectRabbitMQ()
	if err != nil {
		log.Panic(err)
	}

	defer rabbitConn.Close()

	log.Println("Connected to RabbitMQ")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
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
