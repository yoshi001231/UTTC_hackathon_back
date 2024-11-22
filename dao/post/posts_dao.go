// dao/post/posts_dao.go

package post

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"twitter/model"
)

type PostsDAO struct {
	db *sql.DB
}

func NewPostsDAO(db *sql.DB) *PostsDAO {
	return &PostsDAO{db: db}
}

// CreatePost 新しい投稿を作成
func (dao *PostsDAO) CreatePost(post model.Post) (*model.Post, error) {
	var parentPostID interface{}
	if post.ParentPostID == "" {
		parentPostID = nil
	} else {
		parentPostID = post.ParentPostID
	}

	_, err := dao.db.Exec(
		"INSERT INTO posts (post_id, user_id, content, img_url, created_at, parent_post_id) VALUES (?, ?, ?, ?, ?, ?)",
		post.PostID, post.UserID, post.Content, post.ImgURL, post.CreatedAt, parentPostID,
	)
	if err != nil {
		log.Printf("投稿作成失敗 (post_id: %s, user_id: %s): %v", post.PostID, post.UserID, err)
		return nil, err
	}
	return &post, nil
}

// GetPost 投稿の詳細を取得
func (dao *PostsDAO) GetPost(postID string) (*model.Post, error) {
	var post model.Post
	var parentPostID sql.NullString
	var deletedAt sql.NullTime
	var editedAt sql.NullTime

	err := dao.db.QueryRow(
		`SELECT post_id, user_id, content, img_url, created_at, edited_at, parent_post_id, deleted_at 
         FROM posts WHERE post_id = ?`,
		postID,
	).Scan(
		&post.PostID,
		&post.UserID,
		&post.Content,
		&post.ImgURL,
		&post.CreatedAt,
		&editedAt,
		&parentPostID,
		&deletedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("投稿が見つかりません (post_id: %s)", postID)
		return nil, errors.New("投稿が見つかりません")
	} else if err != nil {
		log.Printf("投稿取得失敗: %v", err)
		return nil, err
	}

	// 投稿が削除済みかどうか確認
	if deletedAt.Valid {
		return nil, errors.New("投稿が削除されています")
	}

	// parent_post_id の処理
	if parentPostID.Valid {
		post.ParentPostID = parentPostID.String
	} else {
		post.ParentPostID = ""
	}

	// edited_at の処理
	if editedAt.Valid {
		post.EditedAt = &editedAt.Time
	}

	return &post, nil
}

// UpdatePost 投稿を更新
func (dao *PostsDAO) UpdatePost(post model.Post) error {
	editedAt := time.Now()
	_, err := dao.db.Exec("UPDATE posts SET content = ?, img_url = ?, edited_at = ? WHERE post_id = ? AND deleted_at IS NULL",
		post.Content, post.ImgURL, editedAt, post.PostID)
	if err != nil {
		log.Printf("投稿更新失敗: %v", err)
	}
	return err
}

// DeletePost 投稿を削除 (論理削除)
func (dao *PostsDAO) DeletePost(postID string) error {
	_, err := dao.db.Exec("UPDATE posts SET deleted_at = ? WHERE post_id = ?", time.Now(), postID)
	if err != nil {
		log.Printf("投稿削除失敗: %v", err)
	}
	return err
}
