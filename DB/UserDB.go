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

func ChangePaw(email string, new_paw string) int {
	sqlStr := fmt.Sprintf("select * from user where email='%s'", email)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	res, ok := Query(sqlStr)
	if ok {
		//获取id
		id := res[0]["id"]
		sqlStr = fmt.Sprintf("update user set password='%s',update_at='%s' where id='%s'", new_paw, nowTime, id)
		Exec(sqlStr)
		return 0
	} else {
		//用户不存在
		return -1
	}
	return 0
}

// Exec 增、删、改
func Exec(SQL string) {
	ret, err := Db.Exec(SQL) //增、删、改就靠这一条命令就够了，很简单
	if err != nil {
		fmt.Println(err)
	}
	_, _ = ret.LastInsertId()
	//fmt.Println(insID)
}

// Query 通用查询
func Query(SQL string) ([]map[string]string, bool) {
	rows, err := Db.Query(SQL) //执行SQL语句，比如select * from users
	if err != nil {
		panic(err)
	}
	columns, _ := rows.Columns()            //获取列的信息
	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i, _ := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}
	ret := make([]map[string]string, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		err := rows.Scan(values...)  //开始读行，Scan函数只接受指针变量
		m := make(map[string]string) //用于存放1列的 [键/值] 对
		if err != nil {
			panic(err)
		}
		for i, colName := range columns {
			var raw_value = *(values[i].(*interface{})) //读出raw数据，类型为byte
			b, _ := raw_value.([]byte)
			v := string(b) //将raw数据转换成字符串
			m[colName] = v //colName是键，v是值
		}
		ret = append(ret, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	defer rows.Close()

	if len(ret) != 0 {
		return ret, true
	}
	return nil, false
}
