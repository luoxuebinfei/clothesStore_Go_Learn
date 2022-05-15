package DB

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

//订单表数据库操作
//包括生成订单，查询订单等操作

// GenerateOrderSn 提交订单，生成订单号
func GenerateOrderSn(uid string, data map[string]interface{}) (error, string) {
	//20220401014810 + 4位随机码 + id
	nowTime := time.Now().Format("20060102150405")
	rand.Seed(time.Now().Unix())
	randCode := rand.Intn(10000)
	orderSn := fmt.Sprintf("%s%d%s", nowTime, randCode, uid)
	//fmt.Println(orderSn)
	a := data["data"].([]interface{})
	x := make(map[string][]string) //商品id数组{商品id:[商品数量,单价]}

	for _, value := range a {
		//skuId := strconv.FormatFloat(value.(map[string]interface{})["skuId"].(float64), 'E', -1, 64)
		skuId := value.(map[string]interface{})["skuId"].(string)
		buyNum, _ := strconv.Atoi(value.(map[string]interface{})["buyNum"].(string))
		if buyNum <= 0 {
			return errors.New("参数不正确！"), ""
		}
		//先查询库存，如果库存不足则返回
		sqlStr := fmt.Sprintf("select goods_name,stock_num,price from sku where id='%s'", skuId)
		res, ok := Query(sqlStr)
		if ok {
			stockNum, _ := strconv.Atoi(res[0]["stock_num"])
			price := res[0]["price"]
			if buyNum <= stockNum {
				//购买数量小于等于库存数量时更新库存
				//fmt.Println("1111")
				sqlStr = fmt.Sprintf("update sku set stock_num='%d' where id='%s'", stockNum-buyNum, skuId)
				err := Exec(sqlStr)
				if err != nil {
					return err, ""
				}
				x[skuId] = []string{strconv.Itoa(buyNum), price}
			} else {
				//购买数量大于库存时，返回
				//fmt.Println("商品库存不足")
				return errors.New(fmt.Sprintf("%s 此商品库存不足", res[0]["goods_name"])), ""
			}
		} else {
			//fmt.Println("商品不存在")
			return errors.New(fmt.Sprintf("商品id:%s 此商品不存在", skuId)), ""
		}
	}

	//获取收货地址信息
	b := data["area"].(map[string]interface{})
	addressName := b["name"].(string)
	phone := b["phone"].(string)
	address := b["address"].(string)
	//fmt.Println(addressName,phone,address)
	nowTime = time.Now().Format("2006-01-02 15:04:05")
	total := decimal.NewFromFloat(0.00)
	for key, value := range x {
		buyNum, _ := strconv.ParseFloat(value[0], 64)
		price, _ := strconv.ParseFloat(value[1], 64)
		total = total.Add(decimal.NewFromFloat(buyNum).Mul(decimal.NewFromFloat(price)))
		sqlStr := fmt.Sprintf("insert into order_detail (order_sn, sku_id, price, num, created_time, update_time) value ('%s','%s','%s','%s','%s','%s')", orderSn, key, value[1], value[0], nowTime, nowTime)
		err := Exec(sqlStr)
		if err != nil {
			return err, ""
		}
		//删除购物车中商品
		DeleteCart(key, uid)
	}
	total_ := fmt.Sprintf("%.2f", total.InexactFloat64())
	//将订单号存入订单主表
	sqlStr := fmt.Sprintf("insert into order_master (order_sn, uid, address_name, address_phone, address_info, created_time, order_money, pay_money, order_status, update_time) value ('%s','%s','%s','%s','%s','%s','%s','%s',2,'%s')", orderSn, uid, addressName, phone, address, nowTime, total_, total_, nowTime)
	err := Exec(sqlStr)
	if err != nil {
		return err, ""
	}
	return nil, orderSn
}

