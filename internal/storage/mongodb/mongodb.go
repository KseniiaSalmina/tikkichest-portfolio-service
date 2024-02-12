package mongodb

type Database struct {
	db mongo.Client
}

func NewDB() {
	client, err := mongo.Connect
}
