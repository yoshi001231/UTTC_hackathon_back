// dao/auth/users_dao.go

package auth

import (
	"database/sql"
	"log"
	"twitter/model"
)

type UsersDAO struct {
	db *sql.DB
}

func NewUsersDAO(db *sql.DB) *UsersDAO {
	return &UsersDAO{db: db}
}

func (dao *UsersDAO) RegisterUser(user model.User) error {
	_, err := dao.db.Exec("INSERT INTO users (user_id, name, bio, profile_img_url) VALUES (?, ?, ?, ?)",
		user.UserID, user.Name, user.Bio, user.ProfileImgURL)
	if err != nil {
		log.Printf("ユーザー登録失敗: %v", err)
	}
	return err
}
