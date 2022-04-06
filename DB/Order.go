package DB

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"time"
)

//订单表数据库操作
//包括生成订单，查询订单等操作

// GenerateOrderSn 提交订单，生成订单号
func GenerateOrderSn(uid string, data map[string]interface{}) (int, string) {
	//20220401014810 + 4位随机码 + id
	nowTime := time.Now().Format("20060102150405")
	rand.Seed(time.Now().Unix())
	randCode := rand.Intn(10000)
	orderSn := fmt.Sprintf("%s%d%s", nowTime, randCode, uid)
	//fmt.Println(orderSn)
	a := data["data"].([]interface{})
	x := make(map[string][]string) //商品id数组{商品id:[商品数量,单价]}
	for _, value := range a {
		skuId := strconv.FormatFloat(value.(map[string]interface{})["shopId"].(float64), 'E', -1, 64)
		buyNum := int(value.(map[string]interface{})["quantity"].(float64))
		if buyNum <= 0 {
			return -3, "参数不正确！"
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
					return -1, ""
				}
				x[skuId] = []string{strconv.Itoa(buyNum), price}
			} else {
				//购买数量大于库存时，返回
				//fmt.Println("商品库存不足")
				return -2, fmt.Sprintf("%s 此商品库存不足", res[0]["goods_name"])
			}
		} else {
			//fmt.Println("商品不存在")
			return -3, fmt.Sprintf("商品id:%s 此商品不存在", skuId)
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
			return -1, ""
		}
		//删除购物车中商品
		DeleteCart(key, uid)
	}
	total_ := fmt.Sprintf("%.2f", total.InexactFloat64())
	//将订单号存入订单主表
	sqlStr := fmt.Sprintf("insert into order_master (order_sn, uid, address_name, address_phone, address_info, created_time, order_money, pay_money, order_status, update_time) value ('%s','%s','%s','%s','%s','%s','%s','%s',2,'%s')", orderSn, uid, addressName, phone, address, nowTime, total_, total_, nowTime)
	err := Exec(sqlStr)
	if err != nil {
		return -1, ""
	}

	return 1, ""
}

// QueryOrderSn 查询订单
func QueryOrderSn(orderSn string, uid string) (int, string, map[string]interface{}) {
	sqlStr := fmt.Sprintf("select a.*,b.sku_id,b.num,b.price,c.title,c.first_image from (order_master a inner join order_detail b on a.order_sn=b.order_sn and a.order_sn='%s' and a.uid='%s' and a.order_status!=0) inner join sku c on b.sku_id=c.id", orderSn, uid)
	res, ok := Query(sqlStr)
	if ok {
		//fmt.Println(res)
		data := make(map[string]interface{})
		area := make(map[string]string)
		var shoplist []map[string]string
		for _, i := range res {
			//收货信息
			area["address_name"] = i["address_name"]
			area["address_phone"] = i["address_phone"]
			area["address_info"] = i["address_info"]
			//订单创建时间
			data["created_time"] = i["created_time"]
			//订单状态
			data["order_status"] = i["order_status"]
			//订单号
			data["order_sn"] = i["order_sn"]
			//商品信息
			shop := make(map[string]string)
			shop["id"] = i["sku_id"]
			shop["price"] = i["price"]
			shop["buy_num"] = i["num"]
			shop["title"] = i["title"]
			shop["first_image"] = i["first_image"]
			shoplist = append(shoplist, shop)
			//金额信息
			data["order_money"] = i["order_money"]
			data["pay_money"] = i["pay_money"]
			//物流信息
			data["shipping_comp_name"] = i["shipping_comp_name"]
			data["shipping_sn"] = i["shipping_sn"]
			data["shipping_time"] = i["shipping_time"] //发货时间
			//收货时间
			data["receive_time"] = i["receive_time"]
		}
		data["shop"] = shoplist
		data["area"] = area
		return 1, "", data
	} else {
		return -2, "订单号不存在！", nil
	}
	return 0, "", nil
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
