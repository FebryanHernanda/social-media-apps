package models

import "time"

type User struct {
	ID         int       `form:"id"`
	Email      string    `form:"email"`
	Password   string    `form:"password"`
	Name       string    `form:"name"`
	AvatarPath *string   `form:"avatar_path,omitempty"`
	Biography  *string   `form:"biography,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type AllUser struct {
	ID         int       `form:"id"`
	Email      string    `form:"email"`
	Name       string    `form:"name"`
	AvatarPath *string   `form:"avatar_path,omitempty"`
	Biography  *string   `form:"biography,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Follows struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	FollowedUserID int       `json:"followed_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateFollowRequest struct {
	FollowedUserID int `json:"followed_user_id" binding:"required"`
}

type Notifications struct {
	ID          int       `json:"id"`
	ActorID     int       `json:"actor_id"`
	ActorName   string    `json:"actor_name"`
	ActorAvatar *string   `json:"actor_avatar,omitempty"`
	Action      string    `json:"action"`
	PostID      *int      `json:"post_id,omitempty"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}
