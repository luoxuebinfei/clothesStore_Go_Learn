package DB

import (
	"database/sql"
	"fmt"
	"time"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
)

var (
	Db  *sql.DB
	err error
)

// User 用户数据库的字段定义
type User struct {
	ID        int
	Password  string
	Email     string
	Username  string
	CreatedAt string
	UpdateAt  string
}

//初始化运行,连接数据库
func init() {
	Db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/clothes_store")
	if err != nil {
		panic(err.Error())
	}

}

// AddNewUser 添加新用户，即用户注册
func AddNewUser(email string, password string, username string) (int64, error) {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updateAt := time.Now().Format("2006-01-02 15:04:05")
	//判断邮箱是否被注册
	rows, err := Db.Query(`select id from user where email=?`, email)
	defer rows.Close()
	if rows.Next() {
		return -2, nil
	}
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	//没有被注册则继续执行
	sqlStr := `insert into user(email,password ,name,created_at,update_at) values(?,?,?,?,?)`
	stmt, err := Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("预编译异常：", err)
		return -1, err
	}
	res, err := stmt.Exec(email, password, username, createdAt, updateAt)
	if err != nil {
		fmt.Println("执行出现异常：", err)
		return -1, err
	}
	userId, _ := res.LastInsertId() //返回注册时生成的id
	return userId, nil
}
