package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/viper"
)

const (
	privKeyPath           = "./name_of_private_key.pem" // openssl genrsa -out app.rsa keysize
	publicKeyPath         = "./name_of_public_key.pem"
	refreshPrivateKeyPath = "./refresh_private_key.pem"
	refreshPublicKeyPath  = "./refresh_public_key.pem"
)

type Config struct {
	DBUri                  string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri               string `mapstructure:"REDIS_URL"`
	Port                   string `mapstructure:"PORT"`
	AccessTokenPrivateKey  []byte
	AccessTokenPublicKey   []byte
	RefreshTokenPrivateKey []byte
	RefreshTokenPublicKey  []byte
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	GrpcServerAddress      string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	Origin                 string        `mapstructure:"CLIENT_ORIGIN"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

// ? Struct Config

func LoadConfig(path string) (config Config, err error) {
	//Private key
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		fmt.Println("Reading file error: ", err)
	}
	config.AccessTokenPrivateKey = signBytes
	//Public key
	publicKeyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		fmt.Println("Reading file error: ", err)
	}
	config.AccessTokenPublicKey = publicKeyBytes
	//Private refresh key
	privateRefreshKey, err := ioutil.ReadFile(refreshPrivateKeyPath)
	if err != nil {
		fmt.Println("Reading file error: ", err)
	}
	config.RefreshTokenPrivateKey = privateRefreshKey
	//Public refresh key
	publicRefreshKey, err := ioutil.ReadFile(refreshPublicKeyPath)
	if err != nil {
		fmt.Println("Reading file error: ", err)
	}
	config.RefreshTokenPublicKey = publicRefreshKey
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
