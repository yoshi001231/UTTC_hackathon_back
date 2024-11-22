// usecase/follow_usecase.go

package usecase

import (
	"twitter/dao/follow"
	"twitter/model"
)

type FollowUseCase struct {
	followersDAO *follow.FollowersDAO
}

func NewFollowUseCase(followersDAO *follow.FollowersDAO) *FollowUseCase {
	return &FollowUseCase{followersDAO: followersDAO}
}

// AddFollow 指定ユーザーをフォロー
func (uc *FollowUseCase) AddFollow(userID, followingUserID string) error {
	return uc.followersDAO.AddFollow(userID, followingUserID)
}

// RemoveFollow 指定ユーザーのフォローを解除
func (uc *FollowUseCase) RemoveFollow(userID, followingUserID string) error {
	return uc.followersDAO.RemoveFollow(userID, followingUserID)
}

// GetFollowers 指定ユーザーのフォロワー一覧を取得
func (uc *FollowUseCase) GetFollowers(userID string) ([]model.User, error) {
	return uc.followersDAO.GetFollowers(userID)
}

// GetFollowing 指定ユーザーのフォロー中一覧を取得
func (uc *FollowUseCase) GetFollowing(userID string) ([]model.User, error) {
	return uc.followersDAO.GetFollowing(userID)
}
