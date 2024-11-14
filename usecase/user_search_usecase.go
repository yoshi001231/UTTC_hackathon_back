package usecase

import (
	"kaizen/dao"
	"kaizen/model"
)

type FindUserByNameUseCase struct {
	userDAO *dao.UserDAO
}

func NewFindUserByNameUseCase(userDAO *dao.UserDAO) *FindUserByNameUseCase {
	return &FindUserByNameUseCase{userDAO}
}

func (uc *FindUserByNameUseCase) Execute(name string) ([]model.User, error) {
	return uc.userDAO.FindUserByName(name)
}
