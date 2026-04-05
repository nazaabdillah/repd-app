package controllers

import (
	"net/http"
	"repd-backend/config"
	"repd-backend/models"

	"github.com/gin-gonic/gin"
)

type SyncUserRequest struct {
	ID    string `json:"id" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func SyncUser(c *gin.Context) {
	var req SyncUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Cari user, jika tidak ada, buat baru (FirstOrCreate)
	result := config.DB.Where(models.User{ID: req.ID}).FirstOrCreate(&user, models.User{
		ID:    req.ID,
		Email: req.Email,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal sinkronisasi user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User tersinkronisasi",
		"user":    user,
	})
}