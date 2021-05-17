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
	Admin       Admin  `json:"admin"`
	Port        string `json:"port"`
	ShortLength int    `json:"shortLength"`
}

type Admin struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username"`
	Password string `json:"password"`
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
