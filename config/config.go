package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

type Config struct {
	Listen struct {
		BindIP   string `yaml:"bind_ip" env-default:"127.0.0.1"`
		HttpPort string `yaml:"http_port" env-default:"8080"`
		GrpcPort string `yaml:"grpc_port" env-default:"8080"`
	} `yaml:"listen"`
	Storage         StorageConfig `yaml:"storage"`
	LevelDebug      string        `yaml:"level_debug"`
	TTLAccessToken  int           `yaml:"ttl_access_token"`
	TTLRefreshToken int           `yaml:"ttl_refresh_token"`
	SecretKey       string        `yaml:"secret_key"`
}

type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println("ERROR", help)
			log.Fatal(err)
		}

		loadErr := godotenv.Load()
		if loadErr != nil {
			log.Fatalln("can't load env file from current directory")
		}

		instance.Storage.Database = os.Getenv("DATABASE")
		instance.Storage.Username = os.Getenv("DB_USER")
		instance.Storage.Password = os.Getenv("DB_PASSWORD")
	})
	return instance
}
