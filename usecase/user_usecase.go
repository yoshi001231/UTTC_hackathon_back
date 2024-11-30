package usecase

import (
	"errors"
	"twitter/dao"
	"twitter/model"
)

type UserUseCase struct {
	UserDAO *dao.UserDAO
}

func NewUserUseCase(UserDAO *dao.UserDAO) *UserUseCase {
	return &UserUseCase{UserDAO: UserDAO}
}

// GetUser ユーザー情報を取得する
func (uc *UserUseCase) GetUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("[user_usecase.go] user_id が無効: 必須項目")
	}
	return uc.UserDAO.GetUser(userID)
}

// UpdateProfile プロフィールを更新する
func (uc *UserUseCase) UpdateProfile(user model.User) error {
	if user.UserID == "" {
		return errors.New("[user_usecase.go] user_id が無効: 必須項目")
	}
	if user.Name == "" || len(user.Name) > 50 {
		return errors.New("[user_usecase.go] 名前が無効: 必須項目で50文字以内である必要がある")
	}
	if len(user.Bio) > 160 {
		return errors.New("[user_usecase.go] 自己紹介が無効: 160文字以内である必要がある")
	}
	return uc.UserDAO.UpdateUser(user)
}

// GetUpdatedUser 更新後のユーザー情報を取得する
func (uc *UserUseCase) GetUpdatedUser(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("[user_usecase.go] user_id が無効: 必須項目")
	}
	return uc.UserDAO.GetUser(userID)
}
