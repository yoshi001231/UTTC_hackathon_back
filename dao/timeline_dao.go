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

// FetchUserTimeline ログインユーザーのタイムラインを取得
func (dao *TimelineDAO) FetchUserTimeline(userID string) ([]model.Post, error) {
	rows, err := dao.db.Query(`SELECT p.post_id, p.user_id, p.content, p.img_url, p.created_at, p.parent_post_id FROM posts p JOIN followers f ON p.user_id = f.following_user_id WHERE (f.user_id = ? OR p.user_id = ?) AND p.deleted_at IS NULL ORDER BY p.created_at DESC`, userID)
	if err != nil {
		log.Printf("[timeline_dao.go] 以下のタイムライン取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var parentPostID sql.NullString
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.ImgURL, &post.CreatedAt, &parentPostID); err != nil {
			log.Printf("[timeline_dao.go] 投稿データのScan失敗: %v", err)
			return nil, err
		}
		if parentPostID.Valid {
			post.ParentPostID = parentPostID.String
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// FetchUserPosts 指定ユーザーの投稿一覧を取得
func (dao *TimelineDAO) FetchUserPosts(userID string) ([]model.Post, error) {
	rows, err := dao.db.Query(`SELECT post_id, user_id, content, img_url, created_at, parent_post_id FROM posts WHERE user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Printf("[timeline_dao.go] 以下の投稿一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var parentPostID sql.NullString
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.ImgURL, &post.CreatedAt, &parentPostID); err != nil {
			log.Printf("[timeline_dao.go] 投稿データのScan失敗: %v", err)
			return nil, err
		}
		if parentPostID.Valid {
			post.ParentPostID = parentPostID.String
		}
		posts = append(posts, post)
	}
	return posts, nil
}
