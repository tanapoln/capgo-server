package config

import (
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MongoConnectionString string `yaml:"mongo_connection_string" env:"MONGO_CONNECTION_STRING" env-default:"mongodb://localhost:27017"`
	MongoDatabase         string `yaml:"mongo_database" env:"MONGO_DATABASE" env-default:"capgo"`
	S3Bucket              string `yaml:"s3_bucket" env:"S3_BUCKET" env-default:"capgo-bundle"`
}

var (
	cfg *Config
)

func init() {
	cfg = &Config{}
	err := Reload()
	if err != nil {
		slog.Error("Config error", "error", err)
		os.Exit(1)
	}
}

func Get() Config {
	return *cfg
}

func Reload() error {
	return cleanenv.ReadConfig("config.yml", cfg)
}
