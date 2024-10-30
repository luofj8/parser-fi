package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"os"
	"parsefi/config"
	"parsefi/parser"
)

func main() {
	// 定义命令行参数
	chainType := flag.String("chain", "Ethereum", "链类型 (如 Ethereum, Arbitrum)")
	contractType := flag.String("contract", "EETHABI", "合约类型 (如 EETHABI)")
	userAddress := flag.String("user", "", "用户地址 (用于 balanceOf 方法)")
	flag.Parse()

	// 检查用户地址是否为空
	if *userAddress == "" {
		fmt.Println("请提供用户地址，例如 -user 0xUserAddress")
		flag.Usage()
		os.Exit(1)
	}

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

	// 使用用户输入调用合约的方法
	address := common.HexToAddress(*userAddress)

	// 调用 balanceOf 方法
	balance, err := parser.CallContractMethod(*chainType, *contractType, "balanceOf", []interface{}{address})
	if err != nil {
		fmt.Printf("balanceOf 方法调用失败: %v\n", err)
	} else {
		fmt.Printf("%s 链上用户 %s 的余额: %v\n", *chainType, *userAddress, balance)
	}

	// 调用 totalSupply 方法
	totalSupply, err := parser.CallContractMethod(*chainType, *contractType, "totalSupply", nil)
	if err != nil {
		fmt.Printf("totalSupply 方法调用失败: %v\n", err)
	} else {
		fmt.Printf("%s 链上合约的总供应量: %v\n", *chainType, totalSupply)
	}

	// 调用 symbol 方法
	symbol, err := parser.CallContractMethod(*chainType, *contractType, "symbol", nil)
	if err != nil {
		fmt.Printf("symbol 方法调用失败: %v\n", err)
	} else {
		fmt.Printf("%s 链上合约的符号: %v\n", *chainType, symbol)
	}

	select {}
}
