package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"
)

type AuthController struct {
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(useCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{AuthUseCase: useCase}
}

// Handle ユーザー登録ハンドラー
func (c *AuthController) Handle(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("[auth_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	// デリファレンスを追加
	bio := ""
	if user.Bio != nil {
		bio = *user.Bio
	}

	profileImgURL := ""
	if user.ProfileImgURL != nil {
		profileImgURL = *user.ProfileImgURL
	}

	// ユーザー登録
	if _, err := c.AuthUseCase.RegisterUser(user.UserID, user.Name, bio, profileImgURL); err != nil {
		log.Printf("[auth_controller.go] ユーザー登録失敗: %v", err)
		http.Error(w, "ユーザー登録に失敗しました", http.StatusInternalServerError)
		return
	}

	// 成功時はステータスコードのみを返す
	w.WriteHeader(http.StatusCreated)
}
