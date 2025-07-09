package models

type ClothingIdentificationRequest struct {
	Image []byte `json:"image"`
}

type ClothingItem struct {
	Title           string           `json:"title"`
	SearchValue     string           `json:"search_value"`
	ShoppingResults []ShoppingResult `json:"shopping_results"`
}

type ShoppingResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image"`
	Price string `json:"price"`
}

type ProductSearchResponse struct {
	Data struct {
		Products []struct {
			Title string `json:"product_title"`
			URL   string `json:"product_offer_page_url"`
			Image string `json:"product_photo"`
			Price string `json:"price"`
		} `json:"products"`
	} `json:"data"`
}

type ClothingIdentificationResponse struct {
	Items  []ClothingItem `json:"items"`
	Status string         `json:"status"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}
