package config

type Kafka struct {
	Host  string `env:"KAFKA_HOST" envDefault:"localhost"`
	Port  uint16 `env:"KAFKA_PORT" envDefault:"9092"`
	Topic string `env:"KAFKA_TOPIC"`
}
