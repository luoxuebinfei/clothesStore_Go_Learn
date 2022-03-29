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
	r.POST("/login", login)
	r.POST("/reset_Password", resetPassword)
	r.GET("/cart_index", JWTAuthMiddleware(), cartIndex)
	r.POST("/delete_cart_shop", JWTAuthMiddleware(), deleteCart)
	r.GET("get_address", JWTAuthMiddleware(), getADDress)
	r.POST("add_address", JWTAuthMiddleware(), addNewADDress)
	r.POST("delete_address", JWTAuthMiddleware(), deleteADDress)
	r.POST("update_address", JWTAuthMiddleware(), updateADDress)
	r.POST("change_addressDefault", JWTAuthMiddleware(), changeADDressDefault)
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
	nowTime := time.Now().Unix()
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
	if (code != Code[email].CodeNum) || (nowTime-Code[email].CodeTime) > 600 {
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
		c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
		c.JSON(200, gin.H{
			"code": http.StatusOK,
			"msg":  "注册成功！",
			"data": gin.H{
				"userinfo": username,
			},
		})
		return
	}

}

//登录函数
func login(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	if len(d) == 0 {
		c.JSON(200, gin.H{
			"code": 2002,
			"msg":  "请求参数不正确",
		})
		return
	}
	//println(d["email"].(string))
	paw := AesEncrypt(d["paw"].(string))
	ok, res := DB.QueryUser(d["email"].(string), paw)
	if ok == 0 {
		id, err := strconv.ParseInt(res[0]["id"], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		username := res[0]["name"]
		token, err := GenToken(id, username, res[0]["email"])
		if err != nil {
			fmt.Println("生成token失败", err)
		}
		c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "登录成功！",
			"data": gin.H{
				"userinfo": username,
			},
		})
	} else if ok == -1 {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "密码错误！请检查输入",
			"data": "",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "邮箱账号不存在，请检查输入或注册新账号",
			"data": "",
		})
	}
}

//重置密码函数
func resetPassword(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	nowTime := time.Now().Unix()
	email := d["email"].(string)
	code, _ := strconv.Atoi(d["code"].(string))
	newPaw := AesEncrypt(d["password"].(string))
	if Code[email].CodeNum != code || nowTime-Code[email].CodeTime > 600 {
		//验证码错误
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "验证码错误或已过期，请重新获取输入",
			"data": "",
		})
		return
	} else {
		ok := DB.ChangePaw(email, newPaw)
		if ok == -1 {
			c.JSON(200, gin.H{
				"code": 2001,
				"msg":  "用户邮箱不存在！",
				"data": "",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "重置密码成功！",
				"data": "data",
			})
			return
		}
	}

}

//获取购物车数据
func cartIndex(c *gin.Context) {
	uid, _ := c.Get("uid")
	data := DB.CartInfo(strconv.FormatInt(uid.(int64), 10))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
}

//删除购物车中商品
func deleteCart(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	uid, _ := c.Get("uid")
	for _, value := range d["data"].([]interface{}) {
		skuId := value.(map[string]interface{})["shopId"]
		ok := DB.DeleteCart(skuId.(string), strconv.FormatInt(uid.(int64), 10))
		if ok == -1 {
			c.JSON(200, gin.H{
				"code": 2001,
				"msg":  "删除商品出现错误",
			})
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "成功删除",
	})
}

//获取收货地址
func getADDress(c *gin.Context) {
	uid, _ := c.Get("uid")
	res := DB.GetAddress(strconv.FormatInt(uid.(int64), 10))
	c.JSON(200, gin.H{
		"cood": 200,
		"msg":  "",
		"data": res,
	})
	return
}

//新增收货地址
func addNewADDress(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	if len(d) == 0 || len(d) != 4 {
		c.JSON(200, gin.H{
			"code": 2001,
			"mag":  "请求参数不足",
		})
		return
	}
	uid, _ := c.Get("uid")
	name := d["name"].(string)
	phone := d["phone"].(string)
	address := d["address"].(string)
	isDefault := d["is_default"].(string)
	ok := DB.AddNewAddress(strconv.FormatInt(uid.(int64), 10), name, phone, address, isDefault)
	if ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "添加新地址成功",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 202,
			"msg":  "添加新地址失败",
		})
		return
	}
}

//删除收货地址
func deleteADDress(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	if len(d) == 0 || len(d) != 1 {
		c.JSON(200, gin.H{
			"code": 2001,
			"mag":  "请求参数不足",
		})
		return
	}
	uid, _ := c.Get("uid")
	id := d["id"].(string)
	ok := DB.DeleteAddress(id, strconv.FormatInt(uid.(int64), 10))
	if ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "删除地址成功",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "删除地址失败",
		})
		return
	}
}

//更新收货地址
func updateADDress(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	if len(d) == 0 || len(d) != 5 {
		c.JSON(200, gin.H{
			"code": 2001,
			"mag":  "请求参数不足",
		})
		return
	}
	uid, _ := c.Get("uid")
	id := d["id"].(string)
	name := d["name"].(string)
	phone := d["phone"].(string)
	address := d["address"].(string)
	isDefault := d["is_default"].(string)
	ok := DB.UpdateAddress(id, strconv.FormatInt(uid.(int64), 10), name, phone, address, isDefault)
	if ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "更新地址成功",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 202,
			"msg":  "更新地址失败",
		})
		return
	}
}

//更换默认收货地址
func changeADDressDefault(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	uid, _ := c.Get("uid")
	id := d["id"].(string)
	ok := DB.ChangeDefault(id, strconv.FormatInt(uid.(int64), 10))
	if ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "请求参数与数据库不匹配",
		})
	}

}
