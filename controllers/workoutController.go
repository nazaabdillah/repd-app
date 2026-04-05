package controllers

import (
	"net/http"
	"repd-backend/config"
	"repd-backend/models"

	"github.com/gin-gonic/gin"
)

type ExerciseRequest struct {
	Name          string `json:"name" binding:"required"`
	Sets          int    `json:"sets" binding:"required"`
	Reps          int    `json:"reps" binding:"required"`
	Done          int    `json:"done"`
	RPE           *int   `json:"rpe"`
	Notes         string `json:"notes"`
}

type WorkoutRequest struct {
	UserID     string            `json:"user_id" binding:"required"`
	TemplateID string            `json:"template_id" binding:"required"`
	DayName    string            `json:"day_name" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Exercises  []ExerciseRequest `json:"exercises" binding:"required"`
}

func CreateWorkout(c *gin.Context) {
	var req WorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// MEMULAI TRANSAKSI DATABASE
	tx := config.DB.Begin()

	// 1. Simpan Sesi Utama
	session := models.WorkoutSession{
		UserID:     req.UserID,
		TemplateID: req.TemplateID,
		DayName:    req.DayName,
		Type:       req.Type,
	}

	if err := tx.Create(&session).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan sesi"})
		return
	}

	// 2. Simpan Semua Detail Latihan (Looping)
	for _, ex := range req.Exercises {
		log := models.ExerciseLog{
			SessionID:     session.ID, // Ambil ID dari sesi yang baru saja di-insert
			ExerciseName:  ex.Name,
			TargetSets:    ex.Sets,
			TargetReps:    ex.Reps,
			CompletedSets: ex.Done,
			RPE:           ex.RPE,
			Notes:         ex.Notes,
		}

		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback() // Batalkan semua jika satu gerakan gagal
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan log gerakan"})
			return
		}
	}

	// SELESAI & SIMPAN PERMANEN
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sesi latihan berhasil disimpan ke Database",
		"id":      session.ID,
	})
}

func GetWorkouts(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id diperlukan"})
		return
	}

	var sessions []models.WorkoutSession
	
	// Preload("Exercises") akan otomatis mengisi array Exercises di dalam struct Session
	// Order by date DESC agar latihan terbaru muncul paling atas
	err := config.DB.Preload("Exercises").
		Where("user_id = ?", userID).
		Order("date desc").
		Find(&sessions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat latihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": sessions,
	})
}