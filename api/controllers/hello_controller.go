package controllers

import (
	"net/http"

	"diddo-api/services"
	"diddo-api/utils"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	helloService := services.NewHelloService()
	
	response := helloService.GetHelloMessage()
	
	utils.SendJSONResponse(w, http.StatusOK, response)
}