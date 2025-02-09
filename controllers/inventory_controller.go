package controllers

import (
	"inventory-management/models"
	"inventory-management/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryController struct {
	DB *gorm.DB
}

func NewInventoryController(db *gorm.DB) *InventoryController {
	return &InventoryController{DB: db}
}

func (ic *InventoryController) UpdateInventory(c *gin.Context) {
	var req models.UpdateInventoryRequest

	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	var inventory models.Inventory
	if err := ic.DB.Where("id = ?", c.Param("id")).First(&inventory).Error; err != nil {
		c.JSON(404, gin.H{"error": "Inventory is not exists"})
		return
	}

	tx := ic.DB.Begin()

	inventory.Stock = req.Stock

	if err := tx.Save(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	if err := ic.DB.First(&inventory).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading inventory data"})
		return
	}

	inventoryResp := models.InventoryDataResponse{
		ID:        inventory.ID,
		Stock:     inventory.Stock,
		Location:  inventory.Location,
		ProductId: inventory.ProductId,
	}

	resp := models.InventoryResponse{
		Message: "Stock updated",
		Data:    inventoryResp,
	}

	c.JSON(200, resp)

}

func (ic *InventoryController) GetInventory(c *gin.Context) {
	var inventories []models.Inventory

	inventoryId := c.Query("inventoryId")
	productId := c.Query("productId")

	if inventoryId != "" && productId != "" {
		if err := ic.DB.Where("id = ? AND product_id = ?", inventoryId, productId).Find(&inventories).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if inventoryId != "" {
		if err := ic.DB.Where("id = ?", inventoryId).Find(&inventories).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if productId != "" {
		if err := ic.DB.Where("product_id = ?", productId).Find(&inventories).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := ic.DB.Find(&inventories).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	resp := make([]models.InventoryDataResponse, len(inventories))
	for i, inventory := range inventories {
		resp[i] = models.InventoryDataResponse{
			ID:        inventory.ID,
			Stock:     inventory.Stock,
			Location:  inventory.Location,
			ProductId: inventory.ProductId,
		}
	}

	c.JSON(200, gin.H{"data": resp})
}
