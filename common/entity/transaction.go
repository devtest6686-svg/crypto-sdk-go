package entity

import "github.com/shopspring/decimal"

type TransactionDetailReq struct {
	WalletId        string         `json:"wallet_id" binding:"required" example:"69fabb9450a30bdc3c1a79a9"`                                                // 钱包id
	Blockchain      BlockchainName `json:"blockchain" binding:"required" example:"tron"`                                                                   // 公链名称, 如：ethereum, solana, bnb-smart-chain, tron
	TransactionHash string         `json:"transaction_hash" binding:"required" example:"8d113e3a120134e8fb4eed24e35ea97023f43af7efd8ea435676f7433cf0b21c"` // 交易哈希
}

type TransactionDetailResp struct {
	RequestId string                    `json:"request_id"`      // 定义请求的 ID，requestId 由crypto服务生成，每个请求的 requestId 都是唯一的
	Error     *APIError                 `json:"error,omitempty"` // 可选，定义请求失败时的错误信息
	Data      TransactionDetailRespData `json:"data,omitempty"`  //
}

type TransactionDetailRespData struct {
	Network        NetworkName       `json:"network"`         // 网络名称，如：mainnet, testnet, sepolia, nile
	BlockHeight    int64             `json:"block_height"`    // 区块高度
	BlockTimestamp int64             `json:"block_timestamp"` // 区块时间戳（秒）
	Status         TransactionStatus `json:"status"`          // 交易状态，如：success, failed
	Direction      Direction         `json:"direction"`       // 交易方向，如：incoming, outgoing, internal
	Senders        []Sender          `json:"senders"`         // 发送方信息
	Recipients     []Recipient       `json:"recipients"`      // 接收方信息
	TokenTransfers []TokenTransfer   `json:"token_transfers"` // 代币发送信息，只对代币交易有效；同一笔 tx 可能含多笔代币转账；无代币转账时不返回该字段
	Fee            TransactionFee    `json:"fee"`             // 交易Gas费用
}

type Direction string

const (
	DirectionIncoming Direction = "incoming"
	DirectionOutgoing Direction = "outgoing"
	DirectionInternal Direction = "internal"
)

type TransactionStatus string

const (
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type TransactionFee struct {
	Amount decimal.Decimal `json:"amount"` // Gas费用
	Symbol string          `json:"symbol"` // Gas符号，如：TRX, ETH, BNB, SOL.
}

type Sender struct {
	Address         string          `json:"address"`           // 发送方地址
	AddressName     string          `json:"address_name"`      // 地址名称
	IsWalletAddress bool            `json:"is_wallet_address"` // 是否为钱包地址，true-主钱包派生的地址
	IsMasterAddress bool            `json:"is_master_address"` // 是否为钱包master地址，true-主钱包master地址（设置为Gas钱包的master地址，归集时可用于支付Gas手续费）
	Amount          decimal.Decimal `json:"amount"`            // 发送原生币数量
	Symbol          string          `json:"symbol"`            // 发送原生币符号：TRX, ETH, BNB, SOL.
}

type Recipient struct {
	Address         string          `json:"address"`           // 接收方地址
	AddressName     string          `json:"address_name"`      // 地址名称
	IsWalletAddress bool            `json:"is_wallet_address"` // 是否为钱包地址，true-主钱包派生的地址
	IsMasterAddress bool            `json:"is_master_address"` // 是否为钱包master地址，true-主钱包master地址（设置为Gas钱包的master地址，归集时可用于支付Gas手续费）
	Amount          decimal.Decimal `json:"amount"`            // 接收原生币数量
	Symbol          string          `json:"symbol"`            // 接收原生币符号，如：TRX, ETH, BNB, SOL
}

type TokenTransfer struct {
	FromAddress string          `json:"from_address"`       // 发送方地址
	ToAddress   string          `json:"to_address"`         // 接收方地址
	Amount      decimal.Decimal `json:"amount"`             // 代币数量
	Contract    string          `json:"contract,omitempty"` // 代币合约地址
	Decimals    int32           `json:"decimals"`           // 代币精度
	Symbol      string          `json:"symbol"`             // 代币符号，如：USDT, USDC
	TokenName   string          `json:"token_name"`         // 代币名称： Tether, USD Coin.
}
