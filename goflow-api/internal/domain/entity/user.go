package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;type:varchar(100)" json:"username"`
	Email     string    `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	Password  string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
