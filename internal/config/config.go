package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/smolathon/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"` // bool обязательно ссылкой
	Listen  struct {
		Type       string `yaml:"type" env-default:"port"`
		BindIP     string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port       string `yaml:"port" env-default:"8080"`
		SocketFile string `yaml:"socket_file" env-default:"app.sock"`
	} `yaml:"listen"`
	MongoDB struct {
		Host            string `json:"host" yaml:"host"`
		Port            string `json:"port" yaml:"port"`
		Database        string `json:"database" yaml:"database"`
		AuthDB          string `json:"auth_db" yaml:"auth_db"`
		Username        string `json:"username" yaml:"username"`
		Password        string `json:"password" yaml:"password"`
		MasterColletion string `yaml:"masterCollection"`
		CardsCollecion  string `yaml:"cardsCollection"`
	} `json:"mongodb" yaml:"mongodb"`
	AppConfig struct {
		LogLevel  string `yaml:"log_level" env-default:"trace"`
		AdminUser struct {
			Email    string `yaml:"admin_email" env-default:"admin"`
			Password string `yaml:"admin_pwd" env-default:"12345"`
		}
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger("trace")
		logger.Info("read app configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
