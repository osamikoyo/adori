package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	ProxyElement struct {
		Prefix     string `yaml:"prefix"`
		SelfPath   string `yaml:"self_path"`
		OutputAddr string `yaml:"output_addr"`
	}

	StatisticConfig struct {
		Write bool   `yaml:"write"`
		Dir   string `yaml:"dir"`
	}

	Cash struct {
		Use               bool `yaml:"use"`
		IntervalInSeconds int  `yaml:"interval"`
	}

	Defence struct {
		Use                 bool     `yaml:"use"`
		MaxRequestFromOneIP uint     `yaml:"max_request_from_ip"`
		BadRequestParts     []string `yaml:"bad_request_parts"`
		BlackList           []string `yaml:"black_list"`
	}

	Static struct {
		Dir          string   `yaml:"dir"`
		Prefix       string   `yaml:"prefix"`
		ExcludeFiles []string `yaml:"exclude_files"`
	}

	Config struct {
		ServiceName     string          `yaml:"service_name"`
		Addr            string          `yaml:"addr"`
		Regime          string          `yaml:"regime"`
		Production      bool            `yaml:"production"`
		Defence         Defence         `yaml:"defence"`
		Cash            Cash            `yaml:"cash"`
		Static          Static          `yaml:"static"`
		StatisticConfig StatisticConfig `yaml:"statistic"`
		ApiGateway      []ProxyElement  `yaml:"api_gateway"`
	}
)

func NewConfig(paths ...string) (*Config, error) {
	var (
		file *os.File
		err  error
	)

	for _, path := range paths {
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

	if len(cfg.ApiGateway) == 0 {
		cfg.ApiGateway = nil
	}

	return &cfg, nil
}

func (c *Config) GetExludeMap(fullFiles []string) map[string]bool {
	resp := make(map[string]bool)

	for _, path := range fullFiles {
		resp[path] = true
	}

	for _, path := range c.Static.ExcludeFiles {
		resp[path] = false
	}

	return resp
}
