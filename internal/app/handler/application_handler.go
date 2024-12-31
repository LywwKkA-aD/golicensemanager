package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/LywwKkA-aD/golicensemanager/internal/models"
	"github.com/LywwKkA-aD/golicensemanager/internal/service"
)

type ApplicationHandler struct {
	BaseHandler
	service service.ApplicationService
}

func NewApplicationHandler(service service.ApplicationService, logger *zap.SugaredLogger) *ApplicationHandler {
	return &ApplicationHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
	}
}

func (h *ApplicationHandler) Create(c *gin.Context) {
	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	createdApp, err := h.service.Create(c.Request.Context(), &app)
	if err != nil {
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.created(c, createdApp)
}

func (h *ApplicationHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid application ID"))
		return
	}

	app, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, app)
}

func (h *ApplicationHandler) List(c *gin.Context) {
	apps, err := h.service.List(c.Request.Context())
	if err != nil {
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, apps)
}

func (h *ApplicationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid application ID"))
		return
	}

	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}
	app.ID = id

	updatedApp, err := h.service.Update(c.Request.Context(), &app)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, updatedApp)
}

func (h *ApplicationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid application ID"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.noContent(c)
}

func (h *ApplicationHandler) GenerateToken(c *gin.Context) {
	var req struct {
		APIKey    string `json:"api_key" binding:"required"`
		APISecret string `json:"api_secret" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.GenerateToken(c.Request.Context(), req.APIKey, req.APISecret)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			h.error(c, http.StatusUnauthorized, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, gin.H{"token": token})
}
