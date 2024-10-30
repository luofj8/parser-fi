package parser

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ContractResult struct {
	BalanceOf   *big.Int `json:"balanceOf,omitempty"`
	TotalSupply *big.Int `json:"totalSupply,omitempty"`
	Symbol      string   `json:"symbol,omitempty"`
	Decimals    uint8    `json:"decimals,omitempty"`
}

func OutputJSON(response Response) {
	respJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("生成 JSON 响应失败:", err)
		return
	}
	fmt.Println(string(respJSON))
}
