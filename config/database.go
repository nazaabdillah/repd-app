package config

import (
	"fmt"
	"log"
	"os"
	"repd-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Ambil variabel dan beri default jika kosong untuk port
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if port == "" {
		port = "5432"
	}

	// Logging untuk debug di Cloud (Hanya host agar password tetap aman)
	fmt.Printf("Attempting connection to host: %s port: %s\n", host, port)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, pass, name, port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		// PENTING: Gunakan log.Println, JANGAN log.Fatal
		// Supaya aplikasi tidak mati dan Health Check Back4App tetap jalan
		log.Println("🚨 Gagal terkoneksi ke database! Error:", err)
		return 
	}

	fmt.Println("✅ Koneksi database berhasil!")

	// Jalankan Auto Migration
	err = database.AutoMigrate(
		&models.User{},
		&models.WorkoutSession{},
		&models.ExerciseLog{},
		&models.WeightLog{},
		&models.ProgressPhoto{},
	)
	if err != nil {
		log.Println("🚨 Gagal menjalankan migrasi database:", err)
		return
	}
	
	fmt.Println("✅ Tabel Database berhasil disinkronisasi!")
	DB = database
}