package controllers

import (
	"inventory-management/models"
	"inventory-management/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
	}

	tx := pc.DB.Begin()
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	if err := pc.DB.Preload("Product").First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading product data"})
		return
	}

	c.JSON(201, gin.H{"message": "Product is successfully created", "data": product})

}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var req models.UpdateProductRequest
	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	var product models.Product
	// Check if post exists and belongs to user
	if err := pc.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product is not exists"})
		return
	}

	tx := pc.DB.Begin()

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Category != "" {
		product.Category = req.Category
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	if err := pc.DB.Preload("Product").First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading product data"})
		return
	}
	c.JSON(200, gin.H{"message": "Product is successfully updated"})

}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	var product models.Product

	if err := pc.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product is not exists"})
		return
	}

	tx := pc.DB.Begin()

	if err := tx.Model(&product).Association("Inventory").Clear(); err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Unscoped().Delete(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(204, gin.H{})

}

func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product

	response := make([]models.ProductResponse, len(products))
	for i, product := range products {
		response[i] = models.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
		}
	}

	c.JSON(200, gin.H{"data": response})
}
