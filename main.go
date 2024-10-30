package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"os"
	"parsefi/config"
	"parsefi/parser"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("请提供链名称、合约名称和用户地址，例如：Ethereum EETH 0xUserAddress")
	}
	chainName, contractName, userAddress := os.Args[1], os.Args[2], os.Args[3]

	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	multiChainParser, err := parser.NewMultiChainParser(cfg)
	if err != nil {
		log.Fatalf("初始化多链解析器失败: %v", err)
	}

	result := parser.ContractResult{}
	log.Println("参数：", chainName, contractName, userAddress)

	// 调用合约方法并赋值给result
	callContractMethodAndSetResult(multiChainParser, chainName, contractName, userAddress, &result)

	// 检查结果是否为空
	response := parser.Response{Code: 0, Message: "操作成功", Data: result}
	if result.BalanceOf == nil && result.TotalSupply == nil && result.Symbol == "" {
		response = parser.Response{Code: 1, Message: "合约或方法不存在", Data: nil}
	}
	parser.OutputJSON(response)
}

// 调用合约方法并设置结果
func callContractMethodAndSetResult(multiChainParser *parser.MultiChainParser, chainName, contractName, userAddress string, result *parser.ContractResult) {
	var err error

	// 查询余额
	result.BalanceOf, err = callBigIntMethod(multiChainParser, chainName, contractName, "balanceOf", common.HexToAddress(userAddress))
	if err != nil {
		log.Printf("获取余额失败: %v", err)
	}

	// 查询总供应量
	result.TotalSupply, err = callBigIntMethod(multiChainParser, chainName, contractName, "totalSupply", nil)
	if err != nil {
		log.Printf("获取总供应量失败: %v", err)
	}

	// 查询代币符号
	if symbol, err := multiChainParser.CallContractMethod(chainName, contractName, "symbol", nil); err == nil {
		result.Symbol = symbol.(string)
	} else {
		log.Printf("获取代币符号失败: %v", err)
	}

	// 查询精度
	if decimals, err := multiChainParser.CallContractMethod(chainName, contractName, "decimals", nil); err == nil {
		result.Decimals = decimals.(uint8)
	} else {
		log.Printf("获取代币精度失败: %v", err)
	}
}

// 统一调用返回*big.Int类型的合约方法
func callBigIntMethod(multiChainParser *parser.MultiChainParser, chainName, contractName, methodName string, param interface{}) (*big.Int, error) {
	params := []interface{}{}
	if param != nil {
		params = append(params, param)
	}

	res, err := multiChainParser.CallContractMethod(chainName, contractName, methodName, params)
	if err != nil {
		return nil, fmt.Errorf("%s 方法调用失败: %v", methodName, err)
	}
	return res.(*big.Int), nil
}
