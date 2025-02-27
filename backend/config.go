package bilicomicdownloader

import (
	"encoding/json"
	"fmt"
	"os"
)

const configPath = "bcconfig.json"

type Config struct {
	UrlBase     string `json:"urlBase"`
	OutputPath  string `json:"outputPath"`
	PackageType string `json:"packageType"`
	ImageFormat string `json:"imageFormat"`
	NamingStyle string `json:"namingStyle"`
	Cookie      string `json:"cookie"`
}

var ConfigInstance *Config = Load()

func NewConfig() *Config {
	return &Config{
		UrlBase:     "www.bilicomic.net",
		OutputPath:  "./",
		PackageType: "cbz",
		ImageFormat: "source",
		NamingStyle: "title",
		Cookie:      "",
	}
}

func Load() *Config {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return NewConfig()
	}

	var config *Config = &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		fmt.Println("Error deserializing config:", err)
		return NewConfig()
	}

	return config
}

func (c *Config) Save() {
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println("Error serializing config:", err)
		return
	}

	err = os.WriteFile(configPath, content, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
	}
}

func (c *Config) SaveConfig(config *Config) {
	*ConfigInstance = *config
	ConfigInstance.Save()
}

func (c *Config) GetConfig() *Config {
	return ConfigInstance
}
