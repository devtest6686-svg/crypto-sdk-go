# crypto-sdk-go

Go 语言版加密业务 SDK，提供：

- 服务端 API 请求签名与响应验签
- 交易与钱包相关接口封装
- 回调通知验签与事件分发
- Gin / net/http 两种回调接入方式

模块地址：`github.com/lbe-io/crypto-sdk-go`

## 1. 功能简介

### 1.1 API 请求能力

通过 `request.IRequest` 暴露统一调用能力，当前包含：

- 交易详情查询：`TransactionDetail`
- 钱包列表查询：`WalletList`
- 钱包地址生成：`WalletAddressGenerate`
- 钱包地址列表：`WalletAddressList`
- 钱包地址资产：`WalletAddressAssets`

内部由 `client.CryptoClient` 统一处理：

- 请求体 JSON 序列化
- `Authorization` 签名生成（Crypto-V1）
- HTTP 发起请求
- 响应签名校验
- 响应 JSON 反序列化

### 1.2 回调通知能力

通过 `callback.CallbackManager` 提供：

- 回调请求签名校验
- 事件类型识别与数据解析
- 泛型安全处理函数注册与调用

支持两种适配器：

- `callback.HTTPHandler`（标准库 `net/http`）
- `callback.GinHandler`（Gin 框架）

### 1.3 签名能力

`sdk/signx` 提供 Crypto-V1 签名算法：

- HMAC-SHA256
- 签名时间窗默认 ±5 分钟
- CanonicalURL 规范化（去掉 Host，仅保留 RequestURI）
- `Authorization` 头格式统一为 `Crypto-V1 ...`

---

## 2. 项目结构

```text
crypto-sdk-go/
├── client/                 # SDK 客户端（初始化、请求、验签、选项）
├── request/                # 业务 API 请求封装
├── callback/               # 回调管理与 HTTP/Gin 适配
├── common/
│   ├── consts/             # API 路径、回调事件常量
│   └── entity/             # 请求/响应与回调事件实体
├── sdk/
│   ├── signx/              # 签名与验签
│   └── logger/             # 日志接口与默认实现
└── examples/               # 可运行示例
```

---

## 3. 如何接入

### 3.1 环境要求

- Go >= 1.25（以 `go.mod` 为准）

### 3.2 安装依赖

```bash
go get github.com/lbe-io/crypto-sdk-go
```

### 3.3 准备配置

初始化客户端需要 4 个必填参数（缺失会报错）：

- `ServerHost`：服务端主机地址，例如 `https://api.example.com`
- `MerchantID`：商户 ID
- `APIKey`：应用 Key
- `APISecret`：签名密钥

### 3.4 初始化客户端

```go
package main

import (
	"log"
	"time"

	"github.com/lbe-io/crypto-sdk-go/client"
)

func main() {
	cfg := client.SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	}

	sdkClient, err := client.NewClient(
		cfg,
		client.WithTimeout(20*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = sdkClient
}
```

可选配置项：

- `client.WithTimeout(timeout)`：请求超时（默认 30s）
- `client.WithHeaderKV(key, value)` / `client.WithHeaders(map[string]string)`：预留 Header 扩展配置
- `client.WithIgnoreLog()`：忽略内部请求/响应日志
- `client.WithLogger(log)`：替换默认日志实现

> 说明：当前版本 `WithHeaderKV` 和 `WithHeaders` 的 Header 数据已保存在 `clientOptions` 中，但尚未注入到实际请求头中。

---

## 4. 如何使用

### 4.1 发起 API 请求

```go
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

	sdkClient, err := client.NewClient(client.SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	})
	if err != nil {
		log.Fatal(err)
	}

	api := request.NewRequest(sdkClient)

	resp, err := api.TransactionDetail(ctx, &entity.TransactionDetailReq{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("transaction detail response: %+v", resp)
}
```

当前 SDK 业务接口映射为 POST 请求：

- `TransactionDetail` -> `POST /transaction/detail`
- `WalletList` -> `POST /wallet/list`
- `WalletAddressGenerate` -> `POST /wallet/address/generate`
- `WalletAddressList` -> `POST /wallet/address/list`
- `WalletAddressAssets` -> `POST /wallet/address/assets`

SDK 也提供通用请求方法：

- `CryptoClient.Get(ctx, path, req, resp)`
- `CryptoClient.Post(ctx, path, req, resp)`

### 4.2 接入回调（net/http）

```go
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
	sdkClient, err := client.NewClient(client.SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	})
	if err != nil {
		log.Fatal(err)
	}

	manager := callback.NewCallbackManager(sdkClient, "/callback")

	callback.Register(manager, consts.EventTypeTransactionCreated,
		func(ctx context.Context, data entity.TransactionCreatedEvent) error {
			fmt.Printf("tx created: %s\n", data.TransactionID)
			return nil
		})

	h := callback.NewHTTPHandler(manager)
	mux := http.NewServeMux()
	mux.Handle("/callback", h.HandlerFunc())

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### 4.3 接入回调（Gin）

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lbe-io/crypto-sdk-go/callback"
	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/consts"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

func main() {
	sdkClient, err := client.NewClient(client.SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	})
	if err != nil {
		log.Fatal(err)
	}

	manager := callback.NewCallbackManager(sdkClient, "/callback")

	callback.Register(manager, consts.EventTypeTest,
		func(ctx context.Context, data entity.TestEvent) error {
			fmt.Println("test callback:", data.Message)
			return nil
		})

	g := gin.Default()
	g.POST("/callback", callback.NewGinHandler(manager).HandleFunc)
	log.Fatal(g.Run(":8080"))
}
```

---

## 5. 回调事件类型

可注册事件：

- `INCOMING_CONFIRMED`
- `OUTGOING_CONFIRMED`
- `OUTGOING_FAILED`
- `KEEPALIVE`

对应常量位于 `common/consts/callback.go`和 `common/consts/event.go`：

- `consts.Event`

---

## 6. 签名与安全说明

### 6.1 请求签名

SDK 在每次请求前自动生成 `Authorization`，签名内容包含：

- merchant_id
- appid（即 APIKey）
- method
- nonce
- sign_time（毫秒）
- canonical_url
- sign_body

签名 Header 格式为：

`Crypto-V1 merchant_id:xxx&appid:xxx&nonce:xxx&sign_time:xxx&sign:xxx`

### 6.2 响应验签

SDK 会读取响应头 `Authorization` 并校验签名。

- 若响应头缺失，返回错误
- 若验签失败，返回错误

### 6.3 回调验签

`callback.CallbackManager` 的 `HandleRequest` 会对回调请求进行验签。

- 默认使用回调请求的 `RequestURI`
- 若网关或代理改写了 URI，请显式调用 `callback.NewCallbackManager(sdkClient, callbackUri)`

---

## 7. 示例程序

仓库内已提供 3 个示例：

- `examples/request`：基础 API 调用
- `examples/http_callback`：`net/http` 回调接入
- `examples/gin_callback`：Gin 回调接入

运行方式：

```bash
go run ./examples/request
go run ./examples/http_callback
go run ./examples/gin_callback
```

---

## 8. 注意事项

- 初始化报错：优先检查 4 个必填配置项是否正确设置。
- `WithHeaderKV/WithHeaders` 目前未注入到实际 HTTP 请求头，如需附加自定义 Header，需补充 `client.doRequest` 实现。
- 签名时间窗为 ±5 分钟，请确保客户端与服务端时钟基本一致。
- 回调验证失败时，优先检查 `APISecret`、回调 URI、请求方法、事件类型注册情况。
