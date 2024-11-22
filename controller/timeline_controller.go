// controller/timeline_controller.go

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

type TimelineController struct {
	timelineUseCase *usecase.TimelineUseCase
}

func NewTimelineController(timelineUseCase *usecase.TimelineUseCase) *TimelineController {
	return &TimelineController{timelineUseCase: timelineUseCase}
}

// HandleGetUserTimeline ログインユーザーのタイムラインを取得
func (c *TimelineController) HandleGetUserTimeline(w http.ResponseWriter, r *http.Request) {
	var req model.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id が必要です", http.StatusBadRequest)
		return
	}

	posts, err := c.timelineUseCase.GetUserTimeline(req.UserID)
	if err != nil {
		log.Printf("タイムライン取得失敗: %v", err)
		http.Error(w, "タイムラインの取得に失敗しました", http.StatusInternalServerError)
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

	posts, err := c.timelineUseCase.GetUserPosts(userID)
	if err != nil {
		log.Printf("ユーザー投稿一覧取得失敗: %v", err)
		http.Error(w, "投稿一覧の取得に失敗しました", http.StatusInternalServerError)
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
