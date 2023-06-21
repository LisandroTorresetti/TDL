package config

import (
	"bot-telegram/utils"
	"fmt"
	"gopkg.in/yaml.v3"
)

const configFilepath = "./config.yaml"

type Endpoint struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

type GPTConfig struct {
	Model       string              `yaml:"model"`
	Temperature float64             `yaml:"temperature"`
	Endpoints   map[string]Endpoint `yaml:"endpoints"`
}

func LoadConfig() (*GPTConfig, error) {
	configFile, err := utils.GetConfigFileAsBytes(configFilepath)
	if err != nil {
		return nil, err
	}

	var gptConfig GPTConfig
	err = yaml.Unmarshal(configFile, &gptConfig)
	if err != nil {
		return nil, fmt.Errorf("error parsing GPT config file: %s", err)
	}

	return &gptConfig, nil
}
