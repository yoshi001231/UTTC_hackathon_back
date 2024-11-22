// usecase/user_usecase.go

package usecase

import (
	"errors"
	"twitter/dao/user"
	"twitter/model"
)

type GetUserUseCase struct {
	UsersDAO *user.UsersDAO
}

func NewGetUserUseCase(UsersDAO *user.UsersDAO) *GetUserUseCase {
	return &GetUserUseCase{UsersDAO: UsersDAO}
}

func (uc *GetUserUseCase) Execute(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user_id が無効です")
	}
	return uc.UsersDAO.GetUser(userID)
}

type UpdateProfileUseCase struct {
	UsersDAO *user.UsersDAO
}

func NewUpdateProfileUseCase(UsersDAO *user.UsersDAO) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{UsersDAO: UsersDAO}
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
	return uc.UsersDAO.UpdateUser(user)
}

// GetUser 更新後のユーザー情報を取得する
func (uc *UpdateProfileUseCase) GetUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user_id が無効です")
	}
	return uc.UsersDAO.GetUser(userID)
}
