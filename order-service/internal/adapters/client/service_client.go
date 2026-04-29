package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ServiceClient struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *ServiceClient {
	return &ServiceClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ServiceClient) CheckUser(ctx context.Context, userID, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/users/"+userID, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("user check failed: %s", resp.Status)
	}
	return nil
}

func (c *ServiceClient) CheckProduct(ctx context.Context, productID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/product/"+productID, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("product check failed: %s", resp.Status)
	}
	var body map[string]any
	return json.NewDecoder(resp.Body).Decode(&body)
}
