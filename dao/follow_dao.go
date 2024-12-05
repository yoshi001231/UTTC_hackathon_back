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
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url, u.header_img_url FROM users u INNER JOIN followers f ON u.user_id = f.user_id WHERE f.following_user_id = ?`, userID)
	if err != nil {
		log.Printf("[follow_dao.go] 以下のフォロワー一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio, profileImgURL, headerImgURL sql.NullString

		// Scan 時に NullString を使用
		if err := rows.Scan(&user.UserID, &user.Name, &bio, &profileImgURL, &headerImgURL); err != nil {
			log.Printf("[follow_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}

		// NullString をポインタ型に変換
		user.Bio = nullableToPointer(bio)
		user.ProfileImgURL = nullableToPointer(profileImgURL)
		user.HeaderImgURL = nullableToPointer(headerImgURL)

		users = append(users, user)
	}
	return users, nil
}

func (dao *FollowDAO) GetFollowing(userID string) ([]model.User, error) {
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url, u.header_img_url FROM users u INNER JOIN followers f ON u.user_id = f.following_user_id WHERE f.user_id = ?`, userID)
	if err != nil {
		log.Printf("[follow_dao.go] 以下のフォロー中一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio, profileImgURL, headerImgURL sql.NullString

		// Scan 時に NullString を使用
		if err := rows.Scan(&user.UserID, &user.Name, &bio, &profileImgURL, &headerImgURL); err != nil {
			log.Printf("[follow_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}

		// NullString をポインタ型に変換
		user.Bio = nullableToPointer(bio)
		user.ProfileImgURL = nullableToPointer(profileImgURL)
		user.HeaderImgURL = nullableToPointer(headerImgURL)

		users = append(users, user)
	}
	return users, nil
}

// GetFollowGraph フォローグラフを取得
func (dao *FollowDAO) GetFollowGraph() ([]model.Follow, error) {
	rows, err := dao.db.Query("SELECT user_id, following_user_id FROM followers")
	if err != nil {
		log.Printf("[follow_dao.go] フォローグラフの取得失敗: %v", err)
		return nil, err
	}
	defer rows.Close()

	var follows []model.Follow
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.UserID, &follow.FollowingUserID); err != nil {
			log.Printf("[follow_dao.go] フォローグラフデータのScan失敗: %v", err)
			return nil, err
		}
		follows = append(follows, follow)
	}

	return follows, nil
}
