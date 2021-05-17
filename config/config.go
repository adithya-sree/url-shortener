package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Redis Redis `json:"redis"`
	App   App   `json:"app"`
}

type App struct {
	Port        string `json:"port"`
	ShortLength int    `json:"shortLength"`
}

type Redis struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	DB       int64  `json:"db"`
	Timeout  int64  `json:"timeout"`
	TTL      int    `json:"ttl"`
}

func ReadInConfig(filepath string) (*Config, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
