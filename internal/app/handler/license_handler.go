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

type LicenseHandler struct {
	BaseHandler
	service service.LicenseService
}

func NewLicenseHandler(service service.LicenseService, logger *zap.SugaredLogger) *LicenseHandler {
	return &LicenseHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
	}
}

func (h *LicenseHandler) Create(c *gin.Context) {
	var license models.License
	if err := c.ShouldBindJSON(&license); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	createdLicense, err := h.service.Create(c.Request.Context(), &license)
	if err != nil {
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.created(c, createdLicense)
}

func (h *LicenseHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid license ID"))
		return
	}

	license, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, license)
}

func (h *LicenseHandler) List(c *gin.Context) {
	var filters service.LicenseFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	licenses, err := h.service.List(c.Request.Context(), filters)
	if err != nil {
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, licenses)
}

func (h *LicenseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid license ID"))
		return
	}

	var license models.License
	if err := c.ShouldBindJSON(&license); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}
	license.ID = id

	updatedLicense, err := h.service.Update(c.Request.Context(), &license)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, updatedLicense)
}

func (h *LicenseHandler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.error(c, http.StatusBadRequest, errors.New("invalid license ID"))
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Revoke(c.Request.Context(), id, req.Reason); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.error(c, http.StatusNotFound, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.noContent(c)
}

func (h *LicenseHandler) Validate(c *gin.Context) {
	var req struct {
		LicenseKey string `json:"license_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.error(c, http.StatusBadRequest, err)
		return
	}

	validationResult, err := h.service.Validate(c.Request.Context(), req.LicenseKey)
	if err != nil {
		if errors.Is(err, service.ErrLicenseInvalid) {
			h.error(c, http.StatusUnauthorized, err)
			return
		}
		h.error(c, http.StatusInternalServerError, err)
		return
	}

	h.success(c, validationResult)
}
