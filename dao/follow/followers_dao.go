// dao/follow/followers_dao.go

package follow

import (
	"database/sql"
	"log"
	"time"
	"twitter/model"
)

type FollowersDAO struct {
	db *sql.DB
}

func NewFollowersDAO(db *sql.DB) *FollowersDAO {
	return &FollowersDAO{db: db}
}

// AddFollow 指定ユーザーをフォロー
func (dao *FollowersDAO) AddFollow(userID, followingUserID string) error {
	_, err := dao.db.Exec(
		"INSERT INTO followers (user_id, following_user_id, created_at) VALUES (?, ?, ?)",
		userID, followingUserID, time.Now(),
	)
	if err != nil {
		log.Printf("フォロー追加失敗 (user_id: %s, following_user_id: %s): %v", userID, followingUserID, err)
	}
	return err
}

// RemoveFollow 指定ユーザーのフォローを解除
func (dao *FollowersDAO) RemoveFollow(userID, followingUserID string) error {
	_, err := dao.db.Exec(
		"DELETE FROM followers WHERE user_id = ? AND following_user_id = ?",
		userID, followingUserID,
	)
	if err != nil {
		log.Printf("フォロー解除失敗 (user_id: %s, following_user_id: %s): %v", userID, followingUserID, err)
	}
	return err
}

// GetFollowers 指定ユーザーのフォロワー一覧を取得
func (dao *FollowersDAO) GetFollowers(userID string) ([]model.User, error) {
	rows, err := dao.db.Query(
		`SELECT u.user_id, u.name, u.bio, u.profile_img_url 
         FROM users u
         INNER JOIN followers f ON u.user_id = f.user_id
         WHERE f.following_user_id = ?`,
		userID,
	)
	if err != nil {
		log.Printf("フォロワー一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Bio, &user.ProfileImgURL); err != nil {
			log.Printf("ユーザーデータの読み取り失敗: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetFollowing 指定ユーザーのフォロー中一覧を取得
func (dao *FollowersDAO) GetFollowing(userID string) ([]model.User, error) {
	rows, err := dao.db.Query(
		`SELECT u.user_id, u.name, u.bio, u.profile_img_url 
         FROM users u
         INNER JOIN followers f ON u.user_id = f.following_user_id
         WHERE f.user_id = ?`,
		userID,
	)
	if err != nil {
		log.Printf("フォロー中一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Bio, &user.ProfileImgURL); err != nil {
			log.Printf("ユーザーデータの読み取り失敗: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
