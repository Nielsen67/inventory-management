package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name         string        `json:"name" gorm:"unique"`
	Description  string        `json:"description"`
	Price        uint          `json:"price"`
	Category     string        `json:"category"`
	Inventories  []Inventory   `gorm:"foreignKey:ProductId"`
	OrderDetails []OrderDetail `gorm:"foreignKey:ProductId"`
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"content"`
	Price       uint   `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" binding:"omiempty"`
	Description string `json:"description" binding:"omiempty"`
	Price       uint   `json:"price" binding:"omiempty"`
	Category    string `json:"category" binding:"omiempty"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
}
