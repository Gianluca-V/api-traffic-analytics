package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"api-traffic-analytics/cmd/api-gateway/internal/config"
)

type ProxyService struct {
	client *http.Client
	config *config.Config
}

func NewProxyService(cfg *config.Config) *ProxyService {
	return &ProxyService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: cfg,
	}
}

func (s *ProxyService) ProxyToService(ctx context.Context, service string, method, path string, body []byte) (*http.Response, error) {
	var baseURL string

	switch service {
	case "traffic":
		baseURL = s.config.TrafficIngestorURL
	case "analytics":
		baseURL = s.config.AnalyticsServiceURL
	case "alerts":
		baseURL = s.config.AlertingServiceURL
	default:
		return nil, fmt.Errorf("unknown service: %s", service)
	}

	// Construct full URL
	fullURL, err := url.JoinPath(baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	// Create request
	var req *http.Request
	if body != nil {
		req, err = http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, fullURL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers
	// Note: Don't copy Authorization header for security
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "API-Gateway")

	// Make request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to proxy request: %w", err)
	}

	return resp, nil
}

func (s *ProxyService) GetAnalytics(ctx context.Context, path string) ([]byte, error) {
	resp, err := s.ProxyToService(ctx, "analytics", "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("analytics service returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (s *ProxyService) GetAlerts(ctx context.Context, path string) ([]byte, error) {
	resp, err := s.ProxyToService(ctx, "alerts", "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("alerts service returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (s *ProxyService) GetTrafficData(ctx context.Context, path string) ([]byte, error) {
	resp, err := s.ProxyToService(ctx, "traffic", "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("traffic service returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
