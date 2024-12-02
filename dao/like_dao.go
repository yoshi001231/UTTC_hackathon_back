package dao

import (
	"database/sql"
	"log"
	"time"
	"twitter/model"
)

type LikeDAO struct {
	db *sql.DB
}

func NewLikeDAO(db *sql.DB) *LikeDAO {
	return &LikeDAO{db: db}
}

// AddLike 投稿にいいねを追加
func (dao *LikeDAO) AddLike(userID, postID string) error {
	_, err := dao.db.Exec("INSERT INTO likes (user_id, post_id, created_at) VALUES (?, ?, ?)", userID, postID, time.Now())
	if err != nil {
		log.Printf("[like_dao.go] 以下のいいね追加失敗 (user_id: %s, post_id: %s): %v", userID, postID, err)
	}
	return err
}

// RemoveLike 投稿のいいねを削除
func (dao *LikeDAO) RemoveLike(userID, postID string) error {
	_, err := dao.db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		log.Printf("[like_dao.go] 以下のいいね削除失敗 (user_id: %s, post_id: %s): %v", userID, postID, err)
	}
	return err
}

// GetUsersByPostID 投稿にいいねしたユーザー一覧を取得
func (dao *LikeDAO) GetUsersByPostID(postID string) ([]model.User, error) {
	rows, err := dao.db.Query(`SELECT u.user_id, u.name, u.bio, u.profile_img_url FROM users u INNER JOIN likes l ON u.user_id = l.user_id WHERE l.post_id = ?`, postID)
	if err != nil {
		log.Printf("[like_dao.go] 以下のいいねユーザー一覧取得失敗 (post_id: %s): %v", postID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio, profileImgURL sql.NullString

		// Scan 時に NullString を使用
		if err := rows.Scan(&user.UserID, &user.Name, &bio, &profileImgURL); err != nil {
			log.Printf("[like_dao.go] ユーザーデータのScan失敗: %v", err)
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
