// dao/user_dao.go

package dao

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

func (dao *UserDAO) GetUser(userID string) (*model.User, error) {
	var user model.User
	err := dao.db.QueryRow("SELECT user_id, name, bio, profile_img_url FROM users WHERE user_id = ?", userID).Scan(
		&user.UserID, &user.Name, &user.Bio, &user.ProfileImgURL,
	)
	if err != nil {
		log.Printf("ユーザー取得失敗: %v", err)
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) UpdateUser(user model.User) error {
	_, err := dao.db.Exec("UPDATE users SET name = ?, bio = ?, profile_img_url = ? WHERE user_id = ?",
		user.Name, user.Bio, user.ProfileImgURL, user.UserID)
	if err != nil {
		log.Printf("ユーザー更新失敗: %v", err)
	}
	return err
}
