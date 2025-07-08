package models

type ClothingIdentificationRequest struct {
	Image []byte `json:"image"`
}

type ClothingItemDescription struct {
	Description string `json:"description,omitempty"`
}

type ClothingIdentificationResponse struct {
	Items  []ClothingItemDescription `json:"items"`
	Status string                    `json:"status"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}
