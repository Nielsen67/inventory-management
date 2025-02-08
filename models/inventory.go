package models

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	Stock     uint   `json:"stock"`
	Location  string `json:"location"`
	ProductId uint   `json:"productId"`
}

type CreateInventoryRequest struct {
	Stock     uint   `json:"stock" binding:"required"`
	Location  string `json:"location" binding:"required"`
	ProductId uint   `json:"productId" binding:"required"`
}

type UpdateInventoryRequest struct {
	Stock uint `json:"stock" binding:"required"`
}

type InventoryResponse struct {
	Stock     uint   `json:"stock"`
	Location  string `json:"location"`
	ProductId uint   `json:"productId"`
}
