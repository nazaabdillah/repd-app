package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"repd-backend/config"
	"repd-backend/controllers"
	"repd-backend/middlewares"

	"github.com/gin-contrib/cors" // <-- Import CORS
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()

	r := gin.Default()

	// --- SETUP CORS ---
	// Mengizinkan Next.js (port 3000) untuk mengakses API ini
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 2. SETUP FOLDER STATIS (DI BAWAH CORS)
	// Gunakan http.Dir agar pathing-nya absolut dan aman
	r.StaticFS("/uploads", http.Dir("uploads"))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "REPD API is running."})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "online", "message": "REPD API is running on cloud"})
	})

	api := r.Group("/api/v1")
	api.Use(middlewares.RequireAuth())
	{
		api.POST("/users/sync", controllers.SyncUser)
		api.POST("/weights", controllers.AddWeight)
		api.GET("/weights", controllers.GetWeights)
		api.POST("/workouts", controllers.CreateWorkout)
		api.GET("/workouts", controllers.GetWorkouts)
		api.POST("/photos", controllers.UploadPhoto)
		api.GET("/photos", controllers.GetPhotos)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback jika dijalankan di lokal
	}
	fmt.Println("🚀 Server berjalan di port:", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
