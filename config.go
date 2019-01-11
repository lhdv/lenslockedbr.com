package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type PostgresConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Name string `json:"name"`
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig {
		Host: "192.168.56.101",
		Port: 5432,
		User: "developer",
		Password: "1234qwer",
		Name: "lenslockedbr_dev",
	}
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s " +
                                   "sslmode=disable", c.Host, c.Port,
                                   c.User, c.Name)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s " + 
                           "dbname=%s sslmode=disable", c.Host, c.Port,
                           c.User, c.Password, c.Name)
}

type Config struct {
	Port int `json:"port"`
	Env string `json:"env"`
	Pepper string `json:"pepper"`
	HMACKey string `json:"hmac_key"`

	Database PostgresConfig `json:"database"`
}

func DefaultConfig() Config {
	return Config {
		Port: 3000,
		Env: "dev",
		Pepper: "foobar",
		HMACKey: "secret-hmac-key", 
		Database: DefaultPostgresConfig(),
	}
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func LoadConfig(configReq bool) Config {
	// Open the config file
	f, err := os.Open(".config")
	if err != nil {
		if configReq {
			panic(err)
		}

		log.Println("Using the default config...")
		return DefaultConfig()		
	}

	var c Config

	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully loaded .config")

	return c
}


