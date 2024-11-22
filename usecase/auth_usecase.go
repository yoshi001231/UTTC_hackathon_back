// usecase/auth_usecase.go

package usecase

import (
	"errors"
	"twitter/dao/auth"
	"twitter/model"
)

type RegisterUserUseCase struct {
	UsersDAO *auth.UsersDAO
}

func NewRegisterUserUseCase(UsersDAO *auth.UsersDAO) *RegisterUserUseCase {
	return &RegisterUserUseCase{UsersDAO: UsersDAO}
}

func (uc *RegisterUserUseCase) Execute(userID, name, bio, profileImgURL string) (string, error) {
	// バリデーション
	if userID == "" {
		return "", errors.New("user_id が無効: 必須項目")
	}
	if name == "" || len(name) > 50 {
		return "", errors.New("名前が無効: 必須項目で50文字以内である必要がある")
	}
	if len(bio) > 160 {
		return "", errors.New("自己紹介が無効: 160文字以内である必要がある")
	}

	// ユーザー登録
	user := model.User{
		UserID:        userID,
		Name:          name,
		Bio:           bio,
		ProfileImgURL: profileImgURL,
	}
	if err := uc.UsersDAO.RegisterUser(user); err != nil {
		return "", err
	}

	return user.UserID, nil
}
