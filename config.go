package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Quality       string   `yaml:"quality"`
	Refresh       int      `yaml:"refresh"`
	SubscribedAll bool     `yaml:"subscribedAll"`
	Subscriptions []string `yaml:"subscriptions"`
}

var DefaultConfig = Config{
	Quality:       "720",
	Refresh:       10,
	SubscribedAll: true,
}

func initConfig() {
	// Default config
	config := DefaultConfig

	bytes, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalln("Failed to create config")
	}

	ioutil.WriteFile("config.yaml", bytes, 0644)

	// TODO: First run question
}

func (n *Notificator) GetConfig() Config {
	_, err := os.Stat("config.yaml")
	if os.IsNotExist(err) {
		initConfig()
		return DefaultConfig
	}

	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("Failed to read config")
	}

	config := Config{}
	err = yaml.Unmarshal(file, &config)

	return config
}

func (n *Notificator) SetConfig(config Config) error {
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return errors.New("Failed to save config")
	}

	ioutil.WriteFile("config.yaml", bytes, 0644)
	return nil
}

// TODO
func autorun() {

}
