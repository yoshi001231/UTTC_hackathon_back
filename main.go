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
	userDAO := dao.GetUserDAO()
	authDAO := dao.GetAuthDAO()
	postDAO := dao.GetPostDAO()
	likeDAO := dao.GetLikeDAO()
	followDAO := dao.GetFollowDAO()
	timelineDAO := dao.GetTimelineDAO()

	// UseCase, Controller初期化
	authUseCase := usecase.NewAuthUseCase(authDAO)
	userUseCase := usecase.NewUserUseCase(userDAO)
	postUseCase := usecase.NewPostUseCase(postDAO)
	likeUseCase := usecase.NewLikeUseCase(likeDAO)
	followUseCase := usecase.NewFollowUseCase(followDAO)
	timelineUseCase := usecase.NewTimelineUseCase(timelineDAO)

	authController := controller.NewAuthController(authUseCase)
	userController := controller.NewUserController(userUseCase)
	postController := controller.NewPostController(postUseCase)
	likeController := controller.NewLikeController(likeUseCase)
	followController := controller.NewFollowController(followUseCase)
	timelineController := controller.NewTimelineController(timelineUseCase) // タイムラインコントローラを追加

	// ルーター初期化
	router := mux.NewRouter()

	// ユーザー関連エンドポイント
	router.HandleFunc("/auth/register", authController.Handle).Methods("POST")
	router.HandleFunc("/user/{user_id}", userController.HandleGetUser).Methods("GET")
	router.HandleFunc("/user/update-profile", userController.HandleUpdateProfile).Methods("PUT")

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

	// タイムライン関連エンドポイント
	router.HandleFunc("/timeline", timelineController.HandleGetUserTimeline).Methods("GET")        // ログインユーザーのタイムライン取得
	router.HandleFunc("/timeline/{user_id}", timelineController.HandleGetUserPosts).Methods("GET") // 指定ユーザーの投稿一覧取得

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
