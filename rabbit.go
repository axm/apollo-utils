package apollo

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	User     string
	Password string
	Host     string
	Port     uint16
}

func (rc *RabbitConnection) GetConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", rc.User, rc.Password, rc.Host, rc.Port)
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
	PublishMessage(settings *RabbitPublisherSettings, contents *[]byte) error
}

func (rc *RabbitConnection) CreateConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rc.User, rc.Password, rc.Host, rc.Port))
	if err != nil {
		return nil, fmt.Errorf("error creating Rabbit connection: %w", err)
	}

	return conn, nil
}

func (rc *RabbitConnection) PublishMessage(settings *RabbitPublisherSettings, contents *[]byte) error {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.User, rc.Password, rc.Host, rc.Port)
	conn, error := amqp.Dial(connString)
	if error != nil {
		return fmt.Errorf("unable to connect to rabbit instance: %w", error)
	}
	defer conn.Close()

	ch, error := conn.Channel()
	if error != nil {
		return fmt.Errorf("unable to create rabbit channel: %w", error)
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
		return fmt.Errorf("error declaring queue: %w", error)
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
		return fmt.Errorf("error publishing message: %w", error)
	}

	return nil
}
