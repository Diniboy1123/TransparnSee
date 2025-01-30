package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	CtLogURL         string   `json:"ctLogURL"`
	CommonNameSuffix string   `json:"commonNameSuffix"`
	OutputFile       string   `json:"outputFile"`
	BatchSize        int64    `json:"batchSize"`
	TrustedIssuers   []string `json:"trustedIssuers"`
	CloudflareNS     []string `json:"cloudflareNS"`
}

var AppConfig Config

func LoadConfig() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
}
