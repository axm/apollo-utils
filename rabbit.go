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

func (rc *RabbitConnection) Conn() (*amqp.Connection, error) {
	connString := rc.GetConnectionString()
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create rabbit connection: %w", err)
	}

	return conn, err
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

func (rc RabbitConnection) CreateConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rc.User, rc.Password, rc.Host, rc.Port))
	if err != nil {
		return nil, fmt.Errorf("error creating Rabbit connection: %w", err)
	}

	return conn, nil
}

func (rc RabbitConnection) PublishMessage(settings *RabbitPublisherSettings, contents *[]byte) error {
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

type RabbitConsumerApp struct {
	Rc  *RabbitConnection
	Rcs *RabbitConsumerSettings
	conn *amqp.Connection
}

func (app *RabbitConsumerApp) Close() {
	app.conn.Close()
}

func (app *RabbitConsumerApp) Consume() (<-chan amqp.Delivery, error) {
	conn, err := app.Rc.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to create rabbit connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("unable to create rabbit channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		app.Rcs.Queue,
		app.Rcs.AutoAck,
		app.Rcs.Exclusive,
		app.Rcs.NoLocal,
		app.Rcs.NoWait,
		nil)
	if err != nil {
		return nil, fmt.Errorf("unable to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		app.Rcs.Queue,
		app.Rcs.Consumer,
		app.Rcs.AutoAck,
		app.Rcs.Exclusive,
		app.Rcs.NoLocal,
		app.Rcs.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create consumer: %w", err)
	}

	return msgs, nil
}
