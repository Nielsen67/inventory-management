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

	var existingProduct models.Product
	if err := pc.DB.Where("name = ?", req.Name).First(&existingProduct).Error; err == nil {
		c.JSON(400, gin.H{"error": "Product name is already used"})
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

	productResp := models.ProductDataResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	resp := models.ProductResponse{
		Message: "Product is successfully created",
		Data:    productResp,
	}

	if err := pc.DB.First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading product data"})
		return
	}

	c.JSON(201, resp)

}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var req models.UpdateProductRequest
	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	var product models.Product

	if err := pc.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product is not exists"})
		return
	}

	if err := pc.DB.Where("name = ?", req.Name).First(&product).Error; err == nil {
		c.JSON(400, gin.H{"error": "Product name is already used"})
		return
	}

	tx := pc.DB.Begin()

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
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

	if err := pc.DB.First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading product data"})
		return
	}

	productResp := models.ProductDataResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	resp := models.ProductResponse{
		Message: "Product is successfully created",
		Data:    productResp,
	}

	c.JSON(200, resp)

}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	var product models.Product

	if err := pc.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product is not exists"})
		return
	}

	tx := pc.DB.Begin()

	if err := tx.Model(&product).Association("Inventories").Clear(); err != nil {
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

	id := c.Query("id")
	category := c.Query("category")

	if id != "" && category != "" {
		if err := pc.DB.Where("id = ? AND category = ?", id, category).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if id != "" {
		if err := pc.DB.Where("id = ?", id).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if category != "" {
		if err := pc.DB.Where("category = ?", category).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := pc.DB.Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	resp := make([]models.ProductDataResponse, len(products))
	for i, product := range products {
		resp[i] = models.ProductDataResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(200, gin.H{"data": resp})
}
