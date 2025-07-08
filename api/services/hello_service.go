package services

import "diddo-api/models"

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) GetHelloMessage() models.HelloResponse {
	return models.HelloResponse{
		Message: "Hello diddo!",
		Status:  "success",
	}
}