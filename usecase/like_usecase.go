package usecase

import (
	"twitter/dao"
	"twitter/model"
)

type LikeUseCase struct {
	LikeDAO *dao.LikeDAO
}

func NewLikeUseCase(LikeDAO *dao.LikeDAO) *LikeUseCase {
	return &LikeUseCase{LikeDAO: LikeDAO}
}

// AddLike 投稿にいいねを追加
func (uc *LikeUseCase) AddLike(userID, postID string) error {
	return uc.LikeDAO.AddLike(userID, postID)
}

// RemoveLike 投稿のいいねを削除
func (uc *LikeUseCase) RemoveLike(userID, postID string) error {
	return uc.LikeDAO.RemoveLike(userID, postID)
}

// GetUsersByPostID 投稿にいいねしたユーザー一覧を取得
func (uc *LikeUseCase) GetUsersByPostID(postID string) ([]model.User, error) {
	return uc.LikeDAO.GetUsersByPostID(postID)
}
