package config

import (
	"encoding/json"
	"log"
	"os"
)

type sslConfig struct {
	Certificate string
	Key         string
}
type tokenConfig struct {
	Size   int
	Expire int
}
type mysqlConfig struct {
	User     string
	Host     string
	Pass     string
	Database string
}
type passwordConfig struct {
	Cost int
}

//Configuration is the Main configuration struct for chat options
type Configuration struct {
	Port             string
	SSLSettings      *sslConfig
	TokenSettings    tokenConfig
	MYSQLSettings    mysqlConfig
	PasswordSettings passwordConfig
}

var config = readConfig()

//Settings returns a pointer to the internal Configuration variable
func Settings() *Configuration {
	return config
}

func readConfig() *Configuration {
	file, _ := os.Open("conf.json")
	config := Configuration{}
	err := json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalf("Could not read configuration file: %s", err.Error())
	}
	return &config
}
