package usecase

import (
	"cloud.google.com/go/vertexai/genai"
	"fmt"
	"log"
	"strings"
	"twitter/dao"
)

type GeminiUseCase struct {
	geminiDAO *dao.GeminiDAO
}

func NewGeminiUseCase(geminiDAO *dao.GeminiDAO) *GeminiUseCase {
	return &GeminiUseCase{geminiDAO: geminiDAO}
}

// GenerateBio 過去ツイートと指示から自己紹介を生成
func (uc *GeminiUseCase) GenerateBio(authID, instruction string) (*genai.Part, error) {
	tweets, err := uc.geminiDAO.FetchUserPostContents(authID)
	if err != nil {
		return nil, fmt.Errorf("過去ツイートの取得失敗: %w", err)
	}

	prompt := "以下のツイート内容をもとに、Twitterの自己紹介文を日本語で150字以内で生成してください。'#'はつけないでください。"
	prompt += "\nツイート内容:\n" + strings.Join(tweets, "\n")
	if instruction != "" {
		prompt += fmt.Sprintf(" 追加の指示: %s", instruction)
	}

	return uc.geminiDAO.GenerateResponseFromPrompt(prompt)
}

// GenerateName 過去ツイートと指示から名前を生成
func (uc *GeminiUseCase) GenerateName(authID, instruction string) (*genai.Part, error) {
	tweets, err := uc.geminiDAO.FetchUserPostContents(authID)
	if err != nil {
		return nil, fmt.Errorf("過去ツイートの取得失敗: %w", err)
	}

	prompt := "以下のツイート内容をもとに、Twitterの名前を日本語で15字以内で1つだけ生成してください。"
	prompt += "\nツイート内容:\n" + strings.Join(tweets, "\n")
	if instruction != "" {
		prompt += fmt.Sprintf(" 追加の指示: %s", instruction)
	}

	return uc.geminiDAO.GenerateResponseFromPrompt(prompt)
}

// GenerateTweetContinuation 過去ツイート、指示、現在の入力からツイートの続きを生成
func (uc *GeminiUseCase) GenerateTweetContinuation(authID, instruction, tempText string) (*genai.Part, error) {
	tweets, err := uc.geminiDAO.FetchUserPostContents(authID)
	if err != nil {
		return nil, fmt.Errorf("過去ツイートの取得失敗: %w", err)
	}

	// プロンプト作成
	prompt := "以下のツイート内容を基に、Twitterの新しいツイートを合計200字以内で生成してください。'#'はつけないでください。"
	prompt += "\n過去のツイート内容:\n" + strings.Join(tweets, "\n")
	prompt += "過去のツイートがタメ口中心ならタメ口中心、敬語中心なら敬語中心にしてください。"
	if instruction != "" {
		prompt += fmt.Sprintf(" 追加の指示: %s", instruction)
	}
	if tempText != "" {
		prompt += " 現在のツイートの続きを生成してください:"
		prompt += fmt.Sprintf("\n現在のツイート: %s", tempText)
	}

	return uc.geminiDAO.GenerateResponseFromPrompt(prompt)
}

// CheckIfPostIsBad 指定した投稿の内容を検査して Gemini の結果を返す
func (uc *GeminiUseCase) CheckIfPostIsBad(postID string) (*genai.Part, error) {
	// DAO から投稿内容を取得
	content, err := uc.geminiDAO.GetPostContent(postID)
	if err != nil {
		return nil, fmt.Errorf("投稿内容の取得失敗: %w", err)
	}

	// プロンプトを作成
	prompt := fmt.Sprintf("次の投稿が良識に反している場合は 'YES' を、そうでない場合は 'NO' を返してください:\n\n投稿内容: %s", content)

	// Gemini API を使用して判定
	return uc.geminiDAO.GenerateResponseFromPrompt(prompt)
}

// UpdateIsBad 指定した投稿の is_bad カラムを更新
func (uc *GeminiUseCase) UpdateIsBad(postID string, isBad bool) error {
	return uc.geminiDAO.UpdateIsBad(postID, isBad)
}

// RecommendUsers 指示からおすすめユーザーを生成
func (uc *GeminiUseCase) RecommendUsers(authID, instruction string) (*genai.Part, error) {
	// 未フォローのユーザー情報を取得
	unfollowedUsers, err := uc.geminiDAO.FetchUnfollowedUsers(authID)
	if err != nil {
		return nil, fmt.Errorf("未フォローのユーザー取得失敗: %w", err)
	}

	// 未フォローのユーザーをプロンプトに追加
	userInfo := ""
	for _, user := range unfollowedUsers {
		userInfo += fmt.Sprintf("- ID: %s, 名前: %s, 自己紹介: %s\n", user.UserID, user.Name, nullableToString(user.Bio))
	}

	// プロンプト作成
	prompt := "ユーザー情報をもとに、フォローすべきおすすめのユーザーを推薦してください。\n"
	prompt += "ユーザー情報:\n" + userInfo
	if instruction != "" {
		prompt += fmt.Sprintf(" 追加の指示: %s\n", instruction)
	}
	prompt += "必ず上記IDのいずれかから選択するようにしてください。\n"

	log.Printf("prompt", prompt, instruction)

	// Gemini APIで生成
	responsePart, err := uc.geminiDAO.GenerateResponseFromPrompt(prompt)
	if err != nil {
		return nil, fmt.Errorf("Geminiによる推薦生成失敗: %w", err)
	}

	return responsePart, nil
}

// nullableToString ヘルパー関数: *string を文字列に変換
func nullableToString(s *string) string {
	if s != nil {
		return *s
	}
	return "N/A"
}
