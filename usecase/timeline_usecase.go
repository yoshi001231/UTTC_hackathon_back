// usecase/timeline_usecase.go

package usecase

import (
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
	return uc.TimelineDAO.GetUserTimeline(userID)
}

// GetUserPosts 指定ユーザーの投稿一覧を取得
func (uc *TimelineUseCase) GetUserPosts(userID string) ([]model.Post, error) {
	return uc.TimelineDAO.GetUserPosts(userID)
}
