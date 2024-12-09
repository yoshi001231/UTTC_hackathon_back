package usecase

import (
	"errors"
	"twitter/dao"
	"twitter/model"
)

// FindUseCase 検索用のUseCase
type FindUseCase struct {
	FindDAO *dao.FindDAO
}

func NewFindUseCase(findDAO *dao.FindDAO) *FindUseCase {
	return &FindUseCase{FindDAO: findDAO}
}

// FindUsers 指定したキーワードを含むユーザーを検索
func (uc *FindUseCase) FindUsers(key string) ([]model.User, error) {
	if key == "" {
		return nil, errors.New("[find_usecase.go] キーワードが空です")
	}
	return uc.FindDAO.FindUsersByKey(key)
}

// FindPosts 指定したキーワードを含む投稿を検索
func (uc *FindUseCase) FindPosts(key string) ([]model.Post, error) {
	if key == "" {
		return nil, errors.New("[find_usecase.go] キーワードが空です")
	}
	return uc.FindDAO.FindPostsByKey(key)
}
