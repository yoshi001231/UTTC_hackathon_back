package dao

import (
	"database/sql"
	"log"
	"twitter/model"
)

// FindDAO 検索用のDAO
type FindDAO struct {
	db *sql.DB
}

func NewFindDAO(db *sql.DB) *FindDAO {
	return &FindDAO{db: db}
}

// FindUsersByKey 指定したキーワードを name または bio に含むユーザーを検索
func (dao *FindDAO) FindUsersByKey(key string) ([]model.User, error) {
	rows, err := dao.db.Query(`
		SELECT user_id, name, bio, profile_img_url, header_img_url
		FROM users
		WHERE name LIKE ? OR bio LIKE ?`,
		"%"+key+"%", "%"+key+"%",
	)
	if err != nil {
		log.Printf("[find_dao.go] ユーザー検索失敗 (key: %s): %v", key, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio, profileImgURL, headerImgURL sql.NullString

		if err := rows.Scan(
			&user.UserID,
			&user.Name,
			&bio,
			&profileImgURL,
			&headerImgURL,
		); err != nil {
			log.Printf("[find_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}

		user.Bio = nullableToPointer(bio)
		user.ProfileImgURL = nullableToPointer(profileImgURL)
		user.HeaderImgURL = nullableToPointer(headerImgURL)

		users = append(users, user)
	}
	return users, nil
}

// FindPostsByKey 指定したキーワードを content に含む投稿を検索
func (dao *FindDAO) FindPostsByKey(key string) ([]model.Post, error) {
	rows, err := dao.db.Query(`
		SELECT post_id, user_id, content, img_url, created_at, edited_at, parent_post_id
		FROM posts
		WHERE content LIKE ? AND deleted_at IS NULL`,
		"%"+key+"%",
	)
	if err != nil {
		log.Printf("[find_dao.go] 投稿検索失敗 (key: %s): %v", key, err)
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
		); err != nil {
			log.Printf("[find_dao.go] 投稿データのScan失敗: %v", err)
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
