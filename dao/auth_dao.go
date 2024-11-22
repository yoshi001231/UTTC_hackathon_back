// dao/auth_dao.go

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
	_, err := dao.db.Exec("INSERT INTO users (user_id, name, bio, profile_img_url) VALUES (?, ?, ?, ?)",
		user.UserID, user.Name, user.Bio, user.ProfileImgURL)
	if err != nil {
		log.Printf("ユーザー登録失敗: %v", err)
	}
	return err
}
