package DB

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

//购物车相关数据库实现

// CartInfo 获取购物车的全部数据及返回构成json的格式
func CartInfo(uid string) []map[string]string {
	sqlStr := fmt.Sprintf("select goods_id,goods_num,status,checked from cart_item where uid='%s'", uid)
	res, ok := Query(sqlStr)
	if ok {
		var x []map[string]string
		for _, value := range res {
			y := make(map[string]string)
			a := GetGoodsInfoAll(value["goods_id"])
			s := GetSpuId(a[0]["spu_id"])[0]["product_specs"]
			y["shopname"] = a[0]["title"]                               //商品标题
			y["price"] = a[0]["price"]                                  //商品价格
			y["buyNum"] = value["goods_num"]                            //商品数量
			y["shopId"] = value["goods_id"]                             //商品id
			y["props"] = fmt.Sprintf("%s %s", s, a[0]["product_specs"]) //规格
			y["imgurl"] = a[0]["first_image"]                           //商品图片
			y["checked"] = value["checked"]                             //选中情况
			x = append(x, y)
		}
		return x
	} else {
		return nil
	}
}

// DeleteCart 删除购物车中商品
func DeleteCart(skuId string, uid string) error {
	//skuId是商品id，uid是用户id
	sqlStr := fmt.Sprintf("delete from cart_item where goods_id='%s' and uid='%s'", skuId, uid)
	err := Exec(sqlStr)
	if err != nil {
		fmt.Println("删除商品错误：", err)
		return err
	}
	return nil
}

// UpdateCartNum 更新购物车中商品的数量
func UpdateCartNum(skuId string, num string, uid string) error {
	queryNum, _ := strconv.Atoi(num)
	stockNum := QuerySkuStock(skuId)
	if queryNum > stockNum && stockNum != -1 {
		return errors.New("商品库存不足！")
	} else if stockNum == -1 {
		return errors.New("商品不存在！")
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := fmt.Sprintf("update cart_item set goods_num = '%s',update_time='%s' where uid='%s' and goods_id='%s';", num, nowTime, uid, skuId)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// ChangeChecked 切换当前商品的选中状态
func ChangeChecked(skuId string, status string, uid string) error {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := fmt.Sprintf("update cart_item set checked = '%s',update_time='%s' where goods_id='%s' and uid='%s';", status, nowTime, skuId, uid)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// AddCart 将商品添加到购物车
func AddCart(skuId string, buyNum string, uid string) error {
	i, _ := strconv.Atoi(buyNum)
	if i < 1 {
		return errors.New("购买数量不合法")
	}
	//查询商品id是否存在
	sqlStr := fmt.Sprintf("select * from sku where id='%s'", skuId)
	_, ok := Query(sqlStr)
	if ok {
		//查询商品id是否存在于购物车中
		sqlStr = fmt.Sprintf("select * from cart_item where uid='%s' and goods_id='%s';", uid, skuId)
		res, ok := Query(sqlStr)
		if ok {
			//存在则更新购物中此商品数量
			nowTime := time.Now().Format("2006-01-02 15:04:05")
			oldNum, _ := strconv.Atoi(res[0]["goods_num"])
			sqlStr = fmt.Sprintf("update cart_item set goods_num='%d',update_time='%s' where uid='%s' and goods_id='%s'", i+oldNum, nowTime, uid, skuId)
			err := Exec(sqlStr)
			if err != nil {
				return err
			}
		} else {
			//不存在则新加入到购物车中
			nowTime := time.Now().Format("2006-01-02 15:04:05")
			sqlStr = fmt.Sprintf("insert into cart_item (uid,goods_id,goods_num,status,checked,create_time,update_time) values ('%s','%s','%s',1,1,'%s','%s')", uid, skuId, buyNum, nowTime, nowTime)
			err := Exec(sqlStr)
			if err != nil {
				return err
			}
		}

	} else {
		return errors.New("没有此商品")
	}
	return nil
}
