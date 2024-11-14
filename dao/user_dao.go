package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"kaizen/model"
)

// UserDAOInterface は UserDAOのインターフェース
type UserDAOInterface interface {
	RegisterUser(user model.User) error
	FindUserByName(name string) ([]model.User, error)
}

type UserDAO struct {
	db *sql.DB
}

var (
	instance *UserDAO
	once     sync.Once
)

// NewUserDAO データベース接続を初期化して UserDAO のインスタンスを返す
func NewUserDAO() *UserDAO {
	once.Do(func() {
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPwd := os.Getenv("MYSQL_PWD")
		mysqlHost := os.Getenv("MYSQL_HOST")
		mysqlDatabase := os.Getenv("MYSQL_DATABASE")

		if mysqlUser == "" || mysqlPwd == "" || mysqlHost == "" || mysqlDatabase == "" {
			log.Fatal("環境変数MYSQL_USER, MYSQL_PWD, MYSQL_HOST, MYSQL_DATABASEが設定されていません")
		}

		connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			log.Fatalf("fail: sql.Open, %v\n", err)
		}

		instance = &UserDAO{db: db}
	})
	return instance
}

// CloseDB データベース接続を閉じる
func (dao *UserDAO) CloseDB() {
	if dao.db != nil {
		if err := dao.db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
	}
}

// FindUserByName ユーザー検索
func (dao *UserDAO) FindUserByName(name string) ([]model.User, error) {
	rows, err := dao.db.Query("SELECT id, name, age FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// RegisterUser ユーザー登録
func (dao *UserDAO) RegisterUser(user model.User) error {
	_, err := dao.db.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", user.Id, user.Name, user.Age)
	return err
}
