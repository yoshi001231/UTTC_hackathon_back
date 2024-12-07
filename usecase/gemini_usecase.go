package usecase

import (
	"cloud.google.com/go/vertexai/genai"
	"fmt"
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

	prompt := "以下のツイート内容と指示をもとに、Twitterの自己紹介文を日本語で150字以内で生成してください。'#'はつけないでください。"
	prompt += "\nツイート内容:\n" + strings.Join(tweets, "\n")
	if instruction != "" {
		prompt += fmt.Sprintf(" 指示: %s", instruction)
	}

	return uc.geminiDAO.GenerateContentWithPastPosts(prompt)
}

// GenerateName 過去ツイートと指示から名前を生成
func (uc *GeminiUseCase) GenerateName(authID, instruction string) (*genai.Part, error) {
	tweets, err := uc.geminiDAO.FetchUserPostContents(authID)
	if err != nil {
		return nil, fmt.Errorf("過去ツイートの取得失敗: %w", err)
	}

	prompt := "以下のツイート内容と指示をもとに、Twitterの名前を日本語で15字以内で1つだけ生成してください。"
	prompt += "\nツイート内容:\n" + strings.Join(tweets, "\n")
	if instruction != "" {
		prompt += fmt.Sprintf(" 指示: %s", instruction)
	}

	return uc.geminiDAO.GenerateContentWithPastPosts(prompt)
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
	if instruction != "" {
		prompt += fmt.Sprintf(" 指示: %s", instruction)
	}
	if tempText != "" {
		prompt += " 現在のツイートの続きを生成してください:"
		prompt += fmt.Sprintf("\n現在のツイート: %s", tempText)
	}

	return uc.geminiDAO.GenerateContentWithPastPosts(prompt)
}
