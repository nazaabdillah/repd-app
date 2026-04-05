package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"` // Menyimpan ID dari Clerk
	Email     string    `gorm:"unique;not null;type:varchar(255)" json:"email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type WorkoutSession struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	UserID     string        `gorm:"type:varchar(255);not null" json:"user_id"`
	User       User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	TemplateID string        `gorm:"type:varchar(50);not null" json:"template_id"`
	DayName    string        `gorm:"type:varchar(20);not null" json:"day_name"`
	Type       string        `gorm:"type:varchar(50);not null" json:"type"`
	Date       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	// TAMBAHKAN BARIS INI:
	Exercises  []ExerciseLog `gorm:"foreignKey:SessionID" json:"exercises"` 
}

type ExerciseLog struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	SessionID     uint           `gorm:"not null" json:"session_id"`
	Session       WorkoutSession `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;" json:"-"`
	ExerciseName  string         `gorm:"type:varchar(100);not null" json:"exercise_name"`
	TargetSets    int            `gorm:"not null" json:"target_sets"`
	TargetReps    int            `gorm:"not null" json:"target_reps"`
	CompletedSets int            `gorm:"default:0" json:"completed_sets"`
	RPE           *int           `json:"rpe"` // Pointer pointer (*) agar bisa menerima nilai null
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type WeightLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(255);not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	Weight    float64   `gorm:"type:decimal(5,2);not null" json:"weight"`
	Date      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
}

type ProgressPhoto struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(255);not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	PhotoURL  string    `gorm:"type:varchar(500);not null" json:"photo_url"`
	Date      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
}