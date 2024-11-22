// main.go

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

	"github.com/gorilla/mux"
)

func main() {
	// DAO初期化
	usersDAO := dao.GetUsersDAO()
	authDAO := dao.GetAuthDAO()
	postsDAO := dao.GetPostsDAO()
	likesDAO := dao.GetLikesDAO()
	followersDAO := dao.GetFollowersDAO()

	// UseCase, Controller初期化
	registerUserUseCase := usecase.NewRegisterUserUseCase(authDAO)
	registerUserController := controller.NewRegisterUserController(registerUserUseCase)
	getUserUseCase := usecase.NewGetUserUseCase(usersDAO)
	getUserController := controller.NewGetUserController(getUserUseCase)
	updateProfileUseCase := usecase.NewUpdateProfileUseCase(usersDAO)
	updateProfileController := controller.NewUpdateProfileController(updateProfileUseCase)
	postUseCase := usecase.NewPostUseCase(postsDAO)
	postController := controller.NewPostController(postUseCase)
	likeUseCase := usecase.NewLikeUseCase(likesDAO)
	likeController := controller.NewLikeController(likeUseCase)
	followUseCase := usecase.NewFollowUseCase(followersDAO)
	followController := controller.NewFollowController(followUseCase)

	// ルーター初期化
	router := mux.NewRouter()

	// ユーザー関連エンドポイント
	router.HandleFunc("/auth/register", registerUserController.Handle).Methods("POST")
	router.HandleFunc("/user/{user_id}", getUserController.Handle).Methods("GET")
	router.HandleFunc("/user/update-profile", updateProfileController.Handle).Methods("PUT")
	// 投稿関連エンドポイント
	router.HandleFunc("/post/create", postController.HandleCreatePost).Methods("POST")
	router.HandleFunc("/post/{post_id}", postController.HandleGetPost).Methods("GET")
	router.HandleFunc("/post/{post_id}/update", postController.HandleUpdatePost).Methods("PUT")
	router.HandleFunc("/post/{post_id}/delete", postController.HandleDeletePost).Methods("DELETE")
	router.HandleFunc("/post/{post_id}/reply", postController.HandleReplyPost).Methods("POST")
	// いいね関連エンドポイント
	router.HandleFunc("/like/{post_id}", likeController.HandleAddLike).Methods("POST")
	router.HandleFunc("/like/{post_id}/remove", likeController.HandleRemoveLike).Methods("DELETE")
	router.HandleFunc("/like/{post_id}/users", likeController.HandleGetUsersByPostID).Methods("GET")
	// フォロー関連エンドポイント
	router.HandleFunc("/follow/{user_id}", followController.HandleAddFollow).Methods("POST")
	router.HandleFunc("/follow/{user_id}/remove", followController.HandleRemoveFollow).Methods("DELETE")
	router.HandleFunc("/follow/{user_id}/followers", followController.HandleGetFollowers).Methods("GET")
	router.HandleFunc("/follow/{user_id}/following", followController.HandleGetFollowing).Methods("GET")

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
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("サーバー起動失敗")
	}
}
