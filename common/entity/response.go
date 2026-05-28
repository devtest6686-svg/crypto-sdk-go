package entity

import "fmt"

type CommonResponse struct {
	RequestId string      `json:"request_id"`      // 定义请求的 ID，requestId 由crypto服务生成，每个请求的 requestId 都是唯一的
	Error     *APIError   `json:"error,omitempty"` // 可选，定义请求失败时的错误信息
	Data      interface{} `json:"data,omitempty"`  //
}

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *APIError) Detail() error {
	return fmt.Errorf("code:%v, message:%v, details:%+v", e.Code, e.Message, e.Details)
}
