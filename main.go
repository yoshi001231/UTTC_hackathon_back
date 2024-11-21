package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"twitter/controller"
	"twitter/dao"
	"twitter/usecase"
)

func main() {
	// DAO初期化
	userDAO := dao.GetUserDAO()
	authDAO := dao.GetAuthDAO()

	// UseCase, Controller初期化
	registerUserUseCase := usecase.NewRegisterUserUseCase(authDAO)
	registerUserController := controller.NewRegisterUserController(registerUserUseCase)

	getUserUseCase := usecase.NewGetUserUseCase(userDAO)
	updateProfileUseCase := usecase.NewUpdateProfileUseCase(userDAO)

	getUserController := controller.NewGetUserController(getUserUseCase)
	updateProfileController := controller.NewUpdateProfileController(updateProfileUseCase)

	// エンドポイント設定
	http.HandleFunc("/auth/register", registerUserController.Handle)
	http.HandleFunc("/user/", getUserController.Handle)
	http.HandleFunc("/user/update-profile", updateProfileController.Handle)

	// シグナル処理
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		dao.CloseDB()
		os.Exit(0)
	}()

	// サーバー起動
	log.Println("サーバー起動中...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("サーバー起動失敗")
	}
}
