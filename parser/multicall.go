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
	clients   map[string]*ethclient.Client // 存储不同链的客户端
	abis      map[string]abi.ABI           // 存储 ABI 实例
	contracts map[string]common.Address    // 存储合约地址
}

// 初始化多链多合约解析器
func NewMultiChainParser(cfg *config.Config) (*MultiChainParser, error) {
	clients := make(map[string]*ethclient.Client)
	abis := make(map[string]abi.ABI)
	contracts := make(map[string]common.Address)

	for _, chain := range cfg.Chains {
		client, err := ethclient.Dial(chain.RPCUrl)
		if err != nil {
			log.Printf("连接 %s 链失败: %v", chain.Name, err)
			return nil, err
		}
		clients[chain.Name] = client

		for _, contract := range chain.Contracts {
			// 加载 ABI
			abiData, err := ioutil.ReadFile(contract.AbiPath)
			if err != nil {
				return nil, fmt.Errorf("加载合约 %s 的 ABI 失败: %v", contract.Name, err)
			}
			parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
			if err != nil {
				return nil, fmt.Errorf("解析合约 %s 的 ABI 失败: %v", contract.Name, err)
			}
			abis[fmt.Sprintf("%s:%s", chain.Name, contract.Name)] = parsedABI
			contracts[fmt.Sprintf("%s:%s", chain.Name, contract.Name)] = common.HexToAddress(contract.Address)
		}
	}

	return &MultiChainParser{
		clients:   clients,
		abis:      abis,
		contracts: contracts,
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

	// 构建调用数据
	callData, err := parsedABI.Pack(methodName, params...)

	if err != nil {
		log.Println("Pack：", err)
		return nil, err
	}

	// 执行合约调用
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Println("CallContract,err：", err)
		return nil, err
	}

	// 解析返回结果
	var output interface{}
	err = parsedABI.UnpackIntoInterface(&output, methodName, result)
	if err != nil {
		log.Println("UnpackIntoInterface，err：", err)
		return nil, err
	}

	return output, nil
}
