package main

import (
	"encoding/json"
	"log"
	"os"
)
type SSL struct {
	Certificate string
	Key string
}
type token struct {
	size int
	expire int
}
type mysql struct {
	user string
	host string
	pass string
	db   string
}
type Configuration struct {
	Port string
	SSLSettings *SSL
	TokenSettings token
	MYSQLSettings mysql
}

var config = read_config()
func read_config() *Configuration{
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Could not read configuration file: ", err)
	}
	return &config
}
