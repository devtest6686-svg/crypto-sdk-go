package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lbe-io/crypto-sdk-go/callback"
	"github.com/lbe-io/crypto-sdk-go/client"
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

	// 2. 创建回调处理管理器
	manager := callback.NewCallbackManager(sdkClient, "http://example.com/xxxpath/callback")

	// 3. 注册转入成功事件处理器
	callback.Register(manager, entity.EventIncomingConfirmed, handleIncomingConfirmed)

	// 4. 注册转出成功事件处理器
	callback.Register(manager, entity.EventOutgoingConfirmed, handleOutgoingConfirmed)

	// 5. 注册转出失败事件处理器
	callback.Register(manager, entity.EventOutgoingFailed, handleOutgoingFailed)

	// 6. 注册keepalive事件处理器
	callback.Register(manager, entity.EventKeepalive, handleKeepaliveEvent)

	// 创建Gin处理对象
	ginHandler := callback.NewGinHandler(manager)

	g := gin.Default()
	g.POST("/callback", ginHandler.HandleFunc)
	fmt.Println("回调服务器启动在 http://localhost:8080")
	fmt.Println("回调地址: POST http://localhost:8080/callback")
	log.Fatal(g.Run(":8080"))

}

// handleIncomingConfirmed 处理转入成功事件
func handleIncomingConfirmed(ctx context.Context, data entity.Event) error {

	fmt.Println("trace_id:", ctx.Value("trace_id"))

	// 在这里添加你的业务逻辑
	// 例如：保存到数据库、发送通知、更新缓存等
	fmt.Printf("收到转入成功事件: %v \n", data.Json())
	fmt.Printf("  钱包ID: %s\n", data.WalletID)
	fmt.Printf("  交易ID: %s\n", data.Item.TransactionHash)
	fmt.Printf("  状态: %s\n", data.Item.Status)

	// 返回处理结果
	return nil
}

// handleOutgoingConfirmed 处理转出成功事件
func handleOutgoingConfirmed(ctx context.Context, data entity.Event) error {

	fmt.Println("trace_id:", ctx.Value("trace_id"))

	fmt.Printf("收到转出成功事件: %v \n", data.Json())
	fmt.Printf("  钱包ID: %s\n", data.WalletID)
	fmt.Printf("  交易ID: %s\n", data.Item.TransactionHash)

	// 返回处理结果
	return nil
}

// handleOutgoingFailed 处理转出失败事件
func handleOutgoingFailed(ctx context.Context, data entity.Event) error {

	fmt.Println("trace_id:", ctx.Value("trace_id"))

	fmt.Printf("收到转出失败事件: %v \n", data.Json())
	fmt.Printf("  钱包ID: %s\n", data.WalletID)
	fmt.Printf("  交易ID: %s\n", data.Item.TransactionHash)

	// 返回处理结果
	return nil
}

// handleKeepaliveEvent 处理keepalive事件
func handleKeepaliveEvent(ctx context.Context, data entity.KeepaliveEvent) error {

	fmt.Println("trace_id:", ctx.Value("trace_id"))

	fmt.Printf("收到keepalive事件: %s\n", data.Message)

	// 返回处理结果
	return nil
}
