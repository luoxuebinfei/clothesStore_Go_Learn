package DB

import (
	"database/sql"
	"errors"
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

// QueryUser 查询用户，即用户登录
func QueryUser(email string, paw string) (int, []map[string]string) {
	sqlStr := fmt.Sprintf("select * from user where email='%s'", email)
	res, ok := Query(sqlStr)
	if ok {
		//fmt.Println(res[0]["name"])
		if res[0]["password"] != paw {
			//两者不等，密码错误
			return -1, nil
		} else {
			return 0, res
		}
	} else {
		//邮箱账号不存在
		//fmt.Println("没有找到数据")
		return -2, nil
	}
}

// ChangePaw 重置密码
func ChangePaw(email string, newPaw string) int {
	sqlStr := fmt.Sprintf("select * from user where email='%s'", email)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	res, ok := Query(sqlStr)
	if ok {
		//获取id
		id := res[0]["id"]
		sqlStr = fmt.Sprintf("update user set password='%s',update_at='%s' where id='%s'", newPaw, nowTime, id)
		Exec(sqlStr)
		return 0
	} else {
		//用户不存在
		return -1
	}
	return 0
}

// UpdatePaw 更新密码
func UpdatePaw(oldPass string, newPass string, uid string) error {
	sqlStr := fmt.Sprintf("select * from user where id='%s'", uid)
	res, ok := Query(sqlStr)
	if ok {
		if res[0]["password"] != oldPass {
			return errors.New("原密码错误")
		} else {
			nowTime := time.Now().Format("2006-01-02 15:04:05")
			sqlStr = fmt.Sprintf("update user set password = '%s',update_at='%s' where id='%s';", newPass, nowTime, uid)
			err := Exec(sqlStr)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		return errors.New("无此账号")
	}
}

// AdminLogin 管理员登录
func AdminLogin(email string, paw string) ([]map[string]string, error) {
	if len(email) == 0 || len(paw) == 0 {
		return nil, errors.New("参数不完整")
	}
	sqlStr := fmt.Sprintf("select a.id,a.username,a.email,a.password from admin_user a where a.email='%s'", email)
	res, ok := Query(sqlStr)
	if ok {
		if res[0]["password"] != paw {
			return nil, errors.New("密码错误")
		}
		return res, nil

	} else {
		return nil, errors.New("账号不存在")
	}
}

func AdminGetAllUser(u string) []map[string]string {
	var sqlStr string
	if u == "1" {
		//普通用户
		sqlStr = fmt.Sprintf("select * from user;")
	} else {
		//管理员用户
		sqlStr = fmt.Sprintf("select * from admin_user;")
	}
	res, ok := Query(sqlStr)
	if ok {
		return res
	} else {
		return nil
	}
}

// AdminUpdate 管理员更新用户信息
func AdminUpdate(id string, name string, email string, paw string, u string) error {
	nowTime := time.Now().Format("2006-01-02 15:04-05")
	var sqlStr string
	if u == "1" {
		sqlStr = fmt.Sprintf("update user set name='%s',email='%s',password='%s',update_at='%s' where id='%s'", name, email, paw, nowTime, id)
	} else {
		sqlStr = fmt.Sprintf("update admin_user set username='%s',email='%s',password='%s',update_time='%s' where id='%s'", name, email, paw, nowTime, id)
	}
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}
