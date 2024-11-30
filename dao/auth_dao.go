package dao

import (
	"database/sql"
	"log"
	"twitter/model"
)

type AuthDAO struct {
	db *sql.DB
}

func NewAuthDAO(db *sql.DB) *AuthDAO {
	return &AuthDAO{db: db}
}

func (dao *AuthDAO) RegisterUser(user model.User) error {
	_, err := dao.db.Exec("INSERT INTO users (user_id, name, bio, profile_img_url) VALUES (?, ?, ?, ?)", user.UserID, user.Name, user.Bio, user.ProfileImgURL)
	if err != nil {
		log.Printf("[auth_dao.go] 以下のユーザー登録失敗 (user_id: %s, name: %s, bio: %s, profile_img_url: %s): %v", user.UserID, user.Name, user.Bio, user.ProfileImgURL, err)
	}
	return err
}
