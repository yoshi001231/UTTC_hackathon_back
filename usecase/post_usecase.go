// usecase/post_usecase.go

package usecase

import (
	"errors"
	"time"
	"twitter/dao"
	"twitter/model"
)

type PostUseCase struct {
	PostDAO *dao.PostDAO
}

func NewPostUseCase(PostDAO *dao.PostDAO) *PostUseCase {
	return &PostUseCase{PostDAO: PostDAO}
}

// CreatePost 新しい投稿を作成
func (uc *PostUseCase) CreatePost(post model.Post) (*model.Post, error) {
	if post.Content == "" {
		return nil, errors.New("投稿内容が空です")
	}
	post.CreatedAt = time.Now()
	return uc.PostDAO.CreatePost(post)
}

// GetPost 投稿の詳細を取得
func (uc *PostUseCase) GetPost(postID string) (*model.Post, error) {
	return uc.PostDAO.GetPost(postID)
}

// UpdatePost 投稿を更新
func (uc *PostUseCase) UpdatePost(post model.Post) error {
	return uc.PostDAO.UpdatePost(post)
}

// DeletePost 投稿を削除 (論理削除)
func (uc *PostUseCase) DeletePost(postID string) error {
	return uc.PostDAO.DeletePost(postID)
}

// ReplyPost 指定した投稿にリプライを追加
func (uc *PostUseCase) ReplyPost(post model.Post) (*model.Post, error) {
	if post.ParentPostID == "" {
		return nil, errors.New("リプライ対象の投稿IDが指定されていません")
	}
	post.CreatedAt = time.Now()
	return uc.PostDAO.CreatePost(post)
}
