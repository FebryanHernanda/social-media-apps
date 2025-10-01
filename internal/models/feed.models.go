package models

import "time"

type FeedPost struct {
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	ImagePath    *string   `json:"image_path,omitempty"`
	AuthorID     int       `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	AvatarPath   *string   `json:"author_avatar,omitempty"`
	CreatedAt    time.Time `json:"-"`
	CreatedAtStr string    `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	Comments     []Comment `json:"comments"`
}
