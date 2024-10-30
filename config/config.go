package config

import (
	"encoding/json"
	"os"
)

type ContractConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	AbiPath string `json:"abiPath"`
}

type ChainConfig struct {
	Name      string           `json:"name"`
	RPCUrl    string           `json:"rpcUrl"`
	Contracts []ContractConfig `json:"contracts"`
}

type Config struct {
	Chains []ChainConfig `json:"chains"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
