package controller

import (
	"encoding/json"
	"kaizen/usecase"
	"log"
	"net/http"
)

type SearchUserController struct {
	findUserUseCase *usecase.FindUserByNameUseCase
}

func NewSearchUserController(useCase *usecase.FindUserByNameUseCase) *SearchUserController {
	return &SearchUserController{findUserUseCase: useCase}
}

func (c *SearchUserController) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		log.Println("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := c.findUserUseCase.Execute(name)
	if err != nil {
		log.Printf("fail: find user, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("fail: jdon.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
