package config

type Config struct {
	Name string

	DB struct {
		DataSource string
	}

	Mongo struct {
		URI      string
		Database string
	}

	Kafka struct {
		Brokers []string
		Topic   string
		GroupID string
	}

	Log struct {
		Level  string
		Format string
	}
}
