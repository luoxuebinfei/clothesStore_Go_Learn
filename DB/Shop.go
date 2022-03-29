package DB

import "fmt"

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

func GetSpuId(id string) []map[string]string {
	sqlStr := fmt.Sprintf("select * from spu where id ='%s'", id)
	res, ok := Query(sqlStr)
	if ok {
		return res
	} else {
		return nil
	}
}
