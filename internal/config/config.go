package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		App        App        `yaml:"app"`
		DB         DB         `yaml:"db"`
		Migrations Migrations `yaml:"migrations"`
	}

	App struct {
		Version     string `yaml:"version"`
		Name        string `yaml:"name" env:"APP_NAME"`
		Host        string `yaml:"host" env:"HOST"`
		Port        string `yaml:"port" env:"PORT"`
		Environment string `yaml:"environment" env:"ENVIRONMENT"`
	}

	DB struct {
		Host     string `yaml:"host" env:"CH_ANALYTIC_HOST"`
		Port     string `yaml:"port" env:"CH_ANALYTIC_PORT"`
		Username string `yaml:"username" env:"CH_ANALYTIC_USERNAME"`
		Password string `yaml:"password" env:"CH_ANALYTIC_PASSWORD"`
		Database string `yaml:"database" env:"CH_ANALYTIC_DB"`
	}

	Migrations struct {
		Dir string `yaml:"dir"`
	}
)

func New(configPath string) (*Config, error) {
	var (
		cfg  *Config
		err  error
		once sync.Once
	)

	once.Do(func() {
		cfg, err = parse(configPath)
	})

	return cfg, err
}

func parse(configPath string) (config *Config, err error) {
	filename, err := filepath.Abs(configPath)
	if err != nil {
		return
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	cfg := Config{}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return
	}

	if err = cleanenv.ReadConfig(filename, &cfg); err != nil {
		return
	}

	config = &cfg
	return
}
