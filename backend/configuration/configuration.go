package configuration

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Configutration struct {
	ListenAddress        string `env:"LISTEN_ADDR"`
	MinioEndpoint        string `env:"MINIO_ENDPOINT"`
	MinioAccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
	MinioSecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
	RedisEndpoint        string `env:"REDIS_ENDPOINT"`
	RedisPassword        string `env:"REDIS_PASSWORD"`
	MongoEndpoint        string `env:"MONGO_ENDPOINT"`
	MongoDb              string `env:"MONGO_DB"`
	CasdoorEndpoint      string `env:"CASDOOR_ENDPOINT"`
	CasdoorClientID      string `env:"CASDOOR_CLIENT_ID"`
	CasdoorClientSecret  string `env:"CASDOOR_CLIENT_SECRET"`
	CasdoorCertificate   string `env:"CASDOOR_CERTIFICATE"`
	CasdoorOrganization  string `env:"CASDOOR_ORGANIZATION"`
	CasdoorApplication   string `env:"CASDOOR_APPLICATION"`
	CallbackAddress      string `env:"CALLBACK_ADDRESS"`
}

func SetConfiguration() Configutration {
	var cfg Configutration
	env.Parse(&cfg)

	flag.StringVar(&cfg.ListenAddress, "LISTEN_ADDR", cfg.ListenAddress, "Port on where the Application should be running.")
	flag.StringVar(&cfg.MinioEndpoint, "MINIO_ENDPOINT", cfg.MinioEndpoint, "Address to the Minio service.")
	flag.StringVar(&cfg.MinioAccessKeyID, "MINIO_ACCESS_KEY_ID", cfg.MinioAccessKeyID, "Minio access key ID.")
	flag.StringVar(&cfg.MinioSecretAccessKey, "MINIO_SECRET_ACCESS_KEY", cfg.MinioSecretAccessKey, "Minio secret access key.")
	flag.StringVar(&cfg.RedisEndpoint, "REDIS_ENDPOINT", cfg.RedisEndpoint, "Redis Endpoint.")
	flag.StringVar(&cfg.RedisPassword, "REDIS_PASSWORD", cfg.RedisPassword, "Redis Password.")
	flag.StringVar(&cfg.MongoEndpoint, "MONGO_ENDPOINT", cfg.MongoEndpoint, "Mongo DB Endpoint.")
	flag.StringVar(&cfg.MongoDb, "MONGO_DB", cfg.MongoDb, "Mongo database.")
	flag.StringVar(&cfg.CasdoorEndpoint, "CASDOOR_ENDPOINT", cfg.CasdoorEndpoint, "Endpoint to casdoor service.")
	flag.StringVar(&cfg.CasdoorClientID, "CASDOOR_CLIENT_ID", cfg.CasdoorClientID, "Casdoor client ID.")
	flag.StringVar(&cfg.CasdoorClientSecret, "CASDOOR_CLIENT_SECRET", cfg.CasdoorClientSecret, "Casdoor client secret.")
	flag.StringVar(&cfg.CasdoorCertificate, "CASDOOR_CERTIFICATE", cfg.CasdoorCertificate, "Casdoor certificate.")
	flag.StringVar(&cfg.CasdoorOrganization, "CASDOOR_ORGANIZATION", cfg.CasdoorOrganization, "Casdoor organization.")
	flag.StringVar(&cfg.CasdoorApplication, "CASDOOR_APPLICATION", cfg.CasdoorApplication, "Casdoor application.")
	flag.StringVar(&cfg.CallbackAddress, "CALLBACK_ADDRESS", cfg.CallbackAddress, "callback address")

	flag.Parse()

	return cfg
}
