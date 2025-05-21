package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Host          string     `json:"host"`
	Timeout       string     `json:"timeout"`
	CacheCapacity int        `json:"cacheCapacity"`
	MaxBodySize   int64      `json:"maxBodySize"`
	Logger        LoggerConf `json:"logger"`
	StorageDir    string     `json:"storageDir"`
}

type LoggerConf struct {
	Level  string `json:"level"`
	Output string `json:"output"`
}

func Load(configPath string) (*Config, error) {
	// Проверяем существование файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &cfg, nil
}

func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./configs/config.json", "path to config file")
	flag.Parse()

	return configPath, nil
}
