package apollo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	User     string
	Password string
	Host     string
	Port     uint16
}

type RabbitPublisherSettings struct {
	Queue      string
	Exchange   string
	RoutingKey string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type RabbitConsumerSettings struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      interface{}
}

type RabbitPublisher interface {
	PublishMessage(rc *RabbitConnection, settings *RabbitPublisherSettings, contents *[]byte) error
}

func (rc *RabbitConnection) CreateConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rc.User, rc.Password, rc.Host, rc.Port))
	if err != nil {
		return nil, &Error{message: "Error creating Rabbit connection", err: err}
	}

	return conn, nil
}

func PublishMessage(rc *RabbitConnection, settings *RabbitPublisherSettings, contents *[]byte) error {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.User, rc.Password, rc.Host, rc.Port)
	conn, error := amqp.Dial(connString)
	if error != nil {
		return errors.Wrap(error, "Unable to connect to rabbit instance")
	}
	defer conn.Close()

	ch, error := conn.Channel()
	if error != nil {
		return errors.Wrap(error, "Unable to create rabbit channel")
	}
	defer ch.Close()

	_, error = ch.QueueDeclare(
		settings.Queue,      // name
		settings.Durable,    // durable
		settings.AutoDelete, // delete when usused
		settings.Exclusive,  // exclusive
		settings.NoWait,     // no-wait
		settings.Args,       // arguments
	)
	if error != nil {
		return errors.Wrap(error, "Error declaring queue.")
	}

	error = ch.Publish(settings.Exchange,
		settings.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        *contents,
		})
	if error != nil {
		return errors.Wrap(error, "Error publishing message.")
	}

	return nil
}
