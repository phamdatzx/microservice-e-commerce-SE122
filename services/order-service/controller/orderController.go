package controller

import (
	"order-service/service"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{service: service}
}

// TODO: Implement order controller methods
