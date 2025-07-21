package models

import "time"

type Ad struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    *string   `json:"image_url,omitempty"`
	Price       uint      `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	IsMine      bool      `json:"isMine,omitempty"`
}
