package usecase

import (
	"errors"
	"twitter/dao/user"
	"twitter/model"
)

type GetUserUseCase struct {
	userDAO *user.UserDAO
}

func NewGetUserUseCase(userDAO *user.UserDAO) *GetUserUseCase {
	return &GetUserUseCase{userDAO: userDAO}
}

func (uc *GetUserUseCase) Execute(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user_id が無効です")
	}
	return uc.userDAO.GetUser(userID)
}

type UpdateProfileUseCase struct {
	userDAO *user.UserDAO
}

func NewUpdateProfileUseCase(userDAO *user.UserDAO) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userDAO: userDAO}
}

// Execute プロフィールを更新する
func (uc *UpdateProfileUseCase) Execute(user model.User) error {
	if user.UserID == "" {
		return errors.New("user_id が無効です")
	}
	if user.Name == "" || len(user.Name) > 50 {
		return errors.New("名前が無効です")
	}
	if len(user.Bio) > 160 {
		return errors.New("自己紹介が無効です")
	}
	return uc.userDAO.UpdateUser(user)
}

// GetUser 更新後のユーザー情報を取得する
func (uc *UpdateProfileUseCase) GetUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user_id が無効です")
	}
	return uc.userDAO.GetUser(userID)
}
