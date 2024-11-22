// dao/auth/users_dao.go

package auth

import (
	"database/sql"
	"log"
	"twitter/model"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) RegisterUser(user model.User) error {
	_, err := dao.db.Exec("INSERT INTO users (user_id, name, bio, profile_img_url) VALUES (?, ?, ?, ?)",
		user.UserID, user.Name, user.Bio, user.ProfileImgURL)
	if err != nil {
		log.Printf("ユーザー登録失敗: %v", err)
	}
	return err
}
