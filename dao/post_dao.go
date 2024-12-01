package dao

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"twitter/model"
)

type PostDAO struct {
	db *sql.DB
}

func NewPostDAO(db *sql.DB) *PostDAO {
	return &PostDAO{db: db}
}

// CreatePost 新しい投稿を作成
func (dao *PostDAO) CreatePost(post model.Post) (*model.Post, error) {
	var parentPostID interface{}
	if post.ParentPostID == "" {
		parentPostID = nil
	} else {
		parentPostID = post.ParentPostID
	}
	_, err := dao.db.Exec("INSERT INTO posts (post_id, user_id, content, img_url, created_at, parent_post_id) VALUES (?, ?, ?, ?, ?, ?)", post.PostID, post.UserID, post.Content, post.ImgURL, post.CreatedAt, parentPostID)
	if err != nil {
		log.Printf("[post_dao.go] 以下の投稿作成失敗 (post_id: %s, user_id: %s, content: %s, img_url: %s): %v", post.PostID, post.UserID, post.Content, post.ImgURL, err)
		return nil, err
	}
	return &post, nil
}

// GetPost 投稿の詳細を取得
func (dao *PostDAO) GetPost(postID string) (*model.Post, error) {
	var post model.Post
	var parentPostID sql.NullString
	var deletedAt sql.NullTime
	var editedAt sql.NullTime

	err := dao.db.QueryRow(`SELECT post_id, user_id, content, img_url, created_at, edited_at, parent_post_id, deleted_at  FROM posts WHERE post_id = ?`, postID).Scan(
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
		log.Printf("[post_dao.go] 以下の投稿が見つからない (post_id: %s)", postID)
		return nil, err
	} else if err != nil {
		log.Printf("[post_dao.go] 以下の投稿取得失敗 (post_id: %s): %v", postID, err)
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
func (dao *PostDAO) UpdatePost(post model.Post) error {
	editedAt := time.Now()
	_, err := dao.db.Exec("UPDATE posts SET content = ?, img_url = ?, edited_at = ? WHERE post_id = ? AND deleted_at IS NULL", post.Content, post.ImgURL, editedAt, post.PostID)
	if err != nil {
		log.Printf("[post_dao.go] 以下の投稿更新失敗 (post_id: %s, user_id: %s, content: %s, img_url: %s): %v", post.PostID, post.UserID, post.Content, post.ImgURL, err)
	}
	return err
}

// DeletePost 投稿を削除 (論理削除)
func (dao *PostDAO) DeletePost(postID string) error {
	_, err := dao.db.Exec("UPDATE posts SET deleted_at = ? WHERE post_id = ?", time.Now(), postID)
	if err != nil {
		log.Printf("[post_dao.go] 以下の投稿削除失敗 (post_id: %s): %v", postID, err)
	}
	return err
}

// GetChildrenPosts 子ポストを取得
func (dao *PostDAO) GetChildrenPosts(parentPostID string) ([]model.Post, error) {
	rows, err := dao.db.Query(`SELECT post_id, user_id, content, img_url, created_at, edited_at, parent_post_id, deleted_at FROM posts WHERE parent_post_id = ? AND deleted_at IS NULL`, parentPostID)
	if err != nil {
		log.Printf("[post_dao.go] 子ポスト一覧取得失敗 (parent_post_id: %s): %v", parentPostID, err)
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var editedAt sql.NullTime
		var parentPostID sql.NullString
		var deletedAt sql.NullTime

		if err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.ImgURL, &post.CreatedAt, &editedAt, &parentPostID, &deletedAt); err != nil {
			log.Printf("[post_dao.go] 子ポストデータのScan失敗: %v", err)
			return nil, err
		}

		// `edited_at` の処理
		if editedAt.Valid {
			post.EditedAt = &editedAt.Time
		}

		// `parent_post_id` の処理
		if parentPostID.Valid {
			post.ParentPostID = parentPostID.String
		}

		posts = append(posts, post)
	}
	return posts, nil
}
