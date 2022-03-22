// Package clothesStore_go_Learn 主文件
// Go版本1.15
package main

import (
	"clothesStore_go_Learn/DB"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(Cors()) //使用跨域处理
	r.GET("/get_email_code", func(c *gin.Context) {
		email := c.Query("email")
		SendEmail(email)
		fmt.Println(Code[email])
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
			"data": "",
		})
	})
	r.POST("/register", register)
	err := r.Run(":8100")
	if err != nil {
		return
	}
}

// Cors 跨域处理
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type,Authorization")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
			c.Set("content-type", "multipart/form-data")
			c.Set("content-type", "application/xml")
			c.Set("content-type", "application/x-www-form-urlencoded")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

//注册函数
func register(c *gin.Context) {
	dl := make(map[string]interface{})
	now_time := time.Now().Unix()
	err := c.BindJSON(&dl)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "无效参数",
			"data": "",
		})
		return
	}
	email := dl["email"].(string)
	paw := dl["password"].(string)
	username := dl["username"].(string)
	code, _ := strconv.Atoi(dl["code"].(string))
	if (code != Code[email].CodeNum) || (now_time-Code[email].CodeTime) > 600 {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "验证码错误或已过期，请重新获取输入！",
			"data": "",
		})
		return
	} else {
		paw_ := AesEncrypt(paw) //加密密码
		userID, err := DB.AddNewUser(email, paw_, username)
		if userID == -2 {
			c.JSON(200, gin.H{
				"code": 2001,
				"msg":  "此邮箱已被注册！",
				"data": "",
			})
			return
		}
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"code": 2001,
				"msg":  "注册失败！请重试！",
				"data": "",
			})
			return
		}
		token, err := GenToken(userID, username, email)
		if err != nil {
			return
		}
		delete(Code, email) //使用后销毁验证码
		c.JSON(200, gin.H{
			"code": http.StatusOK,
			"msg":  "注册成功！",
			"data": gin.H{
				"token": token,
			},
		})
		return
	}

}
