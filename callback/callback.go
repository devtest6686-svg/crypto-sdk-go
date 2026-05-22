package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

// HandlerFunc 泛型回调处理函数类型
type HandlerFunc[T any] func(ctx context.Context, data T) error

// handlerEntry 用于存储注册的回调及其类型信息
type handlerEntry struct {
	handler any
	typeOf  reflect.Type
}

// CallbackManager 回调管理器
type CallbackManager struct {
	callbackUri string
	client      *client.CryptoClient
	handlers    map[entity.EventType]handlerEntry
}

// NewCallbackManager 创建回调管理器
// callbackUri 可选参数，指定回调的URI路径，默认为请求的URI路径。
func NewCallbackManager(client *client.CryptoClient, callbackUri string) *CallbackManager {
	return &CallbackManager{
		callbackUri: callbackUri,
		client:      client,
		handlers:    make(map[entity.EventType]handlerEntry),
	}
}

// Register 注册类型安全的回调处理器（包级泛型函数）
func Register[T any](m *CallbackManager, eventType entity.EventType, handler HandlerFunc[T]) {
	m.handlers[eventType] = handlerEntry{
		handler: handler,
		typeOf:  reflect.TypeOf((*T)(nil)).Elem(),
	}
}

// HandleRequest 处理HTTP请求
func (m *CallbackManager) HandleRequest(ctx context.Context, req *http.Request) error {

	m.client.Logger().Infof(ctx, "HandleCallback Header trace_id:%s job_id:%s", req.Header.Get("trace-id"), req.Header.Get("job-id"))

	ctx = context.WithValue(ctx, m.client.TraceIdField(), req.Header.Get("trace-id"))

	body, err := m.client.HandleCallback(req, m.callbackUri)
	if err != nil {
		m.client.Logger().Errorf(ctx, "CallbackManager HandleRequest err:%v", err)
		return fmt.Errorf("回调签名验证失败: %v", err)
	}
	var callbackEvent entity.CallbackEvent
	if err := json.Unmarshal(body, &callbackEvent); err != nil {
		m.client.Logger().Errorf(ctx, "CallbackManager CallbackEvent err:%v", err)
		return fmt.Errorf("解析回调事件失败: %v", err)
	}
	handlerEntry, exists := m.handlers[callbackEvent.Event]
	if !exists {
		m.client.Logger().Errorf(ctx, "CallbackManager Event err:%v", callbackEvent.Event)
		return fmt.Errorf("未注册的事件类型: %s", callbackEvent.Event)
	}
	// 反序列化为注册时指定的类型
	value := reflect.New(handlerEntry.typeOf).Interface()
	if err := json.Unmarshal(body, value); err != nil {
		m.client.Logger().Errorf(ctx, "CallbackManager typeOf err:%v", err)
		return fmt.Errorf("事件数据反序列化失败: %v", err)
	}

	// 调用处理器
	result := reflect.ValueOf(handlerEntry.handler).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(reflect.Indirect(reflect.ValueOf(value)).Interface()),
	})
	if !result[0].IsNil() {
		return result[0].Interface().(error)
	}
	return nil
}

// DefaultHandlers 默认处理器
var DefaultHandlers = struct {
	IncomingConfirmed HandlerFunc[entity.Event]
	OutgoingConfirmed HandlerFunc[entity.Event]
	OutgoingFailed    HandlerFunc[entity.Event]
	KeepaliveEvent    HandlerFunc[entity.KeepaliveEvent]
}{
	IncomingConfirmed: func(ctx context.Context, data entity.Event) error {
		return fmt.Errorf("转入成功事件处理失败: %v", data)
	},
	OutgoingConfirmed: func(ctx context.Context, data entity.Event) error {
		return fmt.Errorf("转出成功事件处理失败: %v", data)
	},
	OutgoingFailed: func(ctx context.Context, data entity.Event) error {
		return fmt.Errorf("转出失败事件处理失败: %v", data)
	},
	KeepaliveEvent: func(ctx context.Context, data entity.KeepaliveEvent) error {
		return fmt.Errorf("keepalive事件处理失败: %v", data)
	},
}
