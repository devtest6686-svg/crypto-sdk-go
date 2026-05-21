package main

import (
	"context"
	"log"

	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
	"github.com/lbe-io/crypto-sdk-go/request"
)

func main() {

	ctx := context.Background()

	// 1. 初始化SDK客户端
	config := client.SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	}

	sdkClient, err := client.NewClient(config)
	if err != nil {
		log.Fatal("初始化SDK失败:", err)
	}

	req := &entity.TransactionDetailReq{}

	resp, err := request.NewRequest(sdkClient).TransactionDetail(ctx, req)
	if err != nil {
		log.Fatal("查询交易详情失败:", err)
	}
	log.Println("交易详情查询成功:", resp)

}
