package apollo

import "github.com/gorilla/mux"

const (
	KAFKA_BOOTSTRAP_SERVERS = "bootstrap.servers"
)

type KafkaConsumerConfig struct {
	bootstrapServers string
	groupId          string
	autoOffsetReset  string
}

type KafkaProducerConfig struct {
	bootstrapServers string
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
