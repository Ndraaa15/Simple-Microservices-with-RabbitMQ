package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange creates a new exchange for the consumer
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}
