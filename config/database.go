package config

import (
	"fmt"
	"log"
	"os"
	"repd-backend/models" // <-- Sesuaikan dengan nama module go.mod kamu jika berbeda

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Gagal terkoneksi ke database! Error: ", err)
	}

	fmt.Println("✅ Koneksi database berhasil!")

	// Menjalankan Auto Migration
	err = database.AutoMigrate(
		&models.User{},
		&models.WorkoutSession{},
		&models.ExerciseLog{},
		&models.WeightLog{},
		&models.ProgressPhoto{},
	)
	if err != nil {
		log.Fatal("Gagal menjalankan migrasi database: ", err)
	}
	
	fmt.Println("✅ Tabel Database berhasil disinkronisasi!")

	DB = database
}