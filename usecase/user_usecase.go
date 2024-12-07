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
	if user.Name == "" || len(user.Name) > 255 {
		return errors.New("[user_usecase.go] 名前が無効: 必須項目で255文字以内である必要がある")
	}
	if user.Location != nil && len(*user.Location) > 100 {
		return errors.New("[user_usecase.go] 位置情報が無効: 100文字以内である必要がある")
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

// GetTopUsersByTweetCount ツイート数の多い順にユーザ一覧を取得
func (uc *UserUseCase) GetTopUsersByTweetCount(limit int) ([]model.User, error) {
	return uc.UserDAO.GetTopUsersByTweetCount(limit)
}

// GetTopUsersByLikes いいね数の多い順にユーザ一覧を取得
func (uc *UserUseCase) GetTopUsersByLikes(limit int) ([]model.User, error) {
	return uc.UserDAO.GetTopUsersByLikes(limit)
}
