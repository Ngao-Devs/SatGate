package webhooks

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	response, err := h.service.Process(c.Request.Context(), c.Request.Header, body)
	if err == nil {
		c.JSON(http.StatusOK, response)
		return
	}

	var processingError *ProcessingError
	if !errors.As(err, &processingError) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch processingError.Kind {
	case ErrSignature:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "signature verification failed"})
	case ErrInvalidBody:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook body"})
	case ErrInvalidStatus:
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhook status is not paid"})
	case ErrPaymentHash:
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment hash does not match invoice"})
	case ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case ErrInternal:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
