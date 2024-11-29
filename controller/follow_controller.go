// controller/follow_controller.go

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

type FollowController struct {
	followUseCase *usecase.FollowUseCase
}

func NewFollowController(followUseCase *usecase.FollowUseCase) *FollowController {
	return &FollowController{followUseCase: followUseCase}
}

// HandleAddFollow 指定ユーザーをフォロー
func (c *FollowController) HandleAddFollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followingUserID := vars["user_id"]

	var follow model.Follow
	if err := json.NewDecoder(r.Body).Decode(&follow); err != nil {
		log.Printf("[follow_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if err := c.followUseCase.AddFollow(follow.UserID, followingUserID); err != nil {
		log.Printf("[follow_controller.go] フォロー追加失敗: %v", err)
		http.Error(w, "フォロー追加に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HandleRemoveFollow 指定ユーザーのフォローを解除
func (c *FollowController) HandleRemoveFollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followingUserID := vars["user_id"]

	var follow model.Follow
	if err := json.NewDecoder(r.Body).Decode(&follow); err != nil {
		log.Printf("[follow_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	if err := c.followUseCase.RemoveFollow(follow.UserID, followingUserID); err != nil {
		log.Printf("[follow_controller.go] フォロー解除失敗: %v", err)
		http.Error(w, "フォロー解除に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetFollowers 指定ユーザーのフォロワー一覧を取得
func (c *FollowController) HandleGetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	users, err := c.followUseCase.GetFollowers(userID)
	if err != nil {
		log.Printf("[follow_controller.go] フォロワー一覧取得失敗: %v", err)
		http.Error(w, "フォロワー一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		log.Printf("[follow_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// HandleGetFollowing 指定ユーザーのフォロー中一覧を取得
func (c *FollowController) HandleGetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	users, err := c.followUseCase.GetFollowing(userID)
	if err != nil {
		log.Printf("[follow_controller.go] フォロー中一覧取得失敗: %v", err)
		http.Error(w, "フォロー中一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		log.Printf("[follow_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
