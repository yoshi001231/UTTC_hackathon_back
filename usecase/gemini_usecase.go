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

	prompt := "以下のツイート内容を基に、Twitterの自己紹介文を日本語で150字以内で生成してください。'#'はつけないでください"
	if instruction != "" {
		prompt += fmt.Sprintf(" 指示: %s", instruction)
	}
	prompt += "\nツイート内容:\n" + strings.Join(tweets, "\n")

	return uc.geminiDAO.GenerateContentWithPastPosts(prompt)
}

// GenerateName 過去ツイートと指示から名前を生成
func (uc *GeminiUseCase) GenerateName(authID, instruction string) (*genai.Part, error) {
	tweets, err := uc.geminiDAO.FetchUserPostContents(authID)
	if err != nil {
		return nil, fmt.Errorf("過去ツイートの取得失敗: %w", err)
	}

	prompt := "以下のツイート内容を基に、Twitterの名前を日本語で15字以内で1つだけ生成してください。"
	if instruction != "" {
		prompt += fmt.Sprintf(" 指示: %s", instruction)
	}
	prompt += "\nツイート内容:\n" + strings.Join(tweets, "\n")

	return uc.geminiDAO.GenerateContentWithPastPosts(prompt)
}
