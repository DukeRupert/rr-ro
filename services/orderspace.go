// services/orderspace.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"orderspace-integration/models"
	"time"
)

type OrderSpaceService struct {
	clientID     string
	clientSecret string
	accessToken  string
	tokenExpiry  time.Time
	baseURL      string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewOrderSpaceService(clientID, clientSecret string) *OrderSpaceService {
	return &OrderSpaceService{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      "https://api.orderspace.com/v1",
	}
}

// GetValidToken returns a valid access token, refreshing if necessary
func (s *OrderSpaceService) GetValidToken() (string, error) {
	if s.accessToken == "" || time.Now().After(s.tokenExpiry) {
		if err := s.refreshToken(); err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}
	}
	return s.accessToken, nil
}

// refreshToken gets a new access token using client credentials
func (s *OrderSpaceService) refreshToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)

	req, err := http.NewRequest("POST", "https://identity.orderspace.com/oauth/token",
		bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	s.accessToken = tokenResp.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// GetProducts fetches all products from the Orderspace API
func (s *OrderSpaceService) GetProducts() ([]models.OrderspaceProduct, error) {
	token, err := s.GetValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products", s.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResponse models.OrderspaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return apiResponse.Products, nil
}

// GetProduct fetches a single product by ID
func (s *OrderSpaceService) GetProduct(id string) (*models.OrderspaceProduct, error) {
	token, err := s.GetValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products/%s", s.baseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var product models.OrderspaceProduct
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product: %w", err)
	}

	return &product, nil
}
