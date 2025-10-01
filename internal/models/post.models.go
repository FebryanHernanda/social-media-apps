package models

import (
	"mime/multipart"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	ImagePath *string   `json:"image_path,omitempty"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePostRequest struct {
	Content string                `form:"content" binding:"required"`
	Image   *multipart.FileHeader `form:"image"`
}

type Like struct {
	ID      int       `json:"id"`
	PostID  int       `json:"post_id"`
	UserID  int       `json:"user_id"`
	LikedAt time.Time `json:"liked_at"`
}

type Comment struct {
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	Name         *string   `json:"name,omitempty"`
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id"`
	CreatedAt    time.Time `json:"-"`
	CreatedAtStr string    `json:"created_at"`
}

type CommentRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}
