package callback

import (
	"encoding/json"
	"net/http"
)

// HTTPHandler HTTP处理器
type HTTPHandler struct {
	manager *CallbackManager
}

// NewHTTPHandler 创建HTTP处理器
func NewHTTPHandler(manager *CallbackManager) *HTTPHandler {
	return &HTTPHandler{
		manager: manager,
	}
}

// ServeHTTP 实现http.Handler接口
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 处理回调请求
	err := h.manager.HandleRequest(ctx, r)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// 返回成功响应
	h.writeSuccess(w, nil)
}

// writeError 写入错误响应
func (h *HTTPHandler) writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	errorResp := map[string]interface{}{
		"code": 400,
		"msg":  err.Error(),
	}

	json.NewEncoder(w).Encode(errorResp)
}

// writeSuccess 写入成功响应
func (h *HTTPHandler) writeSuccess(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	successResp := map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": result,
	}

	json.NewEncoder(w).Encode(successResp)
}

// Middleware 创建HTTP中间件
func (h *HTTPHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果是回调请求，由回调处理器处理
		if r.Method == http.MethodPost && r.URL.Path == "/callback" {
			h.ServeHTTP(w, r)
			return
		}

		// 否则传递给下一个处理器
		next.ServeHTTP(w, r)
	})
}

// HandlerFunc 返回HTTP处理函数
func (h *HTTPHandler) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
