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

	posts, err := c.timelineUseCase.GetUserTimeline(authID)
	if err != nil {
		log.Printf("[timeline_controller.go] タイムライン取得失敗: %v", err)
		http.Error(w, "タイムライン取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		log.Printf("[timeline_controller.go] JSONエンコード失敗: %v", err)
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
		log.Printf("[timeline_controller.go] 投稿一覧取得失敗: %v", err)
		http.Error(w, "投稿一覧取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		log.Printf("[timeline_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
