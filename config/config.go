package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string       `mapstructure:"env"`
	RestServer  RestServer   `mapstructure:"restServer"`
	Sqlite      SqliteConfig `mapstructure:"sqlite"`
	S3          S3Config     `mapstructure:"s3"`
	AWS         AWSConfig    `mapstructure:"aws"`
	HttpClients HttpClients  `mapstructure:"httpClients"`
}

type RestServer struct {
	Port int `mapstructure:"port"`
}

type SqliteConfig struct {
	Path              string `mapstructure:"path"`
	Mode              string `mapstructure:"mode"`              // ro, rw, rwc, memory
	Cache             string `mapstructure:"cache"`             // shared or private
	Immutable         int    `mapstructure:"immutable"`         // 1: immutable, 0: read/write
	EnableForeignKeys string `mapstructure:"enableForeignKeys"` // ON or OFF
	JournalMode       string `mapstructure:"journalMode"`       // DELETE, TRUNCATE, PERSIST, MEMORY, WAL, OFF
	LockingMode       string `mapstructure:"lockingMode"`       // normal, exclusive
	Synchronous       string `mapstructure:"synchronous"`       // OFF, NORMAL, FULL
}

type S3Config struct {
	Region               string `mapstructure:"region"`
	AccessKeyID          string `mapstructure:"accessKeyID"`
	SecretKey            string `mapstructure:"secretKey"`
	PublicBucket         string `mapstructure:"publicBucket"`
	PrivateBucket        string `mapstructure:"privateBucket"`
	DefaultPath          string `mapstructure:"defaultPath"`
	Endpoint             string `mapstructure:"endpoint"`             // Endpoint for the S3 service. Useful for MinIO.
	UseSSL               bool   `mapstructure:"useSSL"`               // Whether to use SSL. By default, true for AWS S3, and depends on your MinIO setup.
	UsePathStyleEndpoint bool   `mapstructure:"usePathStyleEndpoint"` // Needed for MinIO, set to true. For AWS, it's typically false.
	S3BaseURL            string `mapstructure:"s3BaseURL"`            // Base URL for the S3 service. If empty, will use the S3 endpoint.
	CDNBaseURL           string `mapstructure:"cdnBaseURL"`           // Base URL for the CDN. If empty, will use the S3 endpoint.
}

type AWSConfig struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secret"`
}

type HttpClients struct {
	AWSCognito AWSCognitoConfig `mapstructure:"awsCognito"`
}

type AWSCognitoConfig struct {
}

func ProvideConfig() *Config {
	config := &Config{}

	viper.SetConfigName("config.yaml")
	viper.SetConfigFile("./config/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading config file", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln("Unable to decode into struct", err)
	}

	return config
}
