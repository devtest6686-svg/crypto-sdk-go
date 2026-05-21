package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/lbe-io/crypto-sdk-go/callback"
	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/consts"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

func main() {
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

	// 2. 创建回调管理器
	manager := callback.NewCallbackManager(sdkClient, "")

	// 3. 注册交易创建事件处理器
	callback.Register(manager, consts.EventTypeTransactionCreated, handleTransactionCreated)

	// 4. 注册钱包创建事件处理器
	callback.Register(manager, consts.EventTypeWalletCreated, handleWalletCreated)

	// 5. 注册测试事件处理器
	callback.Register(manager, consts.EventTypeTest, handleTestEvent)

	// 6. 创建并启动HTTP服务器
	httpHandler := callback.NewHTTPHandler(manager)

	// 可以添加其他业务路由
	mux := http.NewServeMux()
	mux.Handle("/callback", httpHandler.HandlerFunc())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	fmt.Println("回调服务器启动在 http://localhost:8080")
	fmt.Println("回调地址: POST http://localhost:8080/callback")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// handleTransactionCreated 处理交易创建事件
func handleTransactionCreated(ctx context.Context, data entity.TransactionCreatedEvent) error {

	// 在这里添加你的业务逻辑
	// 例如：保存到数据库、发送通知、更新缓存等
	fmt.Printf("收到交易创建事件:\n")
	fmt.Printf("  交易ID: %s\n", data.TransactionID)
	fmt.Printf("  金额: %s %s\n", data.Amount, data.Currency)
	fmt.Printf("  状态: %s\n", data.Status)

	// 返回处理结果
	return nil
}

// handleWalletCreated 处理钱包创建事件
func handleWalletCreated(ctx context.Context, data entity.WalletCreatedEvent) error {

	fmt.Printf("收到钱包创建事件:\n")
	fmt.Printf("  钱包ID: %s\n", data.WalletID)
	fmt.Printf("  钱包名称: %s\n", data.WalletName)
	fmt.Printf("  币种: %s\n", data.Currency)

	// 返回处理结果
	return nil
}

// handleTestEvent 处理测试事件
func handleTestEvent(ctx context.Context, data entity.TestEvent) error {

	fmt.Printf("收到测试事件: %s\n", data.Message)

	// 返回处理结果
	return nil
}
