package model

type User struct {
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Bio           string `json:"bio"`
	ProfileImgURL string `json:"profile_img_url"`
}
