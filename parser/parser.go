package parser

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"parsefi/config"
	"strings"
)

type MultiChainParser struct {
	clients   map[string]*ethclient.Client  // 存储不同链的客户端
	abis      map[string]abi.ABI            // 存储 ABI 实例
	contracts map[string]common.Address     // 存储合约地址
	configs   map[string]config.ChainConfig // 存储链的配置
}

// 初始化解析器
func NewMultiChainParser(cfg *config.Config) (*MultiChainParser, error) {
	clients := make(map[string]*ethclient.Client)
	abis := make(map[string]abi.ABI)
	contracts := make(map[string]common.Address)
	configs := make(map[string]config.ChainConfig)

	for _, chain := range cfg.Chains {
		client, err := ethclient.Dial(chain.RPCUrl)
		if err != nil {
			log.Printf("连接 %s 链失败: %v", chain.Name, err)
			return nil, err
		}
		clients[chain.Name] = client
		configs[chain.Name] = chain

		for _, contract := range chain.Contracts {
			address := common.HexToAddress(contract.Address)
			contracts[fmt.Sprintf("%s:%s", chain.Name, contract.Name)] = address

			abiData, err := ioutil.ReadFile(contract.AbiPath)
			if err != nil {
				return nil, fmt.Errorf("无法加载合约 %s 的 ABI: %v", contract.Name, err)
			}
			parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
			if err != nil {
				return nil, fmt.Errorf("解析 ABI 失败: %v", err)
			}
			abis[fmt.Sprintf("%s:%s", chain.Name, contract.Name)] = parsedABI
		}
	}

	return &MultiChainParser{
		clients:   clients,
		abis:      abis,
		contracts: contracts,
		configs:   configs,
	}, nil
}

// 通用合约调用方法
func (p *MultiChainParser) CallContractMethod(chainName, contractName, methodName string, params []interface{}) (interface{}, error) {
	client, ok := p.clients[chainName]
	if !ok {
		return nil, fmt.Errorf("未找到链 %s 的客户端", chainName)
	}

	contractKey := fmt.Sprintf("%s:%s", chainName, contractName)
	contractAddress, ok := p.contracts[contractKey]
	if !ok {
		return nil, fmt.Errorf("未找到链 %s 的合约 %s", chainName, contractName)
	}

	parsedABI, ok := p.abis[contractKey]
	if !ok {
		return nil, fmt.Errorf("未找到合约 %s 的 ABI", contractName)
	}

	callData, err := parsedABI.Pack(methodName, params...)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var output interface{}
	err = parsedABI.UnpackIntoInterface(&output, methodName, result)
	if err != nil {
		return nil, err
	}

	return output, nil
}
