package handler

import (
	"net/http"

	"api-traffic-analytics/cmd/traffic-ingestor/internal/service"
	"api-traffic-analytics/internal/shared/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ReceiveTrafficData(c *gin.Context) {
	var data models.TrafficData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ProcessTrafficData(c.Request.Context(), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}