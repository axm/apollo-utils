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

func NewRabbitError(isTransient bool, msg string, orig error) *RabbitError {
	return &RabbitError{
		IsTransient: isTransient,
		Msg:         msg,
		SourceError: orig,
	}
}

type RabbitError struct {
	IsTransient bool
	Msg         string
	SourceError error
}

func (err *RabbitError) Error() string {
	if err.IsTransient {
		return fmt.Sprintf("Transient error ocurred: %s", err.Msg)
	}

	return fmt.Sprintf("Permanent error ocurred: %s", err.Msg)
}

type RabbitPublisher struct {
	RabbitPublisherSettings *RabbitPublisherSettings `json:"RabbitPublisherSettings"`
	RabbitConnection        *RabbitConnection        `json:"RabbitConnection"`
	conn                    *amqp.Connection
	ch                      *amqp.Channel
	err                     *RabbitError
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

func (rp *RabbitPublisher) Close() {
	if rp.ch != nil {
		rp.ch.Close()
	}

	if rp.conn != nil {
		rp.conn.Close()
	}
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

func (rp *RabbitPublisher) Publish(contents *[]byte) error {
	if rp.err != nil && !rp.err.IsTransient {
		if !rp.err.IsTransient {
			return fmt.Errorf("publisher permanently closed: %w", rp.err)
		}

		rp.err = nil
	}

	if rp.conn == nil {
		conn, ch, rabbitErr := createConnection(rp.RabbitConnection, rp.RabbitPublisherSettings)
		if rabbitErr != nil {
			rp.err = rabbitErr
			return fmt.Errorf("unable to create connection: %w", rabbitErr)
		}
		rp.conn = conn
		rp.ch = ch
	}

	err := rp.ch.Publish(rp.RabbitPublisherSettings.Exchange,
		rp.RabbitPublisherSettings.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        *contents,
		})

	if err != nil {
		return fmt.Errorf("unable to publish message: %w", err)
	}

	return nil
}

func createConnection(rc *RabbitConnection, rps *RabbitPublisherSettings) (*amqp.Connection, *amqp.Channel, *RabbitError) {
	conn, err := amqp.Dial(rc.GetConnectionString())
	if err != nil {
		return nil, nil, NewRabbitError(true, "unable to create rabbit connection", err)
	}

	ch, err := conn.Channel()
	if err != nil {

	}

	_, err = ch.QueueDeclare(
		rps.Queue,      // name
		rps.Durable,    // durable
		rps.AutoDelete, // delete when usused
		rps.Exclusive,  // exclusive
		rps.NoWait,     // no-wait
		rps.Args,       // arguments
	)
	if err != nil {
		return nil, nil, &RabbitError{
			IsTransient: false,
			Msg:         "error declaring queue",
			SourceError: err,
		}
	}

	return conn, ch, nil
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
