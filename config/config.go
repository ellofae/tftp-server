package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfiguration struct {
		Retries        string `yaml:"retries"`
		Timeout        string `yaml:"timeout"`
		Address        string `yaml:"address"`
		TFTP_directory string `yaml:"tftp_directory"`
	} `yaml:"ServerConfiguration"`
}

func ConfigureViper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read the configuration file, error: %s\n", err.Error())
	}
	log.Println("Config loaded successfully.")

	return v
}

func ParseConfig(v *viper.Viper) *Config {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatal("Unable to parse the configuration file.")
	}
	log.Println("Configuration file parsed successfully.")

	return cfg
}
