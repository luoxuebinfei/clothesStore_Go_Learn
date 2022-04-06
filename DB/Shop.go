package DB

import (
	"fmt"
	"strconv"
)

//商品数据库相关操作

// GetGoodsInfoAll 获取商品相关的全部信息
func GetGoodsInfoAll(id string) []map[string]string {
	sqlStr := fmt.Sprintf("select * from sku where id='%s'", id)
	res, ok := Query(sqlStr)
	if ok {
		return res
	} else {
		return nil
	}

}

// GetSpuId 获取商品所有信息
func GetSpuId(id string) []map[string]string {
	sqlStr := fmt.Sprintf("select * from spu where id ='%s'", id)
	res, ok := Query(sqlStr)
	if ok {
		return res
	} else {
		return nil
	}
}

// QuerySkuStock 查询商品库存
func QuerySkuStock(id string) int {
	sqlStr := fmt.Sprintf("select stock_num from sku where id='%s';", id)
	res, ok := Query(sqlStr)
	if ok {
		stockNum, _ := strconv.Atoi(res[0]["stock_num"])
		return stockNum
	} else {
		return -1
	}
	//-1为没有此商品
}
