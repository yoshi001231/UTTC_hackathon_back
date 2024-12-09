package usecase

import (
	"errors"
	"twitter/dao"
	"twitter/model"
)

type FollowUseCase struct {
	FollowDAO *dao.FollowDAO
}

func NewFollowUseCase(FollowDAO *dao.FollowDAO) *FollowUseCase {
	return &FollowUseCase{FollowDAO: FollowDAO}
}

// AddFollow 指定ユーザーをフォロー
func (uc *FollowUseCase) AddFollow(userID, followingUserID string) error {
	if userID == "" || followingUserID == "" {
		return errors.New("[follow_usecase.go] user_id または following_user_id が無効: 必須項目")
	}
	return uc.FollowDAO.AddFollow(userID, followingUserID)
}

// RemoveFollow 指定ユーザーのフォローを解除
func (uc *FollowUseCase) RemoveFollow(userID, followingUserID string) error {
	if userID == "" || followingUserID == "" {
		return errors.New("[follow_usecase.go] user_id または following_user_id が無効: 必須項目")
	}
	return uc.FollowDAO.RemoveFollow(userID, followingUserID)
}

// GetFollowers 指定ユーザーのフォロワー一覧を取得
func (uc *FollowUseCase) GetFollowers(userID string) ([]model.User, error) {
	return uc.FollowDAO.GetFollowers(userID)
}

// GetFollowing 指定ユーザーのフォロー中一覧を取得
func (uc *FollowUseCase) GetFollowing(userID string) ([]model.User, error) {
	return uc.FollowDAO.GetFollowing(userID)
}

// GetFollowGraph フォローグラフを取得
func (uc *FollowUseCase) GetFollowGraph() ([]model.Follow, error) {
	return uc.FollowDAO.GetFollowGraph()
}
