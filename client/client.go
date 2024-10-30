// client/client.go
package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"parsefi/config"
)

type ChainClient struct {
	Client *ethclient.Client
	Config config.ChainConfig
}

type MultiChainClient struct {
	clients map[string]*ChainClient
}

func NewMultiChainClient(configs []config.ChainConfig) (*MultiChainClient, error) {
	clients := make(map[string]*ChainClient)
	for _, cfg := range configs {
		client, err := ethclient.Dial(cfg.RPCUrl)
		if err != nil {
			log.Printf("连接 %s 链失败: %v", cfg.Name, err)
			return nil, err
		}
		clients[cfg.Name] = &ChainClient{Client: client, Config: cfg}
	}
	return &MultiChainClient{clients: clients}, nil
}

func (m *MultiChainClient) GetClient(chainName string) (*ChainClient, error) {
	client, ok := m.clients[chainName]
	if !ok {
		return nil, fmt.Errorf("未找到链 %s 的客户端", chainName)
	}
	return client, nil
}
