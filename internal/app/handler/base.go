package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type BaseHandler struct {
	logger *zap.SugaredLogger
}

func NewBaseHandler(logger *zap.SugaredLogger) BaseHandler {
	return BaseHandler{
		logger: logger,
	}
}

func (h *BaseHandler) success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func (h *BaseHandler) error(c *gin.Context, status int, err error) {
	c.JSON(status, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func (h *BaseHandler) created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func (h *BaseHandler) noContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
