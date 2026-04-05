package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"repd-backend/config"
	"repd-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	userID := c.PostForm("user_id") // Ambil text dari Form
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id wajib diisi"})
		return
	}

	// Ambil file dari request dengan key "photo"
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengambil file gambar"})
		return
	}

	// Buat nama file unik menggunakan timestamp agar tidak bentrok
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	savePath := fmt.Sprintf("uploads/photos/%s", filename)

	// Simpan file fisik ke folder "uploads/photos/"
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file secara fisik"})
		return
	}

	// Simpan URL/Path nya ke Database
	photoLog := models.ProgressPhoto{
		UserID:   userID,
		PhotoURL: "/" + savePath, // Contoh: /uploads/photos/123456_foto.jpg
	}

	if err := config.DB.Create(&photoLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat path ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Foto berhasil diunggah!",
		"data":    photoLog,
	})
}

func GetPhotos(c *gin.Context) {
	userID := c.Query("user_id")
	var photos []models.ProgressPhoto

	// Ambil foto dari yang terbaru
	config.DB.Where("user_id = ?", userID).Order("date desc").Find(&photos)

	c.JSON(http.StatusOK, gin.H{"data": photos})
}