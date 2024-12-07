package dao

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"database/sql"
	"fmt"
	"log"
	"twitter/model"
)

const (
	location  = "asia-northeast1"
	modelName = "gemini-1.5-flash-002"
	projectID = "term6-yoshiaki-tanabe" // ① 自分のプロジェクトIDを指定する
)

type GeminiDAO struct {
	db *sql.DB
}

func NewGeminiDAO(db *sql.DB) *GeminiDAO {
	return &GeminiDAO{db: db}
}

// GenerateResponseFromPrompt Geminiを使用してプロンプトに対するレスポンスを生成
func (dao *GeminiDAO) GenerateResponseFromPrompt(prompt string) (*genai.Part, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("Geminiクライアントの初期化失敗: %w", err)
	}

	gemini := client.GenerativeModel(modelName)
	resp, err := gemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("Geminiによる生成失敗: %w", err)
	}

	// Candidates配列を確認
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("Geminiからの応答が空です")
	}

	return &resp.Candidates[0].Content.Parts[0], nil
}

// FetchUserPostContents 指定ユーザーの投稿内容を取得 (content のみ)
func (dao *GeminiDAO) FetchUserPostContents(userID string) ([]string, error) {
	rows, err := dao.db.Query(`
		SELECT content 
		FROM posts 
		WHERE user_id = ? AND deleted_at IS NULL 
		ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Printf("[gemini_dao.go] 以下の投稿一覧取得失敗 (user_id: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var contents []string
	for rows.Next() {
		var content string

		if err := rows.Scan(&content); err != nil {
			log.Printf("[gemini_dao.go] 投稿データのScan失敗: %v", err)
			return nil, err
		}

		contents = append(contents, content)
	}
	return contents, nil
}

// GetPostContent 指定した投稿IDの内容を文字列として取得
func (dao *GeminiDAO) GetPostContent(postID string) (string, error) {
	var content string

	// 投稿内容を取得
	err := dao.db.QueryRow(
		"SELECT content FROM posts WHERE post_id = ? AND deleted_at IS NULL",
		postID,
	).Scan(&content)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[gemini_dao.go] 投稿が見つからない (post_id: %s)", postID)
			return "", fmt.Errorf("投稿が存在しません")
		}
		log.Printf("[gemini_dao.go] 投稿内容の取得失敗 (post_id: %s): %v", postID, err)
		return "", err
	}

	return content, nil
}

// UpdateIsBad 指定した投稿の is_bad カラムを更新
func (dao *GeminiDAO) UpdateIsBad(postID string, isBad bool) error {
	_, err := dao.db.Exec(
		"UPDATE posts SET is_bad = ? WHERE post_id = ? AND deleted_at IS NULL",
		isBad,
		postID,
	)
	if err != nil {
		log.Printf("[gemini_dao.go] is_bad 更新失敗 (post_id: %s, is_bad: %v): %v", postID, isBad, err)
		return fmt.Errorf("is_bad の更新に失敗しました: %w", err)
	}

	return nil
}

// FetchUnfollowedUsers 指定ユーザーがフォローしていないユーザーのID、名前、自己紹介を取得
func (dao *GeminiDAO) FetchUnfollowedUsers(authID string) ([]model.User, error) {
	rows, err := dao.db.Query(`
		SELECT u.user_id, u.name, u.bio
		FROM users u
		WHERE u.user_id NOT IN (
			SELECT f.following_user_id
			FROM follows f
			WHERE f.user_id = ?
		) AND u.user_id != ?
		ORDER BY u.created_at DESC
	`, authID, authID)

	if err != nil {
		log.Printf("[gemini_dao.go] 未フォローのユーザー取得失敗 (auth_id: %s): %v", authID, err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var bio sql.NullString

		if err := rows.Scan(&user.UserID, &user.Name, &bio); err != nil {
			log.Printf("[gemini_dao.go] ユーザーデータのScan失敗: %v", err)
			return nil, err
		}

		// NULL値を処理
		user.Bio = nullableToPointer(bio)

		users = append(users, user)
	}
	return users, nil
}
