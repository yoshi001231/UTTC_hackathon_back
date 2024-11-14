package usecase

import (
	"errors"
	"kaizen/dao"
	"kaizen/model"
)

type RegisterUserUseCase struct {
	userDAO dao.UserDAOInterface
}

func NewRegisterUserUseCase(userDAO dao.UserDAOInterface) *RegisterUserUseCase {
	return &RegisterUserUseCase{userDAO: userDAO}
}

func (uc *RegisterUserUseCase) Execute(id, name string, age int) (string, error) {
	// バリデーション
	if name == "" || len(name) > 50 {
		return "", errors.New("invalid name: must not be empty and max 50 characters")
	}
	if age < 20 || age > 80 {
		return "", errors.New("invalid age: must be between 20 and 80")
	}

	// ユーザー登録
	user := model.User{
		Id:   id,
		Name: name,
		Age:  age,
	}

	if err := uc.userDAO.RegisterUser(user); err != nil {
		return "", err
	}

	return user.Id, nil
}
