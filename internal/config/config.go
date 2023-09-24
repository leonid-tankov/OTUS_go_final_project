package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Grpc struct {
		Port string `yaml:"port"`
	} `yaml:"grpc"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`
	Kafka struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"kafka"`
	Probability float64 `yaml:"probability"`
}

func NewConfig(logger *logrus.Logger, path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		logger.Fatal(err)
	}
	viper.SetConfigType("yaml")
	viper.ReadConfig(file)
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		logger.Fatal(err)
	}

	return conf
}
