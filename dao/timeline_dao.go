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
	rows, err := dao.db.Query(`
		SELECT p.post_id, p.user_id, p.content, p.img_url, p.created_at, p.edited_at, p.parent_post_id, p.is_bad 
		FROM posts p 
		WHERE p.deleted_at IS NULL 
		AND (p.user_id = ? OR EXISTS (
			SELECT 1 FROM followers f WHERE f.user_id = ? AND f.following_user_id = p.user_id
		)) 
		ORDER BY p.created_at DESC`, userID, userID)
	if err != nil {
		log.Printf("[timeline_dao.go] 以下のタイムライン取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var imgURL, parentPostID sql.NullString
		var editedAt sql.NullTime

		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&imgURL,
			&post.CreatedAt,
			&editedAt,
			&parentPostID,
			&post.IsBad,
		); err != nil {
			log.Printf("[timeline_dao.go] 投稿データのScan失敗: %v", err)
			return nil, err
		}

		// NULL 値の処理
		post.ImgURL = nullableToPointer(imgURL)
		post.ParentPostID = nullableToPointer(parentPostID)
		if editedAt.Valid {
			post.EditedAt = &editedAt.Time
		}

		posts = append(posts, post)
	}
	return posts, nil
}

// FetchUserPosts 指定ユーザーの投稿一覧を取得
func (dao *TimelineDAO) FetchUserPosts(userID string) ([]model.Post, error) {
	rows, err := dao.db.Query(`
		SELECT post_id, user_id, content, img_url, created_at, edited_at, parent_post_id, is_bad 
		FROM posts 
		WHERE user_id = ? AND deleted_at IS NULL 
		ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Printf("[timeline_dao.go] 以下の投稿一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var imgURL, parentPostID sql.NullString
		var editedAt sql.NullTime

		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&imgURL,
			&post.CreatedAt,
			&editedAt,
			&parentPostID,
			&post.IsBad,
		); err != nil {
			log.Printf("[timeline_dao.go] 投稿データのScan失敗: %v", err)
			return nil, err
		}

		// NULL 値の処理
		post.ImgURL = nullableToPointer(imgURL)
		post.ParentPostID = nullableToPointer(parentPostID)
		if editedAt.Valid {
			post.EditedAt = &editedAt.Time
		}

		posts = append(posts, post)
	}
	return posts, nil
}

// FetchLikedPosts 指定ユーザーのいいねした投稿一覧を取得
func (dao *TimelineDAO) FetchLikedPosts(userID string) ([]model.Post, error) {
	rows, err := dao.db.Query(`
		SELECT p.post_id, p.user_id, p.content, p.img_url, p.created_at, p.edited_at, p.parent_post_id, p.is_bad 
		FROM posts p
		JOIN likes l ON p.post_id = l.post_id
		WHERE l.user_id = ? AND p.deleted_at IS NULL
		ORDER BY l.created_at DESC`, userID)
	if err != nil {
		log.Printf("[timeline_dao.go] いいねした投稿取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var imgURL, parentPostID sql.NullString
		var editedAt sql.NullTime

		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&imgURL,
			&post.CreatedAt,
			&editedAt,
			&parentPostID,
			&post.IsBad,
		); err != nil {
			log.Printf("[timeline_dao.go] 投稿データのScan失敗: %v", err)
			return nil, err
		}

		// NULL 値の処理
		post.ImgURL = nullableToPointer(imgURL)
		post.ParentPostID = nullableToPointer(parentPostID)
		if editedAt.Valid {
			post.EditedAt = &editedAt.Time
		}

		posts = append(posts, post)
	}
	return posts, nil
}
