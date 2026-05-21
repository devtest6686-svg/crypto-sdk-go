package entity

import "github.com/shopspring/decimal"

type WalletListReq struct {
	PageParamReq
}

type WalletListResp struct {
	RequestId string             `json:"requestId"`       // 定义请求的 ID，requestId 由crypto服务生成，每个请求的 requestId 都是唯一的
	Error     *APIError          `json:"error,omitempty"` // 可选，定义请求失败时的错误信息
	Data      WalletListRespData `json:"data,omitempty"`  //
}

type WalletListRespData struct {
	PageParamResp
	Items []WalletItem `json:"items"`
}

type WalletAddressGenerateReq struct {
	WalletId    string         `json:"walletId" binding:"required"`    // 钱包id
	Blockchain  BlockchainName `json:"blockchain"  binding:"required"` // 区块链名称, 如：Ethereum, Solana, Binance Smart Chain, Tron
	Network     NetworkName    `json:"network"  binding:"required"`    // 网络名称, 如：mainnet, testnet, sepolia, nile
	AddressName string         `json:"addressName" binding:"required"` // 地址名称，如："用户充值"
}

type WalletAddressGenerateResp struct {
	RequestId string                        `json:"requestId"`       // 定义请求的 ID，requestId 由crypto服务生成，每个请求的 requestId 都是唯一的
	Error     *APIError                     `json:"error,omitempty"` // 可选，定义请求失败时的错误信息
	Data      WalletAddressGenerateRespData `json:"data,omitempty"`  //
}

type WalletAddressGenerateRespData struct {
	Id               string `json:"id"`               // 地址id
	Address          string `json:"address"`          // 地址
	AddressName      string `json:"addressName"`      // 地址名称
	CreatedTimestamp int64  `json:"createdTimestamp"` // 创建时间戳（毫秒）
}

type WalletAddressListReq struct {
	PageParamReq
	WalletId   string         `json:"walletId" form:"walletId" binding:"required"`     // 钱包id
	Blockchain BlockchainName `json:"blockchain" form:"blockchain" binding:"required"` // 区块链名称, 如：Ethereum, Solana, Binance Smart Chain, Tron
}

type WalletAddressListResp struct {
	RequestId string                    `json:"requestId"`       // 定义请求的 ID，requestId 由crypto服务生成，每个请求的 requestId 都是唯一的
	Error     *APIError                 `json:"error,omitempty"` // 可选，定义请求失败时的错误信息
	Data      WalletAddressListRespData `json:"data,omitempty"`  //
}

type WalletAddressListRespData struct {
	PageParamResp
	Items []AddressItem `json:"items"`
}

type WalletAddressAssetsReq struct {
	WalletId   string         `json:"walletId" form:"walletId" binding:"required"` // 钱包id
	Address    string         `json:"address" form:"address" binding:"required"`   // 地址
	Blockchain BlockchainName `json:"blockchain" form:"blockchain"`                // 可选，区块链名称, 如：Ethereum, Solana, Binance Smart Chain, Tron
}

type WalletAddressAssetsResp struct {
	RequestId string                      `json:"requestId"`       // 定义请求的 ID，requestId 由crypto服务生成，每个请求的 requestId 都是唯一的
	Error     *APIError                   `json:"error,omitempty"` // 可选，定义请求失败时的错误信息
	Data      WalletAddressAssetsRespData `json:"data,omitempty"`  //
}

type WalletAddressAssetsRespData struct {
	Items []AssetItem `json:"items"`
}

type PageParamReq struct {
	Limit   int    `json:"limit" form:"limit"`     // 可选，返回记录数限制，默认: 50
	StartId string `json:"startId" form:"startId"` // 可选，定义当前列表应从上一个已列出的记录的 ID 开始
}
type PageParamResp struct {
	Limit   int    `json:"limit"`   // 返回记录数限制
	StartId string `json:"startId"` // 下一页应传的游标 ID（本页最后一条记录的 ID），无更多数据时为空
	HasMore bool   `json:"hasMore"` // 是否有更多记录
}

type WalletItem struct {
	Id               string        `json:"id"`               // 钱包id
	Name             string        `json:"name"`             // 钱包名称
	CreatedTimestamp int64         `json:"createdTimestamp"` // 钱包创建时间戳（毫秒）
	IsMaster         bool          `json:"isMaster"`         // 是否为主钱包
	Gas              WalletGasType `json:"gas"`              // 是否为Gas钱包(0-未设置为Gas钱包, 1-已设置为Gas钱包, 2-已设置为默认Gas钱包)
}

// WalletGasType Gas钱包类型枚举 0-未设置为Gas钱包, 1-已设置为Gas钱包, 2-已设置为默认Gas钱包
type WalletGasType int

const (
	WalletNormal     WalletGasType = iota // 0-未设置为Gas钱包
	WalletGas                             // 1-已设置为Gas钱包
	WalletGasDefault                      // 2-已设置为默认Gas钱包
)

type BlockchainName string

const (
	BlockchainEthereum BlockchainName = "ethereum"
	BlockchainBSC      BlockchainName = "bnb-smart-chain"
	BLockchainTron     BlockchainName = "tron"
	BlockchainSolana   BlockchainName = "solana"
)

type NetworkName string

const (
	NetworkMainnet NetworkName = "mainnet"
	NetworkTestnet NetworkName = "testnet"
	NetworkSepolia NetworkName = "sepolia"
	NetworkNile    NetworkName = "nile"
	NewworkDevnet  NetworkName = "devnet"
)

type AddressItem struct {
	Id               string `json:"id"`               // 地址id
	Address          string `json:"address"`          // 地址
	AddressName      string `json:"addressName"`      // 地址名称
	CreatedTimestamp int64  `json:"createdTimestamp"` // 创建时间戳（毫秒）
}

type AssetType string

const (
	AssetToken      AssetType = "token"
	AssetNativeCoin AssetType = "coin"
)

type AssetItem struct {
	AssetId         string          `json:"assetId"`            // 资产id
	Blockchain      BlockchainName  `json:"blockchain"`         // 区块链名称, 如：Ethereum, Solana, Binance Smart Chain, Tron
	Network         NetworkName     `json:"network"`            // 网络名称, 如：mainnet, testnet, sepolia, nile
	AssetType       AssetType       `json:"assetType"`          // 资产类型，原生币（coin）或代币（token）
	Contract        string          `json:"contract,omitempty"` // 代币合约地址，仅对代币（token）有效
	AvailableAmount decimal.Decimal `json:"availableAmount"`    // 可用余额
	AllocatedAmount decimal.Decimal `json:"allocatedAmount"`    // 待处理余额
}
