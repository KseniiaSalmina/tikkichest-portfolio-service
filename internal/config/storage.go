package config

type Storage struct {
	Database string `env:"STORAGE_DATABASE"`
	Postgres Postgres
	Mongo    Mongo
}
