package configuration

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configutration struct {
	AppAddress           string `env:"APP_ADDR" yaml:"APP_ADDR"`
	ListenAddress        string `env:"LISTEN_ADDR" yaml:"LISTEN_ADDR"`
	MinioEndpoint        string `env:"MINIO_ENDPOINT" yaml:"MINIO_ENDPOINT"`
	MinioAccessKeyID     string `env:"MINIO_ACCESS_KEY_ID" yaml:"MINIO_ACCESS_KEY_ID"`
	MinioSecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY" yaml:"MINIO_SECRET_ACCESS_KEY"`
	CasdoorEndpoint      string `env:"CASDOOR_ENDPOINT" yaml:"CASDOOR_ENDPOINT"`
	CasdoorClientID      string `env:"CASDOOR_CLIENT_ID" yaml:"CASDOOR_CLIENT_ID"`
	CasdoorClientSecret  string `env:"CASDOOR_CLIENT_SECRET" yaml:"CASDOOR_CLIENT_SECRET"`
	CasdoorCertificate   string `env:"CASDOOR_CERTIFICATE" yaml:"CASDOOR_CERTIFICATE"`
	CasdoorOrganization  string `env:"CASDOOR_ORGANIZATION" yaml:"CASDOOR_ORGANIZATION"`
	CasdoorApplication   string `env:"CASDOOR_APPLICATION" yaml:"CASDOOR_APPLICATION"`
	CallbackAddress      string `env:"CALLBACK_ADDRESS" yaml:"CALLBACK_ADDRESS"`

	BadgerPath string `env:"BADGER_PATH" yaml:"BADGER_PATH"`
}

func SetConfiguration() Configutration {
	var cfg Configutration
	err := godotenv.Load()
	if err != nil {
		fmt.Println("no .env found")
	}
	err = env.Parse(&cfg)
	if err != nil {
		fmt.Println(".env file could not be parsed")
	}

	flag.StringVar(&cfg.AppAddress, "APP_ADDR", cfg.AppAddress, "Address of the app")
	flag.StringVar(&cfg.ListenAddress, "LISTEN_ADDR", cfg.ListenAddress, "Port on where the Application should be running.")
	flag.StringVar(&cfg.MinioEndpoint, "MINIO_ENDPOINT", cfg.MinioEndpoint, "Address to the Minio service.")
	flag.StringVar(&cfg.MinioAccessKeyID, "MINIO_ACCESS_KEY_ID", cfg.MinioAccessKeyID, "Minio access key ID.")
	flag.StringVar(&cfg.MinioSecretAccessKey, "MINIO_SECRET_ACCESS_KEY", cfg.MinioSecretAccessKey, "Minio secret access key.")
	flag.StringVar(&cfg.CasdoorEndpoint, "CASDOOR_ENDPOINT", cfg.CasdoorEndpoint, "Endpoint to casdoor service.")
	flag.StringVar(&cfg.CasdoorClientID, "CASDOOR_CLIENT_ID", cfg.CasdoorClientID, "Casdoor client ID.")
	flag.StringVar(&cfg.CasdoorClientSecret, "CASDOOR_CLIENT_SECRET", cfg.CasdoorClientSecret, "Casdoor client secret.")
	flag.StringVar(&cfg.CasdoorCertificate, "CASDOOR_CERTIFICATE", cfg.CasdoorCertificate, "Casdoor certificate.")
	flag.StringVar(&cfg.CasdoorOrganization, "CASDOOR_ORGANIZATION", cfg.CasdoorOrganization, "Casdoor organization.")
	flag.StringVar(&cfg.CasdoorApplication, "CASDOOR_APPLICATION", cfg.CasdoorApplication, "Casdoor application.")
	flag.StringVar(&cfg.CallbackAddress, "CALLBACK_ADDRESS", cfg.CallbackAddress, "callback address")

	flag.StringVar(&cfg.BadgerPath, "BADGER_PATH", cfg.BadgerPath, "Badger Path")

	flag.Parse()

	return cfg
}
