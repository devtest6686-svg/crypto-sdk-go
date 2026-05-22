package entity

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

// EventType 回调事件类型：
// 原生币/代币转入成功：INCOMING_CONFIRMED
// 原生币/代币转出成功：OUTGOING_CONFIRMED
// 原生币/代币转出失败：OUTGOING_FAILED
type EventType string

const (
	EventIncomingConfirmed EventType = "INCOMING_CONFIRMED"
	EventOutgoingConfirmed EventType = "OUTGOING_CONFIRMED"
	EventOutgoingFailed    EventType = "OUTGOING_FAILED"
	EventKeepalive         EventType = "KEEPALIVE"
)

type Event struct {
	WalletID string    `json:"wallet_id"` // 钱包id
	Event    EventType `json:"event"`     // 事件类型 INCOMING_CONFIRMED, OUTGOING_CONFIRMED, OUTGOING_FAILED
	Item     EventItem `json:"item"`      // 事件
}

// EventItem
// 当一笔交易的转入转出地址均为主钱包相关地址时，需要发送2次事件：转入地址的incoming事件, 转出地址的outgoing事件
type EventItem struct {
	Direction       Direction         `json:"direction"`        // 交易方向，只支持incoming, outgoing
	Address         string            `json:"address"`          // 产生回调事件的地址
	Blockchain      BlockchainName    `json:"blockchain"`       // 区块链名称, 如：Ethereum, Solana, Binance Smart Chain, Tron
	Network         NetworkName       `json:"network"`          // 网络名称, 如：mainnet, testnet, sepolia, nile
	TransactionHash string            `json:"transaction_hash"` // 交易哈希
	BlockHeight     int64             `json:"block_height"`     // 区块高度
	BlockTimestamp  int64             `json:"block_timestamp"`  // 区块时间戳（秒）
	Status          TransactionStatus `json:"status"`           // 交易状态, 如：success, failed
	Fee             TransactionFee    `json:"fee"`              // 交易Gas费用
	TokenTransfers  TokenTransfer     `json:"token_transfers"`  // 代币转账数量
	Amount          decimal.Decimal   `json:"amount"`           // 原生币转账数量
}

func (e *Event) Json() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}
