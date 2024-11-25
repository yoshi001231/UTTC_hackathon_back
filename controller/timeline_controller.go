package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

type TimelineController struct {
	timelineUseCase *usecase.TimelineUseCase
}

func NewTimelineController(useCase *usecase.TimelineUseCase) *TimelineController {
	return &TimelineController{timelineUseCase: useCase}
}

// HandleGetUserTimeline ログインユーザーのタイムラインを取得
func (c *TimelineController) HandleGetUserTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authID := vars["auth_id"]

	if authID == "" {
		http.Error(w, "auth_id が必要です", http.StatusBadRequest)
		return
	}

	posts, err := c.timelineUseCase.GetUserTimeline(authID)
	if err != nil {
		log.Printf("タイムライン取得失敗: %v", err)
		http.Error(w, "タイムライン取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		log.Printf("JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// HandleGetUserPosts 指定ユーザーの投稿一覧を取得
func (c *TimelineController) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "user_id が必要です", http.StatusBadRequest)
		return
	}

	posts, err := c.timelineUseCase.GetUserPosts(userID)
	if err != nil {
		log.Printf("投稿一覧取得失敗: %v", err)
		http.Error(w, "投稿一覧取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		log.Printf("JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