// QueryOrderSn 查询订单
func QueryOrderSn(orderSn string, uid string) ([]map[string]interface{}, error) {
	var sqlStr string
	re1, _ := regexp.Compile(`^\d+$`)                        // 判断为订单号
	re2, _ := regexp.Compile("^[a-zA-Z0-9\u4E00-\u9FFF]+?$") //匹配中文数字英文
	if re1.MatchString(orderSn) && re2.MatchString(orderSn) {
		sqlStr = fmt.Sprintf("select a.created_time,a.order_sn,a.order_status,a.order_money,a.pay_money,a.address_name,a.address_phone,a.address_info,b.sku_id,c.title,c.first_image from (order_master a inner join order_detail b on a.order_sn=b.order_sn and a.uid='%s' and a.order_status!=0) inner join sku c on b.sku_id=c.id and (a.order_sn regexp '%s' or c.title regexp '%s')", uid, orderSn, orderSn)
	} else if re1.MatchString(orderSn) {
		sqlStr = fmt.Sprintf("select a.created_time,a.order_sn,a.order_status,a.order_money,a.pay_money,a.address_name,a.address_phone,a.address_info,b.sku_id,c.title,c.first_image from (order_master a inner join order_detail b on a.order_sn=b.order_sn and a.order_sn regexp '%s' and a.uid='%s' and a.order_status!=0) inner join sku c on b.sku_id=c.id", orderSn, uid)
	} else if re2.MatchString(orderSn) {
		sqlStr = fmt.Sprintf("select a.created_time,a.order_sn,a.order_status,a.order_money,a.pay_money,a.address_name,a.address_phone,a.address_info,b.sku_id,c.title,c.first_image from ((order_master a inner join order_detail b on a.order_sn=b.order_sn and a.uid='%s' and a.order_status!=0) inner join sku c on b.sku_id=c.id and c.title regexp '%s') inner join order_detail d on a.order_sn=d.order_sn", uid, orderSn)
	}
	res, ok := Query(sqlStr)
	if ok {
		fmt.Println(res)
		var x map[string]interface{}
		var data []map[string]interface{}
		for _, i := range res {
			if x["order_sn"] == i["order_sn"] {
				x["id"] = append(x["id"].([]string), i["sku_id"])
				x["name"] = append(x["name"].([]string), i["title"])
				x["imgurls"] = append(x["imgurls"].([]string), i["first_image"])
				x["shopnums"] = append(x["shopnums"].([]string), i["num"])
				x["unitprice"] = append(x["unitprice"].([]string), i["price"])
			} else if x["order_sn"] != i["order_sn"] && len(x) == 0 {
				x = make(map[string]interface{})
				x["order_sn"] = i["order_sn"]
				x["id"] = []string{i["sku_id"]}
				x["name"] = []string{i["title"]}
				x["date"] = i["created_time"]
				x["price"] = i["pay_money"]
				x["imgurls"] = []string{i["first_image"]}
				x["shopnums"] = []string{i["num"]}
				x["unitprice"] = []string{i["price"]}
				//订单状态
				if i["order_status"] == "1" {
					x["status"] = "已完成"
				} else if i["order_status"] == "2" {
					x["status"] = "未付款"
				} else if i["order_status"] == "3" {
					x["status"] = "待发货"
				} else if i["order_status"] == "4" {
					x["status"] = "已发货"
				} else if i["order_status"] == "5" {
					x["status"] = "已取消"
				}
				//收货信息
				z := make(map[string]string)
				z["name"] = i["address_name"]
				z["phone"] = i["address_phone"]
				z["info"] = i["address_info"]
				x["address"] = z
			} else if x["order_sn"] != i["order_sn"] && len(x) != 0 {
				//下一个
				data = append(data, x)
				x = nil
				x = make(map[string]interface{})
				x["order_sn"] = i["order_sn"]
				x["id"] = []string{i["sku_id"]}
				x["name"] = []string{i["title"]}
				x["date"] = i["created_time"]
				x["price"] = i["pay_money"]
				x["imgurls"] = []string{i["first_image"]}
				x["shopnums"] = []string{i["num"]}
				x["unitprice"] = []string{i["price"]}
				if i["order_status"] == "1" {
					x["status"] = "已完成"
				} else if i["order_status"] == "2" {
					x["status"] = "未付款"
				} else if i["order_status"] == "3" {
					x["status"] = "待发货"
				} else if i["order_status"] == "4" {
					x["status"] = "已发货"
				} else if i["order_status"] == "5" {
					x["status"] = "已取消"
				}
				//收货信息
				z := make(map[string]string)
				z["name"] = i["address_name"]
				z["phone"] = i["address_phone"]
				z["info"] = i["address_info"]
				x["address"] = z
			}
		}
		data = append(data, x)
		return data, nil
	} else {
		return nil, errors.New("搜索数据不存在")
	}
}

