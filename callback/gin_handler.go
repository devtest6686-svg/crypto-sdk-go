package callback

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
	manager *CallbackManager
}

func NewGinHandler(manager *CallbackManager) *GinHandler {
	return &GinHandler{
		manager: manager,
	}
}

func (h *GinHandler) HandleFunc(c *gin.Context) {
	ctx := c.Request.Context()

	// 处理回调请求
	err := h.manager.HandleRequest(ctx, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
