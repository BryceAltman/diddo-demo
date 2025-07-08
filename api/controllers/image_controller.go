package controllers

import (
	"io"
	"net/http"

	"diddo-api/configs"
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

	response, err := openAIService.IdentifyClothingItems(imageData)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to identify clothing items: "+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}