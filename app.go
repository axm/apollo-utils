package apollo

import "github.com/gorilla/mux"

const (
	KAFKA_BOOTSTRAP_SERVERS = "bootstrap.servers"
	KAFKA_GROUP_ID = "group.id"
	KAFKA_AUTO_OFFSET_RESET = "auto.offset.reset"
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
