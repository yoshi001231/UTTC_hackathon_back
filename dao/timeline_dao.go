// dao/timeline_dao.go

package dao

import (
	"database/sql"
	"log"
	"twitter/model"
)

type TimelineDAO struct {
	db *sql.DB
}

func NewTimelineDAO(db *sql.DB) *TimelineDAO {
	return &TimelineDAO{db: db}
}

// GetUserTimeline ログインユーザーのタイムラインを取得
func (dao *TimelineDAO) GetUserTimeline(userID string) ([]model.Post, error) {
	rows, err := dao.db.Query(
		`SELECT p.post_id, p.user_id, p.content, p.img_url, p.created_at, p.edited_at, p.parent_post_id 
		 FROM posts p
		 INNER JOIN followers f ON p.user_id = f.following_user_id
		 WHERE f.user_id = ? AND p.deleted_at IS NULL
		 ORDER BY p.created_at DESC`,
		userID,
	)
	if err != nil {
		log.Printf("タイムライン取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var parentPostID sql.NullString

		err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.ImgURL,
			&post.CreatedAt,
			&post.EditedAt,
			&parentPostID,
		)
		if err != nil {
			log.Printf("タイムライン投稿データの読み取り失敗: %v", err)
			return nil, err
		}
		if parentPostID.Valid {
			post.ParentPostID = parentPostID.String
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetUserPosts 指定ユーザーの投稿一覧を取得
func (dao *TimelineDAO) GetUserPosts(userID string) ([]model.Post, error) {
	rows, err := dao.db.Query(
		`SELECT post_id, user_id, content, img_url, created_at, edited_at, parent_post_id
		 FROM posts
		 WHERE user_id = ? AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		log.Printf("ユーザー投稿一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var parentPostID sql.NullString

		err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.ImgURL,
			&post.CreatedAt,
			&post.EditedAt,
			&parentPostID,
		)
		if err != nil {
			log.Printf("ユーザー投稿データの読み取り失敗: %v", err)
			return nil, err
		}
		if parentPostID.Valid {
			post.ParentPostID = parentPostID.String
		}
		posts = append(posts, post)
	}
	return posts, nil
}
