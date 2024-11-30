package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"
)

var (
	dbInstance *sql.DB
	once       sync.Once

	authDAOInstance     *AuthDAO
	followDAOInstance   *FollowDAO
	likeDAOInstance     *LikeDAO
	postDAOInstance     *PostDAO
	timelineDAOInstance *TimelineDAO
	userDAOInstance     *UserDAO
)

func InitDB() *sql.DB {
	once.Do(func() {
		mode := os.Getenv("MODE")
		var dsn string
		if mode == "production" {
			// ローカル用
			user := "uttc_user"
			password := "uttc_password"
			host := "localhost:3306"
			database := "uttc_hackathon_local_db"
			dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, database)
			log.Println("[init_dao.go] モード: ローカル")
		} else {
			// 環境変数から設定を取得（デプロイ用）
			user := os.Getenv("MYSQL_USER")
			password := os.Getenv("MYSQL_PWD")
			host := os.Getenv("MYSQL_HOST")
			database := os.Getenv("MYSQL_DATABASE")
			dsn = fmt.Sprintf("%s:%s@%s/%s?parseTime=true", user, password, host, database)
			log.Println("[init_dao.go] モード: デプロイ")
		}

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("[init_dao.go] データベース接続失敗: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("[init_dao.go] データベースへのPing失敗: %v", err)
		}
		dbInstance = db
	})
	return dbInstance
}

func CloseDB() {
	if dbInstance != nil {
		if err := dbInstance.Close(); err != nil {
			log.Fatal("[init_dao.go] データベース接続終了失敗")
		}
		log.Println("[init_dao.go] データベース接続終了")
	}
}

func GetAuthDAO() *AuthDAO {
	if authDAOInstance == nil {
		authDAOInstance = NewAuthDAO(InitDB())
	}
	return authDAOInstance
}

func GetFollowDAO() *FollowDAO {
	if followDAOInstance == nil {
		followDAOInstance = NewFollowDAO(InitDB())
	}
	return followDAOInstance
}

func GetLikeDAO() *LikeDAO {
	if likeDAOInstance == nil {
		likeDAOInstance = NewLikeDAO(InitDB())
	}
	return likeDAOInstance
}

func GetPostDAO() *PostDAO {
	if postDAOInstance == nil {
		postDAOInstance = NewPostDAO(InitDB())
	}
	return postDAOInstance
}

func GetTimelineDAO() *TimelineDAO {
	if timelineDAOInstance == nil {
		timelineDAOInstance = NewTimelineDAO(InitDB())
	}
	return timelineDAOInstance
}

func GetUserDAO() *UserDAO {
	if userDAOInstance == nil {
		userDAOInstance = NewUserDAO(InitDB())
	}
	return userDAOInstance
}
