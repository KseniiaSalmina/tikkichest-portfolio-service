package config

type Mongo struct {
	User                   string `env:"MG_USERNAME"`
	Password               string `env:"MG_PASSWORD"`
	Host                   string `env:"MG_HOST" envDefault:"localhost"`
	Port                   int    `env:"MG_PORT" envDefault:"27017"`
	AuthenticationDatabase string `env:"MG_AUTH_DATABASE"`
	Database               string `env:"MG_DATABASE"`
	Collection             string `env:"MG_COLLECTION"`
	HaveIndexes            bool   `env:"MG_HAVE_INDEXES"`
}
