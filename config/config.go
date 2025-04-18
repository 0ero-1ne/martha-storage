package config

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/ini.v1"
)

type ServerConfig struct {
	Hostname string `ini:"hostname"`
	Port     int    `ini:"port"`
}

type DropboxConfig struct {
	Token     string `ini:"token"`
	UsersPath string `ini:"users_path"`
	BookPath  string `ini:"books_path"`
}

type Config struct {
	ServerConfig  ServerConfig  `ini:"server"`
	DropboxConfig DropboxConfig `ini:"dropbox"`
}

func Init(configPath string) (*Config, error) {
	iniFile, err := ini.Load(configPath)

	if err != nil {
		return nil, err
	}

	var config Config
	err = iniFile.MapTo(&config)

	return &config, err
}

func (config ServerConfig) GetFullAddress() string {
	hostname := strings.TrimSpace(config.Hostname)
	port := config.Port

	if len(hostname) == 0 {
		log.Println("Hostname value was not provided in config file. Using default value \"0.0.0.0\"")
		hostname = "0.0.0.0"
	}

	if port == 0 {
		log.Println("Port value was not provided in config file. Using default value 8080")
		port = 8080
	}

	return fmt.Sprintf("%s:%d", hostname, port)
}
