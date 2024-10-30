package test

import (
	"github.com/ethereum/go-ethereum/common"
	"parsefi/config"
	"parsefi/parser"
	"testing"
)

func TestCallContractMethod(t *testing.T) {
	// 加载配置文件
	cfg, err := config.LoadConfig("../config/config.json")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 初始化解析器
	parser, err := parser.NewMultiChainParser(cfg)
	if err != nil {
		t.Fatalf("初始化解析器失败: %v", err)
	}

	// 测试调用合约方法
	result, err := parser.CallContractMethod("Ethereum", "ERC20", "balanceOf", []interface{}{common.HexToAddress("0xUserAddress")})
	if err != nil {
		t.Errorf("调用 balanceOf 方法失败: %v", err)
	} else {
		t.Logf("balanceOf 返回结果: %v", result)
	}
}