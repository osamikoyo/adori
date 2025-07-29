package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var ConfigfileNames = []string{"yadori.yml", "yadorin.yaml", "yadori-config.yaml", "yadori-config.yml"}

type (
	ProxyElement struct {
		SelfPath   string `yaml:"self_path"`
		OutputAddr string `yaml:"output_addr"`
	}

	StatisticConfig struct {
		Write bool   `yaml:"write"`
		Dir   string `yaml:"dir"`
	}

	Config struct {
		ServiceName     string          `yaml:"service_name"`
		DDoSDef         bool            `yaml:"ddos_def"`
		Commress        bool            `yaml:"compress"`
		ExcludeFiles    []string        `yaml:"exclude_files"`
		StatisticConfig StatisticConfig `yaml:"statistic"`
		ApiGateway      []ProxyElement  `yaml:"api_gateway"`
	}
)

func NewConfig() (*Config, error) {
	var (
		file *os.File
		err  error
	)

	for _, path := range ConfigfileNames {
		file, err = os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()
	}

	var cfg Config

	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetExludeMap(fullFiles []string) map[string]bool {
	resp := make(map[string]bool)

	for _, path := range fullFiles {
		resp[path] = true
	}

	for _, path := range c.ExcludeFiles {
		resp[path] = false
	}

	return resp
}
