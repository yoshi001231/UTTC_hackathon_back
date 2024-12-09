package usecase

import (
	"errors"
	"github.com/oklog/ulid"
	"math/rand"
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
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	postID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	post.PostID = postID
	post.CreatedAt = time.Now()

	if post.ParentPostID == nil || *post.ParentPostID == "" { // 修正
		post.ParentPostID = nil
	}

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
	if post.ParentPostID == nil || *post.ParentPostID == "" { // 修正
		return nil, errors.New("[post_usecase.go] リプライ対象の投稿IDが指定されていません")
	}
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))            // 乱数生成器の作成
	replyID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String() // ULIDの生成
	post.PostID = replyID
	post.CreatedAt = time.Now()
	return uc.PostDAO.CreatePost(post)
}

// GetChildrenPosts 子ポストを取得
func (uc *PostUseCase) GetChildrenPosts(parentPostID string) ([]model.Post, error) {
	if parentPostID == "" {
		return nil, errors.New("[post_usecase.go] parent_post_id が無効: 必須項目")
	}
	return uc.PostDAO.GetChildrenPosts(parentPostID)
}
