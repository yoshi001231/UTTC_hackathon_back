package main

import (
	"kaizen/controller"
	"kaizen/dao"
	"kaizen/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// DAO, UseCase, Controllerの初期化
	userDAO := dao.NewUserDAO()
	defer userDAO.CloseDB()

	findUserUseCase := usecase.NewFindUserByNameUseCase(userDAO)
	registerUserUseCase := usecase.NewRegisterUserUseCase(userDAO)
	searchUserController := controller.NewSearchUserController(findUserUseCase)
	registerUserController := controller.NewRegisterUserController(registerUserUseCase)

	// ルートの設定
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			searchUserController.Handle(w, r)
		} else if r.Method == http.MethodPost {
			registerUserController.Handle(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// シグナルハンドリング
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		userDAO.CloseDB()
		os.Exit(0)
	}()

	// サーバーの開始
	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
