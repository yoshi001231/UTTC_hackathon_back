// dao/init_dao.go

package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"twitter/dao/auth"
	"twitter/dao/follow"
	"twitter/dao/like"
	"twitter/dao/post"
	"twitter/dao/user"
)

var (
	dbInstance *sql.DB
	once       sync.Once

	usersDAOInstance     *user.UsersDAO
	authsDAOInstance     *auth.UsersDAO
	postsDAOInstance     *post.PostsDAO
	likesDAOInstance     *like.LikesDAO
	followersDAOInstance *follow.FollowersDAO
)

// InitDB データベース接続の初期化
func InitDB() *sql.DB {
	once.Do(func() {
		// 開発用ローカル接続情報
		user := "uttc_user"
		password := "uttc_password"
		host := "localhost:3306" // ローカルホスト
		database := "uttc_hackathon_local_db"

		// デプロイ用
		//user := os.Getenv("MYSQL_USER")
		//password := os.Getenv("MYSQL_PWD")
		//host := os.Getenv("MYSQL_HOST")
		//database := os.Getenv("MYSQL_DATABASE")

		if user == "" || password == "" || host == "" || database == "" {
			log.Fatal("環境変数MYSQL_USER, MYSQL_PWD, MYSQL_HOST, MYSQL_DATABASEが設定されていない")
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, database)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("データベース接続失敗: %v", err)
		}
		dbInstance = db
	})
	return dbInstance
}

// CloseDB データベース接続の終了
func CloseDB() {
	if dbInstance != nil {
		if err := dbInstance.Close(); err != nil {
			log.Fatal("データベース接続終了失敗")
		}
		log.Println("データベース接続終了")
	}
}

// GetUsersDAO ユーザーDAOのインスタンスを取得
func GetUsersDAO() *user.UsersDAO {
	if usersDAOInstance == nil {
		usersDAOInstance = user.NewUsersDAO(InitDB())
	}
	return usersDAOInstance
}

// GetAuthDAO 認証用ユーザーDAOのインスタンスを取得
func GetAuthDAO() *auth.UsersDAO {
	if authsDAOInstance == nil {
		authsDAOInstance = auth.NewUsersDAO(InitDB())
	}
	return authsDAOInstance
}

// GetPostsDAO 投稿DAOのインスタンスを取得
func GetPostsDAO() *post.PostsDAO {
	if postsDAOInstance == nil {
		postsDAOInstance = post.NewPostsDAO(InitDB())
	}
	return postsDAOInstance
}

// GetLikesDAO いいねDAOのインスタンスを取得
func GetLikesDAO() *like.LikesDAO {
	if likesDAOInstance == nil {
		likesDAOInstance = like.NewLikesDAO(InitDB())
	}
	return likesDAOInstance
}

// GetFollowersDAO フォローDAOのインスタンスを取得
func GetFollowersDAO() *follow.FollowersDAO {
	if followersDAOInstance == nil {
		followersDAOInstance = follow.NewFollowersDAO(InitDB())
	}
	return followersDAOInstance
}
