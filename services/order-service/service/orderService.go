package service

import (
	"order-service/repository"
)

type OrderService interface {
	// TODO: Add order service methods
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{repo: orderRepo}
}

// TODO: Implement order service methods
