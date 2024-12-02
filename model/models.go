package model

import "time"

// User モデル
type User struct {
	UserID        string     `json:"user_id"`
	Name          string     `json:"name"`
	Bio           string     `json:"bio"`
	ProfileImgURL string     `json:"profile_img_url"`
	HeaderImgURL  string     `json:"header_img_url"`
	Location      string     `json:"location"`
	Birthday      *time.Time `json:"birthday,omitempty"`
	TweetCount    int        `json:"tweet_count,omitempty"`
	LikeCount     int        `json:"like_count,omitempty"`
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

// Like モデル
type Like struct {
	UserID string `json:"user_id"`
	PostID string `json:"post_id,omitempty"`
}

// Follow モデル
type Follow struct {
	UserID          string `json:"user_id"`
	FollowingUserID string `json:"following_user_id"`
}
