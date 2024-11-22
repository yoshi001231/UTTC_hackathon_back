// usecase/auth_usecase.go

package usecase

import (
	"errors"
	"twitter/dao"
	"twitter/model"
)

type AuthUseCase struct { // 修正: 名前をAuthUseCaseに変更
	AuthDAO *dao.AuthDAO
}

func NewAuthUseCase(AuthDAO *dao.AuthDAO) *AuthUseCase { // 修正: コンストラクタも変更
	return &AuthUseCase{AuthDAO: AuthDAO}
}

func (uc *AuthUseCase) RegisterUser(userID, name, bio, profileImgURL string) (string, error) { // 修正: メソッド名をわかりやすく変更
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
	if err := uc.AuthDAO.RegisterUser(user); err != nil {
		return "", err
	}

	return user.UserID, nil
}
