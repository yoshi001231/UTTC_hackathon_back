// usecase/like_usecase.go

package usecase

import (
	"twitter/dao/like"
	"twitter/model"
)

type LikeUseCase struct {
	likesDAO *like.LikesDAO
}

func NewLikeUseCase(likesDAO *like.LikesDAO) *LikeUseCase {
	return &LikeUseCase{likesDAO: likesDAO}
}

// AddLike 投稿にいいねを追加
func (uc *LikeUseCase) AddLike(userID, postID string) error {
	return uc.likesDAO.AddLike(userID, postID)
}

// RemoveLike 投稿のいいねを削除
func (uc *LikeUseCase) RemoveLike(userID, postID string) error {
	return uc.likesDAO.RemoveLike(userID, postID)
}

// GetUsersByPostID 投稿にいいねしたユーザー一覧を取得
func (uc *LikeUseCase) GetUsersByPostID(postID string) ([]model.User, error) {
	return uc.likesDAO.GetUsersByPostID(postID)
}
