package controllers

import (
	"net/http"
	"repd-backend/config"
	"repd-backend/models"

	"github.com/gin-gonic/gin"
)

type AddWeightRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Weight float64 `json:"weight" binding:"required,gt=20,lt=300"` // Validasi logika dasar
}

func AddWeight(c *gin.Context) {
	var req AddWeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weightLog := models.WeightLog{
		UserID: req.UserID,
		Weight: req.Weight,
	}

	if err := config.DB.Create(&weightLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data berat badan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berat badan berhasil dicatat",
		"data":    weightLog,
	})
}

func GetWeights(c *gin.Context) {
	userID := c.Query("user_id") // Mengambil dari URL query: /weights?user_id=xxx
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id diperlukan"})
		return
	}

	var weights []models.WeightLog
	// Order by Date descending (terbaru di atas)
	if err := config.DB.Where("user_id = ?", userID).Order("date desc").Find(&weights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": weights,
	})
}