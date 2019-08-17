package apollo

import "github.com/gorilla/mux"

const (
	KAFKA_BOOTSTRAP_SERVERS = "bootstrap.servers"
)

type KafkaConsumerConfig struct {
	BootstrapServers string
	GroupId          string
	AutoOffsetReset  string
}

type KafkaProducerConfig struct {
	BootstrapServers string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       uint8
}

type DefaultApp struct {
	SqlServer *DatabaseConnection
	Rabbit    *RabbitConnection
	Redis     *RedisConfig
	Router    *mux.Router
}
