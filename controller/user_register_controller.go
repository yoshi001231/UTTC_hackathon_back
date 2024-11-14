package controller

import (
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"kaizen/usecase"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type RegisterUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RegisterUserController struct {
	registerUserUseCase *usecase.RegisterUserUseCase
}

func NewRegisterUserController(useCase *usecase.RegisterUserUseCase) *RegisterUserController {
	return &RegisterUserController{registerUserUseCase: useCase}
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("fail: json.NewDecoder, %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ID生成
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	res, _ := ulid.New(ms, entropy)
	userId := res.String()

	// ユーザー登録
	if _, err := c.registerUserUseCase.Execute(userId, req.Name, req.Age); err != nil {
		log.Printf("fail: register user, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// レスポンス
	resp := map[string]string{"id": userId}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
