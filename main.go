// Package clothesStore_go_Learn 主文件
// Go版本1.15
package main

import (
	"clothesStore_go_Learn/DB"
	"clothesStore_go_Learn/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(Cors()) //使用跨域处理
	r.GET("/", index)
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
	r.POST("/add_cart", JWTAuthMiddleware(), addCart)
	r.POST("/update_cart_num", JWTAuthMiddleware(), updateCartNum)
	r.POST("/change_cart_checked", JWTAuthMiddleware(), changeCheckedCart)
	r.POST("/delete_cart_shop", JWTAuthMiddleware(), deleteCart)
	r.GET("/get_address", JWTAuthMiddleware(), getADDress)
	r.POST("/add_address", JWTAuthMiddleware(), addNewADDress)
	r.POST("/delete_address", JWTAuthMiddleware(), deleteADDress)
	r.POST("/update_address", JWTAuthMiddleware(), updateADDress)
	r.POST("/change_addressDefault", JWTAuthMiddleware(), changeADDressDefault)
	r.POST("/get_order_client", JWTAuthMiddleware(), getOrderClient)
	r.POST("/order", JWTAuthMiddleware(), order)
	r.GET("/get_all_order", JWTAuthMiddleware(), getAllOrder)
	r.GET("/query_order", JWTAuthMiddleware(), queryOrder)
	r.GET("/delete_order", JWTAuthMiddleware(), deleteOrder)
	r.GET("/cancel_order", JWTAuthMiddleware(), cancelOrder)
	r.GET("/shopInfo/:id", getShopInfo)
	r.GET("/search", search_)
	r.GET("/Home", JWTAuthMiddleware(), home)
	r.POST("/update_pass", JWTAuthMiddleware(), updatePass)
	r.GET("/pay", tools.WebPageAlipay)
	r.GET("/return", tools.AliPayNotify)
	r.GET("/msg", Msg)

	//管理员路由
	ad := r.Group("/admin")
	ad.POST("/login", adminLogin)
	ad.GET("/shop", JWTAdmin(), adminShop)
	ad.POST("/shop_add_spu", JWTAdmin(), adminAddSpu)
	ad.POST("/shop_delete_spu", JWTAdmin(), adminDeleteSpu)
	ad.POST("/shop_update_spu", JWTAdmin(), adminUpdateSpu)
	ad.POST("/shop_add_sku", JWTAdmin(), adminAddSku)
	ad.POST("/shop_delete_sku", JWTAdmin(), adminDeleteSku)
	ad.POST("/shop_update_sku", JWTAdmin(), adminUpdateSku)
	ad.GET("/order", JWTAdmin(), adminOrder)
	ad.POST("/order_update_address", JWTAdmin(), adminUpdateAddress)
	ad.POST("/order_update_shipping", JWTAdmin(), adminUpdateShipping)
	ad.GET("/user", JWTAdmin(), adminUser)
	ad.POST("/user_update", JWTAdmin(), adminUpdateUser)
	ad.POST("/index_update", JWTAdmin(), adminUpdateIndex)

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

//获取首页内容
func index(c *gin.Context) {
	data := DB.GetIndex()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": gin.H{
			"imageItems": data[0]["imageItems"],
			"adImages":   data[0]["adImages"],
			"good_items": data[0]["good_items"],
		},
	})
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
		token, err := GenToken(userID, username, email, "普通用户")
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
		token, err := GenToken(id, username, res[0]["email"], "普通用户")
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

//添加到购物车
func addCart(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	skuId := d["skuId"].(string)
	buynum := d["buyNum"].(string)
	uid, _ := c.Get("uid")
	err := DB.AddCart(skuId, buynum, strconv.FormatInt(uid.(int64), 10))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err,
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "购物车添加成功",
		"data": "",
	})
	return
}

//更新购物车数量
func updateCartNum(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	skuId := d["skuId"].(string)
	num := d["num"].(string)
	uid, _ := c.Get("uid")
	err := DB.UpdateCartNum(skuId, num, strconv.FormatInt(uid.(int64), 10))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
		})
		return
	}
}

