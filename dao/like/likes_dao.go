// dao/like/likes_dao.go

package like

import (
	"database/sql"
	"log"
	"time"
	"twitter/model"
)

type LikesDAO struct {
	db *sql.DB
}

func NewLikesDAO(db *sql.DB) *LikesDAO {
	return &LikesDAO{db: db}
}

// AddLike 投稿にいいねを追加
func (dao *LikesDAO) AddLike(userID, postID string) error {
	_, err := dao.db.Exec(
		"INSERT INTO likes (user_id, post_id, created_at) VALUES (?, ?, ?)",
		userID, postID, time.Now(),
	)
	if err != nil {
		log.Printf("いいね追加失敗 (user_id: %s, post_id: %s): %v", userID, postID, err)
	}
	return err
}

// RemoveLike 投稿のいいねを削除
func (dao *LikesDAO) RemoveLike(userID, postID string) error {
	_, err := dao.db.Exec(
		"DELETE FROM likes WHERE user_id = ? AND post_id = ?",
		userID, postID,
	)
	if err != nil {
		log.Printf("いいね削除失敗 (user_id: %s, post_id: %s): %v", userID, postID, err)
	}
	return err
}

// GetUsersByPostID 投稿にいいねしたユーザー一覧を取得
func (dao *LikesDAO) GetUsersByPostID(postID string) ([]model.User, error) {
	rows, err := dao.db.Query(
		`SELECT u.user_id, u.name, u.bio, u.profile_img_url 
         FROM users u
         INNER JOIN likes l ON u.user_id = l.user_id
         WHERE l.post_id = ?`,
		postID,
	)
	if err != nil {
		log.Printf("いいねユーザー一覧取得失敗 (post_id: %s): %v", postID, err)
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