// DeleteOrderSn 删除订单，屏蔽订单
func DeleteOrderSn(orderSn string, uid string) string {
	sqlStr := fmt.Sprintf("update order_master set order_status = 0 where order_sn='%s' and uid='%s';", orderSn, uid)
	err := Exec(sqlStr)
	if err != nil {
		return "订单删除失败"
	}
	return "订单删除成功"
}

// CancelOrder 取消订单
func CancelOrder(orderSn string, uid string) string {
	sqlStr := fmt.Sprintf("update order_master set order_status = 5 where order_sn='%s' and uid='%s';", orderSn, uid)
	err := Exec(sqlStr)
	if err != nil {
		return "订单取消失败"
	}
	return "订单取消成功"
}

// GetAllOrder 获取当前账号所有订单
func GetAllOrder(code string, uid string) []map[string]interface{} {
	var sqlStr string
	if code == "All" {
		sqlStr = fmt.Sprintf("select a.created_time,a.order_sn,a.order_status,a.order_money,a.pay_money,a.address_name,a.address_phone,a.address_info,b.sku_id,b.num,b.price,c.title,c.first_image from (order_master a inner join order_detail b on a.uid='%s' and a.order_status!=0 and a.order_sn=b.order_sn) inner join sku c on b.sku_id=c.id;", uid)
	} else if code == "3" || code == "4" {
		sqlStr = fmt.Sprintf("select a.created_time,a.order_sn,a.order_status,a.order_money,a.pay_money,a.address_name,a.address_phone,a.address_info,b.sku_id,b.num,b.price,c.title,c.first_image from (order_master a inner join order_detail b on a.uid='%s' and a.order_sn=b.order_sn and a.order_status in (3,4)) inner join sku c on b.sku_id=c.id;", uid)
	} else if code != "0" {
		sqlStr = fmt.Sprintf("select a.created_time,a.order_sn,a.order_status,a.order_money,a.pay_money,a.address_name,a.address_phone,a.address_info,b.sku_id,b.num,b.price,c.title,c.first_image from (order_master a inner join order_detail b on a.uid='%s' and a.order_status='%s' and a.order_sn=b.order_sn) inner join sku c on b.sku_id=c.id;", uid, code)
	} else {
		return nil
	}
	res, ok := Query(sqlStr)
	if ok {
		var x map[string]interface{}
		var data []map[string]interface{}
		//fmt.Println(len(res))
		for _, i := range res {
			if x["order_sn"] == i["order_sn"] {
				x["id"] = append(x["id"].([]string), i["sku_id"])
				x["name"] = append(x["name"].([]string), i["title"])
				x["imgurls"] = append(x["imgurls"].([]string), i["first_image"])
				x["shopnums"] = append(x["shopnums"].([]string), i["num"])
				x["unitprice"] = append(x["unitprice"].([]string), i["price"])
			} else if x["order_sn"] != i["order_sn"] && len(x) == 0 {
				x = make(map[string]interface{})
				x["order_sn"] = i["order_sn"]
				x["id"] = []string{i["sku_id"]}
				x["name"] = []string{i["title"]}
				x["date"] = i["created_time"]
				x["price"] = i["pay_money"]
				x["imgurls"] = []string{i["first_image"]}
				x["shopnums"] = []string{i["num"]} //购买数量
				x["unitprice"] = []string{i["price"]}
				//订单状态
				if i["order_status"] == "1" {
					x["status"] = "已完成"
				} else if i["order_status"] == "2" {
					x["status"] = "未付款"
				} else if i["order_status"] == "3" {
					x["status"] = "待发货"
				} else if i["order_status"] == "4" {
					x["status"] = "已发货"
				} else if i["order_status"] == "5" {
					x["status"] = "已取消"
				}
				//收货信息
				z := make(map[string]string)
				z["name"] = i["address_name"]
				z["phone"] = i["address_phone"]
				z["info"] = i["address_info"]
				x["address"] = z
			} else if x["order_sn"] != i["order_sn"] && len(x) != 0 {
				//下一个
				data = append(data, x)
				x = nil
				x = make(map[string]interface{})
				x["order_sn"] = i["order_sn"]
				x["id"] = []string{i["sku_id"]}
				x["name"] = []string{i["title"]}
				x["date"] = i["created_time"]
				x["price"] = i["pay_money"]
				x["imgurls"] = []string{i["first_image"]}
				x["shopnums"] = []string{i["num"]}
				x["unitprice"] = []string{i["price"]}
				if i["order_status"] == "1" {
					x["status"] = "已完成"
				} else if i["order_status"] == "2" {
					x["status"] = "未付款"
				} else if i["order_status"] == "3" {
					x["status"] = "待发货"
				} else if i["order_status"] == "4" {
					x["status"] = "已发货"
				} else if i["order_status"] == "5" {
					x["status"] = "已取消"
				}
				//收货信息
				z := make(map[string]string)
				z["name"] = i["address_name"]
				z["phone"] = i["address_phone"]
				z["info"] = i["address_info"]
				x["address"] = z
			}
		}
		data = append(data, x)
		return data
	} else {
		return nil
	}
}

