package services

import (
	"diddo-api/configs"
	"diddo-api/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ProductSearchService struct {
	config *configs.Config
}

func NewProductSearchService(config *configs.Config) *ProductSearchService {
	return &ProductSearchService{
		config: config,
	}
}

func (s *ProductSearchService) SearchProducts(query string) (models.ProductSearchResponse, error) {
	if s.config.ProductSearchAPIKey == "" {
		return models.ProductSearchResponse{}, fmt.Errorf("product search API key not configured")
	}

	apiURL := "https://real-time-product-search.p.rapidapi.com/search-light-v2"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return models.ProductSearchResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("country", "us")
	q.Add("language", "en")
	q.Add("limit", "3")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("x-rapidapi-host", "real-time-product-search.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", s.config.ProductSearchAPIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.ProductSearchResponse{}, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.ProductSearchResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("Response Body:", string(body))

	var productSearchResponse models.ProductSearchResponse
	if err := json.Unmarshal(body, &productSearchResponse); err != nil {
		return models.ProductSearchResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return productSearchResponse, nil
}
