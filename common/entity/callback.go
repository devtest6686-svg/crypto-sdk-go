package entity

import (
	"context"
	"encoding/json"
)

// CallbackEvent 回调事件基础结构
type CallbackEvent struct {
	Event EventType `json:"event"`
}

// CallbackHandler 回调处理器接口
type CallbackHandler interface {
	Handle(ctx context.Context, data json.RawMessage) (any, error)
	EventType() EventType
}

// CallbackResponse 回调响应
type CallbackResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// KeepaliveEvent keepalive事件数据
type KeepaliveEvent struct {
	Message string `json:"message"`
}
