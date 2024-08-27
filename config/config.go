package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MongoConnectionString string        `yaml:"mongo_connection_string" env:"MONGO_CONNECTION_STRING" env-default:"mongodb://localhost:27017"`
	MongoDatabase         string        `yaml:"mongo_database" env:"MONGO_DATABASE" env-default:"capgo"`
	S3BaseEndpoint        string        `yaml:"s3_base_endpoint" env:"S3_BASE_ENDPOINT"`
	S3Bucket              string        `yaml:"s3_bucket" env:"S3_BUCKET" env-required:"true"`
	ManagementAPITokens   []string      `yaml:"management_api_tokens" env:"MANAGEMENT_API_TOKENS" env-required:"true"`
	LimitRequestPerMinute int           `yaml:"limit_request_per_minute" env:"LIMIT_REQUEST_PER_MINUTE" env-default:"100"`
	TrustedProxies        []string      `yaml:"trusted_proxies" env:"TRUSTED_PROXIES"`
	CacheResultDuration   time.Duration `yaml:"cache_result_duration" env:"CACHE_RESULT_DURATION" env-default:"10m"`
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
