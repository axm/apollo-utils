package apollo

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     uint16 `json:"Port"`
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
	Queue      string     `json:"Queue"`
	Exchange   string     `json:"Exchange"`
	RoutingKey string     `json:"RoutingKey"`
	Durable    bool       `json:"Durable"`
	AutoDelete bool       `json:"AutoDelete"`
	Exclusive  bool       `json:"Exclusive"`
	NoWait     bool       `json:"NoWait"`
	Args       amqp.Table `json:"Args"`
}

type RabbitConsumerSettings struct {
	Queue     string      `json:"Queue"`
	Consumer  string      `json:"Consumer"`
	AutoAck   bool        `json:"AutoAck"`
	Exclusive bool        `json:"Exclusive"`
	NoLocal   bool        `json:"NoLocal"`
	NoWait    bool        `json:"NoWait"`
	Args      interface{} `json:"Args"`
}

func (rc RabbitConnection) CreateConnection() (*amqp.Connection, error) {
	connString := rc.GetConnectionString()
	fmt.Printf("rabbit conn string: %s\n", connString)
	conn, err := amqp.Dial(connString)
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
	RabbitConnection *RabbitConnection       `json:"RabbitConnection"`
	RabbitConsumer   *RabbitConsumerSettings `json:"RabbitConsumer"`
	Database         *DatabaseConnection     `json:"Database"`
	conn             *amqp.Connection
}

func (app *RabbitConsumerApp) Close() {
	if app.conn == nil {
		return
	}

	app.conn.Close()
}

func (app *RabbitConsumerApp) Consume() (<-chan amqp.Delivery, error) {
	conn, err := app.RabbitConnection.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to create rabbit connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("unable to create rabbit channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		app.RabbitConsumer.Queue,
		app.RabbitConsumer.AutoAck,
		app.RabbitConsumer.Exclusive,
		app.RabbitConsumer.NoLocal,
		app.RabbitConsumer.NoWait,
		nil)
	if err != nil {
		return nil, fmt.Errorf("unable to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		app.RabbitConsumer.Queue,
		app.RabbitConsumer.Consumer,
		app.RabbitConsumer.AutoAck,
		app.RabbitConsumer.Exclusive,
		app.RabbitConsumer.NoLocal,
		app.RabbitConsumer.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create consumer: %w", err)
	}

	return msgs, nil
}
