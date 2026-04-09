package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogLevel     string   `json:"logLevel"`
	OutputPath   string   `json:"outputPath"`
	RepoTemplate string   `json:"repoTemplate"`
	DownloadRepo string   `json:"downloadRepo"`
	PackagePaths []string `json:"packagePaths"`
	RepoUrls     []string `json:"repoUrls"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel:   "info",
		OutputPath: "./output/",
	}
}

func (c *Config) Parse(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(c)
}
