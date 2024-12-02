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
	// ポインタ型のフィールドを確認して値を取得
	var bio, profileImgURL interface{}
	if user.Bio != nil {
		bio = *user.Bio
	} else {
		bio = nil
	}
	if user.ProfileImgURL != nil {
		profileImgURL = *user.ProfileImgURL
	} else {
		profileImgURL = nil
	}

	_, err := dao.db.Exec(
		"INSERT INTO users (user_id, name, bio, profile_img_url) VALUES (?, ?, ?, ?)",
		user.UserID,
		user.Name,
		bio,
		profileImgURL,
	)
	if err != nil {
		log.Printf("[auth_dao.go] 以下のユーザー登録失敗 (user_id: %s, name: %s): %v", user.UserID, user.Name, err)
	}
	return err
}
