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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// DAO初期化
	authDAO := dao.GetAuthDAO()
	followDAO := dao.GetFollowDAO()
	likeDAO := dao.GetLikeDAO()
	postDAO := dao.GetPostDAO()
	timelineDAO := dao.GetTimelineDAO()
	userDAO := dao.GetUserDAO()
	findDAO := dao.GetFindDAO()
	// UseCase初期化
	authUseCase := usecase.NewAuthUseCase(authDAO)
	followUseCase := usecase.NewFollowUseCase(followDAO)
	likeUseCase := usecase.NewLikeUseCase(likeDAO)
	postUseCase := usecase.NewPostUseCase(postDAO)
	timelineUseCase := usecase.NewTimelineUseCase(timelineDAO)
	userUseCase := usecase.NewUserUseCase(userDAO)
	findUseCase := usecase.NewFindUseCase(findDAO)
	// Controller初期化
	authController := controller.NewAuthController(authUseCase)
	followController := controller.NewFollowController(followUseCase)
	likeController := controller.NewLikeController(likeUseCase)
	postController := controller.NewPostController(postUseCase)
	timelineController := controller.NewTimelineController(timelineUseCase)
	userController := controller.NewUserController(userUseCase)
	findController := controller.NewFindController(findUseCase)

	// ルーター初期化
	router := mux.NewRouter()

	// ユーザー関連エンドポイント
	router.HandleFunc("/auth/register", authController.Handle).Methods("POST")
	router.HandleFunc("/user/{user_id}", userController.HandleGetUser).Methods("GET")
	router.HandleFunc("/user/update-profile", userController.HandleUpdateProfile).Methods("PUT")
	// +ユーザランキング関連エンドポイント
	router.HandleFunc("/users/top/tweets", userController.HandleGetTopUsersByTweetCount).Methods("GET")
	router.HandleFunc("/users/top/likes", userController.HandleGetTopUsersByLikes).Methods("GET")

	// 投稿関連エンドポイント
	router.HandleFunc("/post/create", postController.HandleCreatePost).Methods("POST")
	router.HandleFunc("/post/{post_id}", postController.HandleGetPost).Methods("GET")
	router.HandleFunc("/post/{post_id}/update", postController.HandleUpdatePost).Methods("PUT")
	router.HandleFunc("/post/{post_id}/delete", postController.HandleDeletePost).Methods("DELETE")
	router.HandleFunc("/post/{post_id}/reply", postController.HandleReplyPost).Methods("POST")
	router.HandleFunc("/post/{post_id}/children", postController.HandleGetChildrenPosts).Methods("GET")

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
	router.HandleFunc("/timeline/{auth_id}", timelineController.HandleGetUserTimeline).Methods("GET")
	router.HandleFunc("/timeline/posts_by/{user_id}", timelineController.HandleGetUserPosts).Methods("GET")
	router.HandleFunc("/timeline/liked_by/{user_id}", timelineController.HandleGetLikedPosts).Methods("GET")

	// 検索関連エンドポイント
	router.HandleFunc("/find/user/{key}", findController.HandleFindUsers).Methods("GET")
	router.HandleFunc("/find/post/{key}", findController.HandleFindPosts).Methods("GET")

	// OPTIONSリクエストに対応
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// CORS設定
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                       // 許可するURL
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // 許可するヘッダー
		handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"}), // 許可するHTTPメソッド
	)

	// CORSエラーロギングミドルウェア
	loggingHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rec, r)
			if rec.statusCode == http.StatusForbidden || rec.statusCode == http.StatusMethodNotAllowed {
				log.Printf("[CORS Error] メソッド: %s, パス: %s, ステータスコード: %d", r.Method, r.URL.Path, rec.statusCode)
			}
		})
	}

	// シグナル処理
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		dao.CloseDB()
		os.Exit(0)
	}()

	// サーバー起動
	log.Println("[main.go] サーバー起動中...")
	wrappedRouter := loggingHandler(corsOptions(router))
	if err := http.ListenAndServe(":8080", wrappedRouter); err != nil {
		log.Fatal("[main.go] サーバー起動失敗")
	}
}

// responseRecorder レスポンスのステータスコードを記録するためのラッパー
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader ステータスコードを記録
func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