//更新购物车商品选中状态
func changeCheckedCart(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	skulist := d["data"].([]interface{})
	status := d["status"].(string)
	uid, _ := c.Get("uid")
	for _, i := range skulist {
		err := DB.ChangeChecked(i.(string), status, strconv.FormatInt(uid.(int64), 10))
		if err != nil {
			c.JSON(200, gin.H{
				"code": 2001,
				"msg":  err.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
	})
	return
}

//删除购物车中商品
func deleteCart(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	uid, _ := c.Get("uid")
	for _, value := range d["data"].([]interface{}) {
		skuId := value.(string)
		err := DB.DeleteCart(skuId, strconv.FormatInt(uid.(int64), 10))
		if err != nil {
			c.JSON(200, gin.H{
				"code": 2001,
				"msg":  err.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

//获取收货地址
func getADDress(c *gin.Context) {
	uid, _ := c.Get("uid")
	res := DB.GetAddress(strconv.FormatInt(uid.(int64), 10))
	c.JSON(200, gin.H{
		"code": 200,
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
			"msg":  "请求参数不足",
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
			"msg":  "请求参数不足",
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
			"msg":  "请求参数不足",
		})
		return
	}
	uid, _ := c.Get("uid")
	id := d["id"].(string)
	name := d["name"].(string)
	phone := d["phone"].(string)
	address := d["address"].(string)
	isDefault := d["is_default"].(string)
	err := DB.UpdateAddress(id, strconv.FormatInt(uid.(int64), 10), name, phone, address, isDefault)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新地址成功",
		"data": "",
	})
}

//更换默认收货地址
func changeADDressDefault(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	uid, _ := c.Get("uid")
	id := d["id"].(string)
	data, err := DB.ChangeDefault(id, strconv.FormatInt(uid.(int64), 10))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
	return
}

//进入订单提交页面
func getOrderClient(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	//fmt.Println(d["data"])
	switch t := d["data"].(type) {
	case []interface{}:
		var data []map[string]string
		for _, value := range t {
			m := make(map[string]string)
			res := DB.GetGoodsInfoAll(value.(map[string]interface{})["skuid"].(string))
			m["skuId"] = res[0]["id"]
			m["imgurl"] = res[0]["first_image"]
			m["price"] = res[0]["price"]
			m["shopname"] = res[0]["title"]
			m["props"] = res[0]["product_specs"]
			m["buyNum"] = value.(map[string]interface{})["buyNum"].(string)
			buyNum, _ := strconv.Atoi(value.(map[string]interface{})["buyNum"].(string))
			stockNum, _ := strconv.Atoi(res[0]["stock_num"])
			if stockNum < buyNum {
				c.JSON(200, gin.H{
					"code": 2001,
					"msg":  fmt.Sprintf("%s 商品库存不足", res[0]["id"]),
					"data": "",
				})
				return
			}
			data = append(data, m)
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
			"data": data,
		})
		return
	}

}

//提交订单
func order(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	//fmt.Println(d["data"].(map[string]interface{})["data"])
	uid, _ := c.Get("uid")
	err, data := DB.GenerateOrderSn(strconv.FormatInt(uid.(int64), 10), d)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	newUrl := fmt.Sprintf("/pay?orderSn=%s", data)
	c.Redirect(301, newUrl)
	//c.JSON(200, gin.H{
	//	"code": 200,
	//	"msg":  "订单提交成功，请尽快付款！",
	//	"data":data,
	//})
	//return
}

//获取所有订单
func getAllOrder(c *gin.Context) {
	code := c.Query("code")
	uid, _ := c.Get("uid")
	res := DB.GetAllOrder(code, strconv.FormatInt(uid.(int64), 10))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": res,
	})
	return
}

//查询订单
func queryOrder(c *gin.Context) {
	//d := make(map[string]interface{})
	//c.BindJSON(&d)
	//orderSn := d["orderSn"].(string)
	orderSn := c.Query("search")
	uid, _ := c.Get("uid")
	data, err := DB.QueryOrderSn(orderSn, strconv.FormatInt(uid.(int64), 10))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "",
			"data": data,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
}

//删除订单
func deleteOrder(c *gin.Context) {
	orderSn := c.Query("orderSn")
	uid, _ := c.Get("uid")
	msg := DB.DeleteOrderSn(orderSn, strconv.FormatInt(uid.(int64), 10))
	if msg == "订单删除成功" {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  msg,
		})
	} else {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  msg,
		})
	}
}

//取消订单
func cancelOrder(c *gin.Context) {
	orderSn := c.Query("orderSn")
	uid, _ := c.Get("uid")
	msg := DB.CancelOrder(orderSn, strconv.FormatInt(uid.(int64), 10))
	if msg == "订单取消成功" {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  msg,
		})
	} else {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  msg,
		})
	}
}

//获取商品页面
func getShopInfo(c *gin.Context) {
	skuId := c.Param("id")
	data, err := DB.GetShopInfo(skuId)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err,
			"data": data,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
	return
}

//搜索商品
func search_(c *gin.Context) {
	///search?keyword=&
	keyword := c.Query("keyword")
	queryBrand := c.Query("brand")
	querySize := c.Query("size")
	page := c.DefaultQuery("page", "1")
	res := DB.SearchShop(keyword, queryBrand, querySize, page)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": res,
	})
}

//进入个人中心页面
func home(c *gin.Context) {
	username, _ := c.Get("username")
	email, _ := c.Get("email")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": map[string]interface{}{"username": username, "email": email},
	})
}

//更新密码
func updatePass(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	oldPass := AesEncrypt(d["oldPass"].(string))
	newPass := AesEncrypt(d["newPass"].(string))
	uid, _ := c.Get("uid")
	err := DB.UpdatePaw(oldPass, newPass, strconv.FormatInt(uid.(int64), 10))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "密码更新成功！",
		"data": "",
	})
	return
}

