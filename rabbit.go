package apollo

import (
	"fmt"
	"github.com/streadway/amqp"
)


type RabbitConnection struct {
	User     string
	Password string
	Host     string
	Port     string
}

type RabbitQueueSettings struct {
	Name       string
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


func (rc *RabbitConnection) CreateConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rc.User, rc.Password, rc.Host, rc.Port))
	if err != nil {
		return nil, &Error{message: "Error creating Rabbit connection", err: err}
	}

	return conn, nil
}


func DeclareQueue(ch *amqp.Channel, qs *RabbitQueueSettings) (amqp.Queue, error) {
	queue, error := ch.QueueDeclare(
		qs.Name,       // name
		qs.Durable,    // durable
		qs.AutoDelete, // delete when usused
		qs.Exclusive,  // exclusive
		qs.NoWait,     // no-wait
		qs.Args,       // arguments
	)
	if error != nil {
		return queue, &Error{ message: "Error declaring queue", err: error, }
	}

	return queue, nil
}
