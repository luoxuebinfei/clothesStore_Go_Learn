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
			y["quantity"] = value["goods_num"]                          //商品数量
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
func DeleteCart(skuId string, uid string) int {
	//skuId是商品id，uid是用户id
	sqlStr := fmt.Sprintf("delete from cart_item where goods_id='%s' and uid='%s'", skuId, uid)
	err := Exec(sqlStr)
	if err != nil {
		fmt.Println("删除商品错误：", err)
		return -1
	}
	return 0
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
