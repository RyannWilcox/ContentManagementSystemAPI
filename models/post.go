package models

import (
	"errors"
	"time"
)

// This struct includes fields for:
// - ID (unsigned integer, primary key)
// - Title (string, required, with max length)
// - Content (text field, required)
// - Author (string, optional)
// - CreatedAt (timestamp for creation date)
// - UpdatedAt (timestamp for last update)
// - Media (slice of Media, representing a many-to-many relationship)

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title" binding:"required"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	Author    string    `gorm:"size:100" json:"author"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Media     []Media   `gorm:"many2many:post_media" json:"media"`
}

func (p *Post) Validate() error {
	var errorMsg error

	if len(p.Title) == 0 {
		errorMsg = errors.New("post title cannot be empty")
	} else if len(p.Title) > 255 {
		errorMsg = errors.New("post title cannot be more than 255 characters")
	} else if len(p.Content) == 0 {
		errorMsg = errors.New("post conent cannot be empty")
	}

	return errorMsg
}
