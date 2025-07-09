package controllers

import (
	"io"
	"net/http"

	"diddo-api/configs"
	"diddo-api/models"
	"diddo-api/services"
	"diddo-api/utils"
)

func IdentifyClothingItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close()

	if err := utils.ValidateImageFile(header); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsValidImageExtension(header.Filename) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid image file extension")
		return
	}

	imageData, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Unable to read image file")
		return
	}

	config := configs.LoadConfig()
	openAIService := services.NewOpenAIService(config)
	productSearchService := services.NewProductSearchService(config)

	clothingItems, err := openAIService.IdentifyClothingItems(imageData)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to identify clothing items: "+err.Error())
		return
	}

	for i, item := range clothingItems {
		productSearchResponse, err := productSearchService.SearchProducts(item.SearchValue)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to search for products: "+err.Error())
			return
		}

		var shoppingResults []models.ShoppingResult
		for _, product := range productSearchResponse.Data.Products {
			shoppingResults = append(shoppingResults, models.ShoppingResult{
				Title: product.Title,
				URL:   product.URL,
				Image: product.Image,
				Price: product.Price,
			})
		}
		clothingItems[i].ShoppingResults = shoppingResults
	}

	response := models.ClothingIdentificationResponse{
		Items:  clothingItems,
		Status: "success",
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}
