package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	connection *amqp.Connection
}

func NewEventPublisher(conn *amqp.Connection) (*Publisher, error) {
	publisher := &Publisher{
		connection: conn,
	}

	if err := publisher.setup(); err != nil {
		log.Println("Error setting up publisher:", err)
		return nil, err
	}

	return publisher, nil
}

func (p *Publisher) setup() error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}

func (p *Publisher) Publish(data string, key string) error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	log.Println("Publishing event:", data, "with key:", key)

	err = channel.Publish(
		"logs_topic",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	)

	if err != nil {
		log.Println("Error publishing message:", err)
		return err
	}

	return nil
}
