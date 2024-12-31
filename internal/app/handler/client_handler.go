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

type ClientHandler struct {
	BaseHandler
	service service.ClientService
}

func NewClientHandler(service service.ClientService, logger *zap.SugaredLogger) *ClientHandler {
	return &ClientHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
	}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	// Get application ID from context (set by auth middleware)
	appID, exists := c.Get("application_id")
	if !exists {
		h.error(c, http.StatusUnauthorized, errors.New("application ID not found in context"))
		return
	}
	client.ApplicationID = appID.(uuid.UUID)

	createdClient, err := h.service.Create(c.Request.Context(), &client)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			h.error(c, http.StatusConflict, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.created(c, createdClient)
}

func (h *ClientHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid client ID"))
		return
	}

	// Get application ID from context
	appID, exists := c.Get("application_id")
	if !exists {
		h.error(c, http.StatusUnauthorized, errors.New("application ID not found in context"))
		return
	}

	client, err := h.service.GetByID(c.Request.Context(), appID.(uuid.UUID), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, client)
}

func (h *ClientHandler) List(c *gin.Context) {
	// Get application ID from context
	appID, exists := c.Get("application_id")
	if !exists {
		h.error(c, http.StatusUnauthorized, errors.New("application ID not found in context"))
		return
	}

	var filters service.ClientFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}
	filters.ApplicationID = appID.(uuid.UUID)

	clients, err := h.service.List(c.Request.Context(), filters)
	if err != nil {
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, clients)
}

func (h *ClientHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid client ID"))
		return
	}

	// Get application ID from context
	appID, exists := c.Get("application_id")
	if !exists {
		h.error(c, http.StatusUnauthorized, errors.New("application ID not found in context"))
		return
	}

	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	client.ID = id
	client.ApplicationID = appID.(uuid.UUID)

	updatedClient, err := h.service.Update(c.Request.Context(), &client)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			h.error(c, http.StatusNotFound, err)
		case errors.Is(err, service.ErrDuplicateEmail):
			h.error(c, http.StatusConflict, err)
		default:
			h.error(c, http.StatusInternalServerError, err)
		}
		return
	}

	h.success(c, updatedClient)
}

func (h *ClientHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid client ID"))
		return
	}

	// Get application ID from context
	appID, exists := c.Get("application_id")
	if !exists {
		h.error(c, http.StatusUnauthorized, errors.New("application ID not found in context"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), appID.(uuid.UUID), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, service.ErrClientHasActiveLicenses) {
			h.error(c, http.StatusConflict, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.noContent(c)
}

func (h *ClientHandler) GetLicenses(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid client ID"))
		return
	}

	// Get application ID from context
	appID, exists := c.Get("application_id")
	if !exists {
		h.error(c, http.StatusUnauthorized, errors.New("application ID not found in context"))
		return
	}

	licenses, err := h.service.GetClientLicenses(c.Request.Context(), appID.(uuid.UUID), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, licenses)
}
