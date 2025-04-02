package message_queue

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitRabbitMQ() error {
	var err error
	maxRetries := 10
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %s", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return fmt.Errorf("could not connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		"device_exchange", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func CreateQueue(deviceID string) error {
	q, err := ch.QueueDeclare(
		deviceID, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,            // queue name
		deviceID,          // routing key
		"device_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteQueue(deviceID string) error {
	_, err := ch.QueueDelete(
		deviceID, // queue name
		false,    // ifUnused
		false,    // ifEmpty
		false,    // noWait
	)
	if err != nil {
		return err
	}

	return nil
}

func PublishMessage(deviceID, message string) error {
	err := ch.Publish(
		"device_exchange", // exchange
		deviceID,          // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}

	return nil
}

func CloseRabbitMQ() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
