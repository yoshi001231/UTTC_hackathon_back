package dao

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"database/sql"
	"fmt"
	"log"
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

// GenerateContentWithPastPosts Geminiを使用してツイート履歴から名前や自己紹介を生成
func (dao *GeminiDAO) GenerateContentWithPastPosts(prompt string) (*genai.Part, error) {
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
