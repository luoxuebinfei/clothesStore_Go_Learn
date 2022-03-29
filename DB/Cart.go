package DB

import "fmt"

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
			x = append(x, y)
		}
		return x
	} else {
		return nil
	}
}

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
