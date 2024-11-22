// model/models.go

package model

import "time"

// User モデル
type User struct {
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Bio           string `json:"bio"`
	ProfileImgURL string `json:"profile_img_url"`
}

// Post モデル
type Post struct {
	PostID       string     `json:"post_id"`
	UserID       string     `json:"user_id"`
	Content      string     `json:"content"`
	ImgURL       string     `json:"img_url"`
	CreatedAt    time.Time  `json:"created_at"`
	EditedAt     *time.Time `json:"edited_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	ParentPostID string     `json:"parent_post_id,omitempty"`
}
