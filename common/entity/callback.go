package entity

import (
	"context"
	"encoding/json"

	"github.com/lbe-io/crypto-sdk-go/common/consts"
)

// CallbackEvent 回调事件基础结构
type CallbackEvent struct {
	EventType consts.CallbackEventType `json:"event_type"`
	Data      json.RawMessage          `json:"data"`
}

// CallbackHandler 回调处理器接口
type CallbackHandler interface {
	Handle(ctx context.Context, data json.RawMessage) (any, error)
	EventType() consts.CallbackEventType
}

// CallbackResponse 回调响应
type CallbackResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// TransactionCreatedEvent 交易创建事件数据
type TransactionCreatedEvent struct {
	TransactionID string `json:"transaction_id"`
	FromAddress   string `json:"from_address"`
	ToAddress     string `json:"to_address"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// TransactionUpdatedEvent 交易更新事件数据
type TransactionUpdatedEvent struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	UpdatedAt     string `json:"updated_at"`
	BlockHeight   int64  `json:"block_height,omitempty"`
	TxHash        string `json:"tx_hash,omitempty"`
}

// WalletCreatedEvent 钱包创建事件数据
type WalletCreatedEvent struct {
	WalletID   string `json:"wallet_id"`
	WalletName string `json:"wallet_name"`
	Currency   string `json:"currency"`
	CreatedAt  string `json:"created_at"`
}

// WalletUpdatedEvent 钱包更新事件数据
type WalletUpdatedEvent struct {
	WalletID   string `json:"wallet_id"`
	WalletName string `json:"wallet_name,omitempty"`
	Status     string `json:"status,omitempty"`
	UpdatedAt  string `json:"updated_at"`
}

// AddressGeneratedEvent 地址生成事件数据
type AddressGeneratedEvent struct {
	WalletID  string `json:"wallet_id"`
	Address   string `json:"address"`
	Currency  string `json:"currency"`
	CreatedAt string `json:"created_at"`
}

// TestEvent 测试事件数据
type TestEvent struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}
