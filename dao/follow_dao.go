package dao

import (
	"database/sql"
	"log"
	"time"
	"twitter/model"
)

type FollowDAO struct {
	db *sql.DB
}

func NewFollowDAO(db *sql.DB) *FollowDAO {
	return &FollowDAO{db: db}
}

func (dao *FollowDAO) AddFollow(userID, followingUserID string) error {
	_, err := dao.db.Exec("INSERT INTO followers (user_id, following_user_id, created_at) VALUES (?, ?, ?)", userID, followingUserID, time.Now())
	if err != nil {
		log.Printf("[follow_dao.go] 以下のフォロー追加失敗 (user_id: %s, following_user_id: %s): %v", userID, followingUserID, err)
	}
	return err
}

func (dao *FollowDAO) RemoveFollow(userID, followingUserID string) error {
	_, err := dao.db.Exec("DELETE FROM followers WHERE user_id = ? AND following_user_id = ?", userID, followingUserID)
	if err != nil {
		log.Printf("[follow_dao.go] 以下のフォロー解除失敗 (user_id: %s, following_user_id: %s): %v", userID, followingUserID, err)
	}
	return err
}

func (dao *FollowDAO) GetFollowers(userID string) ([]model.User, error) {
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url FROM users u INNER JOIN followers f ON u.user_id = f.user_id WHERE f.following_user_id = ?`, userID)
	if err != nil {
		log.Printf("[follow_dao.go] 以下のフォロワー一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio, profileImgURL sql.NullString

		// Scan 時に NullString を使用
		if err := rows.Scan(&user.UserID, &user.Name, &bio, &profileImgURL); err != nil {
			log.Printf("[follow_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}

		// NullString をポインタ型に変換
		if bio.Valid {
			user.Bio = &bio.String
		} else {
			user.Bio = nil
		}
		if profileImgURL.Valid {
			user.ProfileImgURL = &profileImgURL.String
		} else {
			user.ProfileImgURL = nil
		}

		users = append(users, user)
	}
	return users, nil
}

func (dao *FollowDAO) GetFollowing(userID string) ([]model.User, error) {
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url FROM users u INNER JOIN followers f ON u.user_id = f.following_user_id WHERE f.user_id = ?`, userID)
	if err != nil {
		log.Printf("[follow_dao.go] 以下のフォロー中一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio, profileImgURL sql.NullString

		// Scan 時に NullString を使用
		if err := rows.Scan(&user.UserID, &user.Name, &bio, &profileImgURL); err != nil {
			log.Printf("[follow_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}

		// NullString をポインタ型に変換
		if bio.Valid {
			user.Bio = &bio.String
		} else {
			user.Bio = nil
		}
		if profileImgURL.Valid {
			user.ProfileImgURL = &profileImgURL.String
		} else {
			user.ProfileImgURL = nil
		}

		users = append(users, user)
	}
	return users, nil
}