// AdminGetAllOrder 管理员获取所有订单
func AdminGetAllOrder(query string) []map[string]interface{} {
	var sqlStr string
	if query == "" {
		sqlStr = fmt.Sprintf("select a.order_sn,a.created_time,a.order_status,a.pay_money,a.order_money,a.address_name,a.address_phone,a.address_info,b.sku_id,b.num,b.price,s.title from order_master a inner join order_detail b on a.order_sn = b.order_sn inner join sku s on b.sku_id = s.id;")
	} else {
		sqlStr = fmt.Sprintf("select a.order_sn,a.created_time,a.order_status,a.pay_money,a.order_money,a.address_name,a.address_phone,a.address_info,b.sku_id,b.num,b.price,s.title from order_master a inner join order_detail b on a.order_sn = b.order_sn inner join sku s on b.sku_id = s.id and a.order_sn='%s';", query)
	}
	res, ok := Query(sqlStr)
	if ok {
		var data []map[string]interface{}
		orderMap := make(map[string]interface{})
		for _, value := range res {
			if orderMap["order_sn"] == value["order_sn"] {
				//上一个订单号和这个订单号相同，为同一个订单
				orderMap["children"] = append(orderMap["children"].([]map[string]interface{}), map[string]interface{}{
					"skuID": value["sku_id"],
					"num":   value["num"],
					"price": value["price"],
					"title": value["title"],
				})
			} else {
				if len(orderMap) != 0 {
					data = append(data, orderMap)
					orderMap = make(map[string]interface{})
				}
				orderMap["order_sn"] = value["order_sn"]
				orderMap["created_time"] = value["created_time"]
				orderMap["order_status"] = value["order_status"]
				orderMap["pay_money"] = value["pay_money"]
				orderMap["order_money"] = value["order_money"]
				orderMap["address_name"] = value["address_name"]
				orderMap["address_phone"] = value["address_phone"]
				orderMap["address_info"] = value["address_info"]
				orderMap["children"] = []map[string]interface{}{
					{
						"skuID": value["sku_id"],
						"num":   value["num"],
						"price": value["price"],
						"title": value["title"],
					},
				}
			}
		}
		data = append(data, orderMap)
		return data
	} else {
		return nil
	}
}

// AdminUpdateAddress 管理员更改订单收货地址
func AdminUpdateAddress(orderSn string, name string, phone string, info string) error {
	nowTime := time.Now().Format("2006-01-02 15:03:04")
	sqlStr := fmt.Sprintf("update order_master set address_name='%s',address_phone='%s',address_info='%s',update_time='%s' where order_sn='%s'", name, phone, info, nowTime, orderSn)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// AdminUpdateShipping 管理员发货
func AdminUpdateShipping(orderSn string, sn string) error {
	nowTime := time.Now().Format("2006-01-02 15:03:04")
	sqlStr := fmt.Sprintf("update order_master set shipping_sn='%s',shipping_time='%s',update_time='%s',order_status=4 where order_sn='%s'", sn, nowTime, nowTime, orderSn)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}
