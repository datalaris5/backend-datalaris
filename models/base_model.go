package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type BaseModel struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy *uint      `json:"created_by" gorm:"column:created_by"`
	UpdatedBy *uint      `json:"updated_by" gorm:"column:updated_by"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("🚀 [BeforeCreate] Triggered")
	if uid, ok := tx.Statement.Context.Value(UserIDKey).(uint); ok {
		fmt.Println("✅ Found user_id in context:", uid)
		b.CreatedBy = &uid
	} else {
		fmt.Println("⚠️ No user_id found in context")
	}

	now := time.Now()
	b.UpdatedAt = &now
	return
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("🚀 [BeforeUpdate] Triggered")
	if uid, ok := tx.Statement.Context.Value(UserIDKey).(uint); ok {
		fmt.Println("✅ Found user_id in context:", uid)
		b.UpdatedBy = &uid
	}
	now := time.Now()
	b.UpdatedAt = &now
	return
}
