package forms

import (
	"errors"
	"net/http"

	"satgate/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	ownerID, ok := middleware.OwnerID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	forms, err := h.service.List(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forms)
}

func (h *Handler) Create(c *gin.Context) {
	ownerID, ok := middleware.OwnerID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request CreateFormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := h.service.Create(c.Request.Context(), ownerID, request)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, form)
}

func (h *Handler) Update(c *gin.Context) {
	ownerID, ok := middleware.OwnerID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request UpdateFormRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := h.service.Update(
		c.Request.Context(),
		ownerID,
		c.Param("id"),
		request,
	)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, form)
}

func (h *Handler) Delete(c *gin.Context) {
	ownerID, ok := middleware.OwnerID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.service.Delete(c.Request.Context(), ownerID, c.Param("id"))
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) respondWithError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidName):
		c.JSON(http.StatusBadRequest, gin.H{"error": "form name is required"})
	case errors.Is(err, ErrInvalidDomain):
		c.JSON(http.StatusBadRequest, gin.H{"error": "allowed domain is required"})
	case errors.Is(err, ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount_sats must be between 1 and 10000"})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
