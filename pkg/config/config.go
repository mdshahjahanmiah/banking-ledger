package config

import (
	"flag"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"os"
)

type Config struct {
	HttpAddress    string
	PostgresDSN    string
	MongoURI       string
	KafkaBrokerURL string
	LoggerConfig   logging.LoggerConfig
}

func Load() (Config, error) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	httpAddress := fs.String("http.address", "0.0.0.0:3000", "HTTP listen address for all specified endpoints.")
	mongoURI := fs.String("mongo.uri", os.Getenv("MONGO_URI"), "MongoDB connection URI")
	postgresDSN := fs.String("dsn", os.Getenv("POSTGRES_DSN"), "DB address")
	kafkaBroker := fs.String("kafka.broker", os.Getenv("KAFKA_BROKER_URL"), "Kafka broker URL")

	loggerConfig := logging.LoggerConfig{}
	fs.StringVar(&loggerConfig.CommandHandler, "logger.handler.type", "json", "handler type e.g json, otherwise default will be text type")
	fs.StringVar(&loggerConfig.LogLevel, "logger.log.level", "debug", "log level wise logging with fatal log")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return Config{}, err
	}

	config := Config{
		HttpAddress:    *httpAddress,
		PostgresDSN:    *postgresDSN,
		MongoURI:       *mongoURI,
		KafkaBrokerURL: *kafkaBroker,
		LoggerConfig:   loggerConfig,
	}

	return config, nil
}
