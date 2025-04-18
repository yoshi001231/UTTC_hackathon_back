package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"twitter/model"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

// UserController ユーザー関連エンドポイントのコントローラ
type UserController struct {
	userUseCase *usecase.UserUseCase
}

// NewUserController コントローラの初期化
func NewUserController(useCase *usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: useCase}
}

// HandleGetUser ユーザー詳細取得エンドポイントのハンドラ
func (c *UserController) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	user, err := c.userUseCase.GetUser(userID)
	if err != nil {
		log.Printf("[user_controller.go] ユーザー取得失敗 (user_id: %s): %v", userID, err)
		http.Error(w, "指定されたユーザーが見つかりません", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(user)
	if err != nil {
		log.Printf("[user_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// HandleUpdateProfile プロフィール更新エンドポイントのハンドラ
func (c *UserController) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req model.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[user_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if err := c.userUseCase.UpdateProfile(req); err != nil {
		log.Printf("[user_controller.go] プロフィール更新失敗 (user_id: %s): %v", req.UserID, err)
		http.Error(w, "プロフィール更新に失敗しました", http.StatusInternalServerError)
		return
	}

	updatedUser, err := c.userUseCase.GetUpdatedUser(req.UserID)
	if err != nil {
		log.Printf("[user_controller.go] 更新後のユーザー取得失敗 (user_id: %s): %v", req.UserID, err)
		http.Error(w, "更新後のユーザー情報取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedUser)
	if err != nil {
		log.Printf("[user_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// デフォルトのlimit値
const DefaultLimit = 100

// HandleGetTopUsersByTweetCount ツイート数の多い順にユーザ一覧を取得
func (c *UserController) HandleGetTopUsersByTweetCount(w http.ResponseWriter, r *http.Request) {
	limit := parseLimit(r.URL.Query().Get("limit"))
	users, err := c.userUseCase.GetTopUsersByTweetCount(limit)
	if err != nil {
		http.Error(w, "ユーザ一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// HandleGetTopUsersByLikes いいね数の多い順にユーザ一覧を取得
func (c *UserController) HandleGetTopUsersByLikes(w http.ResponseWriter, r *http.Request) {
	limit := parseLimit(r.URL.Query().Get("limit"))
	users, err := c.userUseCase.GetTopUsersByLikes(limit)
	if err != nil {
		http.Error(w, "ユーザ一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// parseLimit リクエストのlimit値を解析し、デフォルト値や無効値を処理
func parseLimit(limitStr string) int {
	if limitStr == "" {
		return DefaultLimit
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return DefaultLimit
	}
	return limit
}
