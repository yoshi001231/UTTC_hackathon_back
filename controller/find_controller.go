package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

// FindController 検索用のコントローラー
type FindController struct {
	findUseCase *usecase.FindUseCase
}

func NewFindController(findUseCase *usecase.FindUseCase) *FindController {
	return &FindController{findUseCase: findUseCase}
}

// HandleFindUsers 指定したキーワードを含むユーザーを検索
func (c *FindController) HandleFindUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	users, err := c.findUseCase.FindUsers(key)
	if err != nil {
		log.Printf("[find_controller.go] ユーザー検索失敗 (key: %s): %v", key, err)
		http.Error(w, "ユーザー検索に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		log.Printf("[find_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// HandleFindPosts 指定したキーワードを含む投稿を検索
func (c *FindController) HandleFindPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	posts, err := c.findUseCase.FindPosts(key)
	if err != nil {
		log.Printf("[find_controller.go] 投稿検索失敗 (key: %s): %v", key, err)
		http.Error(w, "投稿検索に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		log.Printf("[find_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
