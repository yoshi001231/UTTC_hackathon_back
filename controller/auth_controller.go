// controller/auth_controller.go

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"twitter/model"
	"twitter/usecase"
)

type RegisterUserController struct {
	registerUserUseCase *usecase.RegisterUserUseCase
}

func NewRegisterUserController(useCase *usecase.RegisterUserUseCase) *RegisterUserController {
	return &RegisterUserController{registerUserUseCase: useCase}
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("JSONデコード失敗: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// user_id が送信されていない場合エラー
	if user.UserID == "" {
		log.Println("user_id がリクエストに含まれていない")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ユーザー登録
	if _, err := c.registerUserUseCase.Execute(user.UserID, user.Name, user.Bio, user.ProfileImgURL); err != nil {
		log.Printf("ユーザー登録失敗: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// レスポンス生成
	resp := map[string]string{"user_id": user.UserID}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("JSONエンコード失敗: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
