package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

// GetUserController ユーザー詳細取得用コントローラ
type GetUserController struct {
	getUserUseCase *usecase.GetUserUseCase
}

// NewGetUserController コントローラの初期化
func NewGetUserController(useCase *usecase.GetUserUseCase) *GetUserController {
	return &GetUserController{getUserUseCase: useCase}
}

// Handle ユーザー詳細取得エンドポイントのハンドラ
func (c *GetUserController) Handle(w http.ResponseWriter, r *http.Request) {
	// gorilla/mux を使用して user_id を取得
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		log.Println("user_id がリクエストに含まれていません")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ユーザー詳細取得
	user, err := c.getUserUseCase.Execute(userID)
	if err != nil {
		log.Printf("ユーザー取得失敗: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// レスポンス生成
	bytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("JSONエンコード失敗: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// UpdateProfileController プロフィール更新用コントローラ
type UpdateProfileController struct {
	updateProfileUseCase *usecase.UpdateProfileUseCase
}

// NewUpdateProfileController コントローラの初期化
func NewUpdateProfileController(useCase *usecase.UpdateProfileUseCase) *UpdateProfileController {
	return &UpdateProfileController{updateProfileUseCase: useCase}
}

// Handle プロフィール更新エンドポイントのハンドラ
func (c *UpdateProfileController) Handle(w http.ResponseWriter, r *http.Request) {
	var req model.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		log.Println("user_id がリクエストに含まれていません")
		http.Error(w, "user_id が必須です", http.StatusBadRequest)
		return
	}

	// プロフィール更新
	if err := c.updateProfileUseCase.Execute(req); err != nil {
		log.Printf("プロフィール更新失敗: %v", err)
		http.Error(w, "プロフィール更新に失敗しました", http.StatusInternalServerError)
		return
	}

	// 更新後のユーザー情報を取得してレスポンス
	updatedUser, err := c.updateProfileUseCase.GetUser(req.UserID)
	if err != nil {
		log.Printf("更新後のユーザー取得失敗: %v", err)
		http.Error(w, "更新後のユーザー情報取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedUser)
	if err != nil {
		log.Printf("JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
