package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"

	"api-traffic-analytics/cmd/api-gateway/internal/config"
	"api-traffic-analytics/cmd/api-gateway/internal/service"
	"api-traffic-analytics/internal/shared/models"
)

type Handler struct {
	proxyService *service.ProxyService
	cfg          *config.Config
}

func NewHandler(proxyService *service.ProxyService, cfg *config.Config) *Handler {
	return &Handler{
		proxyService: proxyService,
		cfg:          cfg,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	response := models.HealthCheckResponse{
		Status:    "healthy",
		Timestamp: models.TimeNow(),
		Version:   "1.0.0",
		Services: map[string]string{
			"api-gateway": "healthy",
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	ctx := c.Request.Context()

	path := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	data, err := h.proxyService.GetAnalytics(ctx, path)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Error:   "Internal server error",
			Message: fmt.Sprintf("Failed to get analytics: %v", err),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h *Handler) GetAlerts(c *gin.Context) {
	ctx := c.Request.Context()

	path := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	data, err := h.proxyService.GetAlerts(ctx, path)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Error:   "Internal server error",
			Message: fmt.Sprintf("Failed to get alerts: %v", err),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h *Handler) GetTrafficData(c *gin.Context) {
	ctx := c.Request.Context()

	path := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	data, err := h.proxyService.GetTrafficData(ctx, path)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Error:   "Internal server error",
			Message: fmt.Sprintf("Failed to get traffic data: %v", err),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h *Handler) GetAnalyticsByLocation(c *gin.Context) {
	// Forward to analytics service with location parameter
	h.GetAnalytics(c)
}

func (h *Handler) GetAlertsByLocation(c *gin.Context) {
	// Forward to alerts service with location parameter
	h.GetAlerts(c)
}

func (h *Handler) GetTrafficDataByLocation(c *gin.Context) {
	// Forward to traffic service with location parameter
	h.GetTrafficData(c)
}

func (h *Handler) ReceiveTrafficData(c *gin.Context) {
	// Proxy POST /traffic to traffic-ingestor service
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Error:   "Bad request",
			Message: fmt.Sprintf("Failed to read request body: %v", err),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	resp, err := h.proxyService.ProxyToService(ctx, "traffic", "POST", "/traffic", body)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Error:   "Internal server error",
			Message: fmt.Sprintf("Failed to proxy traffic data: %v", err),
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	defer resp.Body.Close()

	// Copy response from traffic-ingestor
	responseData, _ := io.ReadAll(resp.Body)

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Status(resp.StatusCode)
	c.Writer.Write(responseData)
}

func (h *Handler) ProxyToService(c *gin.Context) {
	// Generic proxy for internal services
	serviceName := c.Query("service")
	if serviceName == "" {
		errorResponse := models.ErrorResponse{
			Error:   "Bad request",
			Message: "Service parameter is required",
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get the target URL from query or config
	var targetURL string
	switch serviceName {
	case "traffic":
		targetURL = h.cfg.TrafficIngestorURL
	case "analytics":
		targetURL = h.cfg.AnalyticsServiceURL
	case "alerts":
		targetURL = h.cfg.AlertingServiceURL
	default:
		errorResponse := models.ErrorResponse{
			Error:   "Bad request",
			Message: "Invalid service name",
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		errorResponse := models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Invalid target URL",
		}
		c.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify request to forward to target service
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		// Remove service parameter from path
		req.URL.Path = c.Param("path")
		if c.Request.URL.RawQuery != "" {
			req.URL.RawQuery = c.Request.URL.RawQuery
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
