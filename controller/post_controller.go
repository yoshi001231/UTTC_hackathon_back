package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"

	"github.com/gorilla/mux"
)

type PostController struct {
	postUseCase *usecase.PostUseCase
}

func NewPostController(useCase *usecase.PostUseCase) *PostController {
	return &PostController{postUseCase: useCase}
}

// HandleCreatePost 新しい投稿を作成
func (c *PostController) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	var req model.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[post_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	createdPost, err := c.postUseCase.CreatePost(req)
	if err != nil {
		log.Printf("[post_controller.go] 投稿作成失敗: %v", err)
		http.Error(w, "投稿作成に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(createdPost)
	if err != nil {
		log.Printf("[post_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// HandleGetPost 投稿の詳細を取得
func (c *PostController) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	post, err := c.postUseCase.GetPost(postID)
	if err != nil {
		if err.Error() == "投稿が削除されています" {
			log.Printf("[post_controller.go] 投稿削除済み")
			http.Error(w, "投稿が削除されています", http.StatusGone)
		} else {
			log.Printf("[post_controller.go] 投稿取得失敗: %v", err)
			http.Error(w, "投稿が見つかりません", http.StatusNotFound)
		}
		return
	}

	resp, err := json.Marshal(post)
	if err != nil {
		log.Printf("[post_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// HandleUpdatePost 投稿内容を更新
func (c *PostController) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	var req model.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[post_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}
	req.PostID = postID

	if err := c.postUseCase.UpdatePost(req); err != nil {
		log.Printf("[post_controller.go] 投稿更新失敗: %v", err)
		http.Error(w, "投稿更新に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleDeletePost 投稿を削除
func (c *PostController) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	if err := c.postUseCase.DeletePost(postID); err != nil {
		log.Printf("[post_controller.go] 投稿削除失敗: %v", err)
		http.Error(w, "投稿削除に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleReplyPost 指定した投稿にリプライを追加
func (c *PostController) HandleReplyPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentPostID := vars["post_id"]

	var req model.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[post_controller.go] JSONデコード失敗: %v", err)
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	// ポインタ型に変換
	req.ParentPostID = &parentPostID

	replyPost, err := c.postUseCase.ReplyPost(req)
	if err != nil {
		log.Printf("[post_controller.go] リプライ投稿失敗: %v", err)
		http.Error(w, "リプライ投稿に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(replyPost)
	if err != nil {
		log.Printf("[post_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// HandleGetChildrenPosts 指定した投稿の子ポストを取得
func (c *PostController) HandleGetChildrenPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentPostID := vars["post_id"]

	posts, err := c.postUseCase.GetChildrenPosts(parentPostID)
	if err != nil {
		log.Printf("[post_controller.go] 子ポスト一覧取得失敗: %v", err)
		http.Error(w, "子ポスト一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		log.Printf("[post_controller.go] JSONエンコード失敗: %v", err)
		http.Error(w, "レスポンス生成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
