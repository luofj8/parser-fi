package main

import (
	"github.com/ethereum/go-ethereum/common"
	"log"
	"parsefi/config"
	"parsefi/parser"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化多链解析器
	parser, err := parser.NewMultiChainParser(cfg)
	if err != nil {
		log.Fatalf("初始化解析器失败: %v", err)
	}

	// 示例调用 Ethereum 链上 ERC20 合约的 balanceOf 方法
	result, err := parser.CallContractMethod("Ethereum", "EETHABI", "balanceOf", []interface{}{common.HexToAddress("0xUserAddress")})
	if err != nil {
		log.Fatalf("调用合约方法失败: %v", err)
	}
	log.Printf("Ethereum 链上用户余额: %v", result)

	result, err = parser.CallContractMethod("Ethereum", "EETH", "totalSupply", nil)
	if err != nil {
		log.Fatalf("调用合约方法失败: %v", err)
	}
	log.Printf("Ethereum 链上总供应量: %v", result)

	result, err = parser.CallContractMethod("Ethereum", "EETH", "symbol", nil)
	if err != nil {
		log.Fatalf("调用合约方法失败: %v", err)
	}
	log.Printf("Ethereum symbol: %v", result)

	// 示例调用 Arbitrum 链上 ERC20 合约的 totalSupply 方法
	//result, err = parser.CallContractMethod("Arbitrum", "ERC20", "totalSupply", nil)
	//if err != nil {
	//	log.Fatalf("调用合约方法失败: %v", err)
	//}
	//log.Printf("Arbitrum 链上总供应量: %v", result)
}
