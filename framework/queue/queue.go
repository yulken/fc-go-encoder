package queue

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	User              string
	Password          string
	Host              string
	Port              string
	VHost             string
	ConsumerQueueName string
	ConsumerName      string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	rabbitMQArgs := amqp.Table{}
	rabbitMQArgs["x-dead-letter-exchange"] = os.Getenv("RABBITMQ_DLX")

	rabbitMQ := RabbitMQ{
		User:              os.Getenv("RABBIT_DEFAULT_USER"),
		Password:          os.Getenv("RABBIT_DEFAULT_PASS"),
		Host:              os.Getenv("RABBIT_DEFAULT_HOST"),
		Port:              os.Getenv("RABBIT_DEFAULT_PORT"),
		VHost:             os.Getenv("RABBIT_DEFAULT_VHOST"),
		ConsumerQueueName: os.Getenv("RABBIT_CONSUMER_QUEUE_NAME"),
		ConsumerName:      os.Getenv("RABBIT_CONSUMER_NAME"),
		AutoAck:           false,
		Args:              rabbitMQArgs,
	}

	return &rabbitMQ
}

func (r *RabbitMQ) Connect() *amqp.Channel {
	dsn := "amqp://" + r.User + ":" + r.Password + "@" + r.Host + ":" + r.Port + "/" + r.VHost
	conn, err := amqp.Dial(dsn)
	failOnError(err, "Failed to connect to RabbitMQ")

	r.Channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	return r.Channel
}

func (r *RabbitMQ) Consume(messageChannel chan amqp.Delivery) {
	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName,
		true,
		false,
		false,
		false,
		r.Args,
	)

	failOnError(err, "failed to declare a queue")

	incomingMessage, err := r.Channel.Consume(
		q.Name,
		r.ConsumerName,
		r.AutoAck,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		for message := range incomingMessage {
			log.Info("Incoming new message")
			messageChannel <- message
		}

		log.Info("RabbitMQ channel closed")
		close(messageChannel)
	}()
}

func (r *RabbitMQ) Notify(message string, contentType string, exchange string, routingKey string) error {
	err := r.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		},
	)

	if err != nil {
		return err
	}

	return nil

}

func failOnError(err error, message string) {

}