//管理员

//管理员登录
func adminLogin(c *gin.Context) {
	d := make(map[string]interface{})
	c.BindJSON(&d)
	email := d["email"].(string)
	paw := d["paw"].(string)
	paw = AesEncrypt(paw) //加密
	res, err := DB.AdminLogin(email, paw)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	id, _ := strconv.ParseInt(res[0]["id"], 10, 64)
	token_, _ := GenToken(id, res[0]["username"], res[0]["email"], "超级管理员")
	c.Header("Authorization", fmt.Sprintf("Bearer %s", token_))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功！",
		"data": gin.H{
			"userinfo": res[0]["username"],
			"group":    "超级管理员",
		},
	})
	return
}

//管理员获取所有商品
func adminShop(c *gin.Context) {
	data := DB.AdminGetAllShop()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
	return
}

//管理员获取所有订单
func adminOrder(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	data := DB.AdminGetAllOrder(query)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
	return
}

//管理员获取所有用户
func adminUser(c *gin.Context) {
	u := c.DefaultQuery("u", "1")
	data := DB.AdminGetAllUser(u)
	c.JSON(200, gin.H{
		"code": 200,
		"mag":  "",
		"data": data,
	})
}

//管理员更新用户信息
func adminUpdateUser(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	id := d["id"].(string)
	name := d["name"].(string)
	email := d["email"].(string)
	password := d["password"].(string)
	fmt.Println(password)
	paw := AesEncrypt(password)
	fmt.Println(paw)
	u := d["u"].(string)
	err := DB.AdminUpdate(id, name, email, paw, u)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新用户信息成功",
	})
}

//管理员更新用户订单收货地址信息
func adminUpdateAddress(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	orderSn := d["order_sn"].(string)
	name := d["name"].(string)
	phone := d["phone"].(string)
	info := d["info"].(string)
	err := DB.AdminUpdateAddress(orderSn, name, phone, info)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "更新收货地址失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新收货地址成功",
	})
	return
}

//管理员发货
func adminUpdateShipping(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	orderSn := d["order_sn"].(string)
	sn := d["sn"].(string)
	err := DB.AdminUpdateShipping(orderSn, sn)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  "更新发货信息失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新发货信息成功",
	})
	return
}

//更新index内容
func adminUpdateIndex(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	imageItems := d["imageItems"].(string)
	adImages := d["adImages"].(string)
	goodItems := d["good_items"].(string)
	err := DB.UpdateIndex(imageItems, adImages, goodItems)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功",
	})
	return
}

//增加商品spu
func adminAddSpu(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	shop_name := d["shop_name"].(string)
	product_specs := d["product_specs"].(string)
	images := d["images"].(string)
	details_images := d["details_images"].(string)
	attribute_list := d["attribute_list"].(string)
	DB.AdminAddSpu(shop_name, product_specs, attribute_list, images, details_images)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功",
		"data": "",
	})
}

//删除商品spu
func adminDeleteSpu(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	spuID := d["spuID"].(string)
	err := DB.AdminDeleteSpu(spuID)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功",
	})

}

//更新spu
func adminUpdateSpu(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	spuID := d["spuID"].(string)
	shop_name := d["shop_name"].(string)
	product_specs := d["product_specs"].(string)
	images := d["images"].(string)
	details_images := d["details_images"].(string)
	attribute_list := d["attribute_list"].(string)
	err := DB.AdminUpdateSpu(spuID, shop_name, product_specs, images, details_images, attribute_list)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功",
	})

}

//管理员增加sku
func adminAddSku(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	spuID := d["spuID"].(string)
	goods_name := d["goods_name"].(string)
	price := d["price"].(string)
	title := d["title"].(string)
	product_specs := d["product_specs"].(string)
	stock_num := d["stock_num"].(string)
	first_image := d["first_image"].(string)
	weight := d["weight"].(string)
	status := d["status"].(bool)
	err := DB.AdminAddSku(spuID, goods_name, price, title, stock_num, product_specs, first_image, weight, status)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
	})
}

//管理员删除sku
func adminDeleteSku(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	skuID := d["skuID"].(string)
	err := DB.AdminDeleteSku(skuID)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功",
	})
}

//管理员更新sku
func adminUpdateSku(c *gin.Context) {
	var d map[string]interface{}
	c.BindJSON(&d)
	skuID := d["skuID"].(string)
	goods_name := d["goods_name"].(string)
	price := d["price"].(string)
	title := d["title"].(string)
	product_specs := d["product_specs"].(string)
	stock_num := d["stock_num"].(string)
	weight := d["weight"].(string)
	status := d["status"].(bool)
	err := DB.AdminUpdateSku(skuID, goods_name, price, title, stock_num, product_specs, weight, status)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
	})
}
