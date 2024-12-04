package usecase

import (
	"errors"
	"twitter/dao"
	"twitter/model"
)

type TimelineUseCase struct {
	TimelineDAO *dao.TimelineDAO
}

func NewTimelineUseCase(TimelineDAO *dao.TimelineDAO) *TimelineUseCase {
	return &TimelineUseCase{TimelineDAO: TimelineDAO}
}

// GetUserTimeline ログインユーザーのタイムラインを取得
func (uc *TimelineUseCase) GetUserTimeline(userID string) ([]model.Post, error) {
	if userID == "" {
		return nil, errors.New("[timeline_usecase.go] auth_id が無効: 必須項目")
	}
	return uc.TimelineDAO.FetchUserTimeline(userID)
}

// GetUserPosts 指定ユーザーの投稿一覧を取得
func (uc *TimelineUseCase) GetUserPosts(userID string) ([]model.Post, error) {
	if userID == "" {
		return nil, errors.New("[timeline_usecase.go] auth_id が無効: 必須項目")
	}
	return uc.TimelineDAO.FetchUserPosts(userID)
}

// GetLikedPosts 指定ユーザーのいいねした投稿一覧を取得
func (uc *TimelineUseCase) GetLikedPosts(userID string) ([]model.Post, error) {
	if userID == "" {
		return nil, errors.New("[timeline_usecase.go] user_id が無効: 必須項目")
	}
	return uc.TimelineDAO.FetchLikedPosts(userID)
}
