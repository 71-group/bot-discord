package helper

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token       string `json:"token"`
	Application string `json:"application"`
	Prefix      string `json:"prefix"`
}

func ReadConfig() Config {
	file, err := os.Open("./token.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return Config{}
	}
	defer file.Close()
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return Config{}
	}
	return config
}

