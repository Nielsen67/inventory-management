package models

import "time"

type Order struct {
	ID              uint          `gorm:"primarykey;autoIncrement"`
	Date            time.Time     `json:"date"`
	CustomerName    string        `json:"customerName"`
	ShippingAddress string        `json:"shippingAddress"`
	Total           uint          `json:"total"`
	OrderDetails    []OrderDetail `gorm:"foreignKey:OrderId"`
}

type OrderDetail struct {
	OrderId     uint `gorm:"primaryKey"`
	ProductId   uint `gorm:"primaryKey"`
	Qty         uint `json:"qty"`
	PriceDetail uint `json:"priceDetail"`
}

type CreateOrderRequest struct {
	CustomerName    string               `json:"customerName"`
	ShippingAddress string               `json:"shippingAddress"`
	OrderDetails    []OrderDetailRequest `json:"orderDetails" binding:"required,dive"`
}

type OrderDetailRequest struct {
	ProductId uint `json:"productId" binding:"required"`
	Qty       uint `json:"qty" binding:"required"`
}

type OrderResponse struct {
	ID                  uint                  `json:"id"`
	Date                time.Time             `json:"date"`
	CustomerName        string                `json:"customerName"`
	ShippingAddress     string                `json:"shippingAddress"`
	Total               uint                  `json:"total"`
	OrderDetailResponse []OrderDetailResponse `json:"orderDetails"`
}

type OrderDetailResponse struct {
	ProductId   uint `json:"productId"`
	Qty         uint `json:"qty"`
	PriceDetail uint `json:"priceDetail"`
}
