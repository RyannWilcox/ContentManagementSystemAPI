package models

import (
	"errors"
	"time"
)

// This struct includes fields for:
// - ID (unsigned integer, primary key)
// - Title (string, required, with max length)
// - Content (text field, required)
// - CreatedAt (timestamp for creation date)
// - UpdatedAt (timestamp for last update)

type Page struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title" binding:"required"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Page) Validate() error {
	var errorMsg error
	titleLen := len(p.Title)
	contentLen := len(p.Content)

	if titleLen == 0 {
		errorMsg = errors.New("page title cannot be empty")

	} else if titleLen > 255 {
		errorMsg = errors.New("page title cannot be more than 255 characters")

	} else if contentLen == 0 {
		errorMsg = errors.New("page content cannot be empty")
	}
	return errorMsg
}
