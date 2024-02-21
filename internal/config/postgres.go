package config

type Postgres struct {
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASSWORD"`
	Host     string `env:"PG_HOST" envDefault:"localhost"`
	Port     uint16 `env:"PG_PORT" envDefault:"5432"`
	Database string `env:"PG_DATABASE"`
}
