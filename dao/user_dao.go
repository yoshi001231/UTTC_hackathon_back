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
		&user.UserID,
		&user.Name,
		&user.Bio,
		&user.ProfileImgURL,
	)
	if err != nil {
		log.Printf("[user_dao.go] 以下のユーザー取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) UpdateUser(user model.User) error {
	_, err := dao.db.Exec("UPDATE users SET name = ?, bio = ?, profile_img_url = ? WHERE user_id = ?", user.Name, user.Bio, user.ProfileImgURL, user.UserID)
	if err != nil {
		log.Printf("[user_dao.go] 以下のユーザー更新失敗 (user_id: %s, name: %s, bio: %s, profile_img_url: %s): %v", user.UserID, user.Name, user.Bio, user.ProfileImgURL, err)
	}
	return err
}

// GetTopUsersByTweetCount ツイート数の多い順にユーザ一覧を取得
func (dao *UserDAO) GetTopUsersByTweetCount(limit int) ([]model.User, error) {
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url, COUNT(p.post_id) AS tweet_count FROM users u LEFT JOIN posts p ON u.user_id = p.user_id AND p.deleted_at IS NULL GROUP BY u.user_id ORDER BY tweet_count DESC LIMIT ?`, limit)
	if err != nil {
		log.Printf("[user_dao.go] ツイート数順ユーザ取得失敗: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Bio, &user.ProfileImgURL, &user.TweetCount); err != nil {
			log.Printf("[user_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetTopUsersByLikes ユーザ一覧をいいね数の多い順に取得
func (dao *UserDAO) GetTopUsersByLikes(limit int) ([]model.User, error) {
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url, COUNT(l.post_id) AS like_count FROM users u LEFT JOIN posts p ON u.user_id = p.user_id LEFT JOIN likes l ON p.post_id = l.post_id WHERE p.deleted_at IS NULL GROUP BY u.user_id ORDER BY like_count DESC LIMIT ?`, limit)
	if err != nil {
		log.Printf("[user_dao.go] いいね数順ユーザ取得失敗: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Bio, &user.ProfileImgURL, &user.LikeCount); err != nil {
			log.Printf("[user_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
