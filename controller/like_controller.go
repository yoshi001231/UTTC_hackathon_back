// controller/like_controller.go

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

type LikeController struct {
	likeUseCase *usecase.LikeUseCase
}

func NewLikeController(likeUseCase *usecase.LikeUseCase) *LikeController {
	return &LikeController{likeUseCase: likeUseCase}
}

// HandleAddLike 投稿にいいねを追加
func (c *LikeController) HandleAddLike(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	var like model.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		log.Printf("JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if like.UserID == "" {
		http.Error(w, "user_id が必要です", http.StatusBadRequest)
		return
	}

	if err := c.likeUseCase.AddLike(like.UserID, postID); err != nil {
		log.Printf("いいね追加失敗: %v", err)
		http.Error(w, "いいね追加に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HandleRemoveLike 投稿のいいねを削除
func (c *LikeController) HandleRemoveLike(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	var like model.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		log.Printf("JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if like.UserID == "" {
		http.Error(w, "user_id が必要です", http.StatusBadRequest)
		return
	}

	if err := c.likeUseCase.RemoveLike(like.UserID, postID); err != nil {
		log.Printf("いいね削除失敗: %v", err)
		http.Error(w, "いいね削除に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetUsersByPostID 投稿にいいねしたユーザー一覧を取得
func (c *LikeController) HandleGetUsersByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	users, err := c.likeUseCase.GetUsersByPostID(postID)
	if err != nil {
		log.Printf("いいねユーザー一覧取得失敗: %v", err)
		http.Error(w, "いいねユーザー一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		log.Printf("JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
