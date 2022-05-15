package DB

import (
	"errors"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"github.com/tidwall/gjson"
	"sort"
	"strconv"
	"strings"
	"time"
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

// GetShopInfo 在shopInfo中返回的商品数据
func GetShopInfo(skuId string) (map[string]interface{}, error) {
	sqlStr := fmt.Sprintf("select a.price,a.title,a.first_image,a.weight,a.brand_id,a.spu_id,a.stock_num,a.product_specs as a_product_specs,b.images,b.details_images,b.attribute_list,b.bro_spu,b.product_specs,c.brand_name from (sku a inner join spu b on a.id='%s' and a.spu_id=b.id) inner join brands c on a.brand_id=c.id", skuId)
	res, ok := Query(sqlStr)
	if ok {
		data := make(map[string]interface{})
		data["detailsImages"] = res[0]["details_images"] //详情图
		data["specImages"] = res[0]["images"]
		stockNum, _ := strconv.Atoi(res[0]["stock_num"])
		if stockNum > 0 {
			data["stock"] = true
		} else {
			data["stock"] = false
		}
		data["skuName"] = res[0]["title"]
		data["price"] = res[0]["price"]
		data["parameter2"] = res[0]["attribute_list"]
		data["brand_name"] = res[0]["brand_name"]
		data["weight"] = res[0]["weight"]
		//对同一个商品不同规格的数据组成
		sqlStr = fmt.Sprintf("select * from sku where spu_id='%s'", res[0]["spu_id"])
		res1, _ := Query(sqlStr)
		var size []map[string]interface{}
		for _, value := range res1 {
			m := make(map[string]interface{})
			m["skuid"] = value["id"]
			m["text"] = value["product_specs"]
			stockNum, _ := strconv.Atoi(value["stock_num"])
			if stockNum > 0 {
				m["stock"] = true
			} else {
				m["stock"] = false
			}
			if value["id"] == skuId {
				m["isChecked"] = true
			} else {
				m["isChecked"] = false
			}
			size = append(size, m)
		}
		//对同一类别的商品的集合
		switch t := gjson.Parse(res[0]["bro_spu"]).Value().(type) {
		case []interface{}:
			var color []map[string]interface{}
			for _, value := range t {
				m := make(map[string]interface{})
				spu, _ := strconv.Atoi(res[0]["spu_id"])
				if int(value.(float64)) == spu {
					//选中的
					m["skuid"] = skuId
					m["text"] = res[0]["product_specs"]
					stockNum, _ := strconv.Atoi(res[0]["stock_num"])
					if stockNum > 0 {
						m["stock"] = true
					} else {
						m["stock"] = false
					}
					m["isChecked"] = true
					m["imgurl"] = res[0]["first_image"]
					color = append(color, m)
				} else {
					sqlStr = fmt.Sprintf("select * from sku where spu_id='%d' and product_specs='%s';", int(value.(float64)), res[0]["a_product_specs"])
					res2, ok := Query(sqlStr)
					if ok {
						m["skuid"] = res2[0]["id"]
						sqlStr = fmt.Sprintf("select * from spu where id='%d'", int(value.(float64)))
						res3, _ := Query(sqlStr)
						m["text"] = res3[0]["product_specs"]
						stockNum, _ := strconv.Atoi(res2[0]["stock_num"])
						if stockNum > 0 {
							m["stock"] = true
						} else {
							m["stock"] = false
						}
						m["isChecked"] = false
						m["imgurl"] = res2[0]["first_image"]
					} else {
						m["skuid"] = ""
						m["text"] = ""
						m["stock"] = false
						m["isChecked"] = false
						m["imgurl"] = ""
					}
					color = append(color, m)
				}
				data["attrs"] = map[string]interface{}{"color": color, "size": size}
			}
		}

		return data, nil
	} else {
		return nil, errors.New("商品不存在")
	}
}

type Brands [][]string

// SearchShop 搜索商品
func SearchShop(keyword string, queryBrand string, querySize string, page string) map[string]interface{} {
	//分页，每次返回20条数据
	start := 0
	page_, _ := strconv.Atoi(page)
	if page_ > 1 {
		start += page_ * 20
	} else if page_ < 1 {
		start = 0
	}
	sqlStr := ""
	if queryBrand == "" && querySize == "" {
		sqlStr = fmt.Sprintf("select a.*,b.brand_name from sku a inner join brands b where a.title like '%%%s%%' and a.brand_id=b.id group by a.spu_id", strings.Join(strings.Split(keyword, " "), "%"))
	} else if queryBrand != "" && querySize == "" {
		sqlStr = fmt.Sprintf("select a.*,b.brand_name from sku a inner join brands b where a.title like '%%%s%%' and a.brand_id=b.id and b.brand_name='%s' group by a.spu_id", strings.Join(strings.Split(keyword, " "), "%"), queryBrand)
	} else if queryBrand != "" && querySize != "" {
		sqlStr = fmt.Sprintf("select a.*,b.brand_name from sku a inner join brands b where a.title like '%%%s%%' and a.brand_id=b.id and b.brand_name='%s' and a.product_specs='%s' group by a.spu_id", strings.Join(strings.Split(keyword, " "), "%"), queryBrand, querySize)
	} else {
		sqlStr = fmt.Sprintf("select a.*,b.brand_name from sku a inner join brands b where a.title like '%%%s%%' and a.brand_id=b.id and a.product_specs='%s' group by a.spu_id", strings.Join(strings.Split(keyword, " "), "%"), querySize)
	}
	res, ok := Query(sqlStr)
	if ok {
		data := make(map[string]interface{})
		var s []map[string]string
		var brands Brands
		var size []string
		for _, i := range res {
			m := make(map[string]string)
			m["sku"] = i["id"]
			m["name"] = i["title"]
			m["price"] = i["price"]
			m["imgurl"] = i["first_image"]
			s = append(s, m)
			size = append(size, i["product_specs"])
			brands = append(brands, []string{strings.ToUpper(pinyin.LazyPinyin(i["brand_name"], pinyin.NewArgs())[0][0:1]), i["brand_name"]})
		}
		if len(s) < 20 {
			end := len(s)
			s = s[start:end]
		} else {
			end := start + 20
			s = s[start:end]
		}
		data["goods_list"] = s
		x := make(map[string]string)
		for _, i := range brands {
			//["Y","雅鹿"]
			x[i[1]] = i[0]
		}
		brands = nil
		for key, value := range x {
			brands = append(brands, []string{value, key})
		}
		sort.Sort(brands)
		data["brands"] = brands
		x = make(map[string]string)
		for index, i := range size {
			x[i] = strconv.Itoa(index)
		}
		size = nil
		for key, _ := range x {
			size = append(size, key)
		}
		data["size"] = size
		return data
	} else {
		return nil
	}
}

// Len 实现sort.Interface接口的获取元素数量方法
func (m Brands) Len() int {
	return len(m)
}

// Less 实现sort.Interface接口的比较元素方法
func (m Brands) Less(i, j int) bool {
	return m[i][0] < m[j][0]
}

// Swap 实现sort.Interface接口的交换元素方法
func (m Brands) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// AdminGetAllShop 管理员获取所有商品信息并返回预定格式
func AdminGetAllShop() []map[string]interface{} {
	sqlStr := fmt.Sprintf("select a.id as spuID,a.shop_name,a.images,a.details_images,a.product_specs as spu_product_specs,a.attribute_list,b.id as skuID,b.goods_name,b.price,b.stock_num,b.title,b.product_specs as sku_product_specs,b.first_image,b.status,b.weight from spu a left join sku b on b.spu_id=a.id order by a.id;")
	res, ok := Query(sqlStr)
	if ok {
		var spuID string
		var data []map[string]interface{}      //返回的数据
		spuMap := make(map[string]interface{}) //spu集合
		for _, value := range res {
			if value["spuID"] == spuID {
				//如果当前的spuID和上一个spuID相等，则将sku信息添加到children中
				spuMap["children"] = append(spuMap["children"].([]map[string]interface{}), map[string]interface{}{
					"skuID":             value["skuID"],
					"goods_name":        value["goods_name"],
					"price":             value["price"],
					"stock_num":         value["stock_num"],
					"title":             value["title"],
					"sku_product_specs": value["sku_product_specs"],
					"first_image":       value["first_image"],
					"status":            value["status"],
					"weight":            value["weight"],
				})
			} else {
				//如果当前的spuID和上一个spuID不相等，则建一个新的，并将之前的加入到data中
				if len(spuMap) != 0 {
					data = append(data, spuMap)
				}
				spuMap = make(map[string]interface{})
				spuMap["spuID"] = value["spuID"]
				spuID = value["spuID"]
				spuMap["shop_name"] = value["shop_name"]
				spuMap["spu_product_specs"] = value["spu_product_specs"]
				spuMap["attribute_list"] = value["attribute_list"]
				spuMap["images"] = value["images"]
				spuMap["details_images"] = value["details_images"]
				spuMap["children"] = []map[string]interface{}{
					{
						"skuID":             value["skuID"],
						"goods_name":        value["goods_name"],
						"price":             value["price"],
						"stock_num":         value["stock_num"],
						"title":             value["title"],
						"sku_product_specs": value["sku_product_specs"],
						"first_image":       value["first_image"],
						"status":            value["status"],
						"weight":            value["weight"],
					},
				}
			}
		}
		if len(spuMap) != 0 {
			data = append(data, spuMap)
		}
		return data
	} else {
		return nil
	}
}

// AdminAddSpu 管理员增加spu信息
func AdminAddSpu(shopName string, product_specs string, attribute_list string, images string, details_images string) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := fmt.Sprintf("insert into spu (shop_name,product_specs,attribute_list,images,details_images,created_time,update_time) values ('%s','%s','%s','%s','%s','%s','%s')", shopName, product_specs, attribute_list, images, details_images, nowTime, nowTime)
	err := Exec(sqlStr)
	if err != nil {
		fmt.Println(err)
	}
	id := strconv.FormatInt(LastInsertId, 10)
	sqlStr = fmt.Sprintf("update spu set bro_spu='%s' where id='%s'", []string{id}, id)
	err = Exec(sqlStr)
	if err != nil {
		fmt.Println(err)
	}
}

// AdminDeleteSpu 删除spu
func AdminDeleteSpu(spuID string) error {
	sqlStr := fmt.Sprintf("delete from spu where id='%s'", spuID)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// AdminUpdateSpu 更新spu
func AdminUpdateSpu(spuID string, shop_name string, product_specs string, images string, details_images string, attribute_list string) error {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := fmt.Sprintf("update spu set shop_name='%s',product_specs='%s',images='%s',details_images='%s',attribute_list='%s',update_time='%s' where id='%s'", shop_name, product_specs, images, details_images, attribute_list, nowTime, spuID)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// AdminAddSku 添加sku
func AdminAddSku(spuID string, goods_name string, price string, title string, stock_num string, product_specs string, first_image string, weight string, status bool) error {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	var a string
	if status {
		a = "1"
	} else {
		a = "0"
	}
	sqlStr := fmt.Sprintf("insert into sku (goods_name, brand_id, spu_id, price, stock_num, title, product_specs, first_image, status, weight, sort, created_time, update_time) values ('%s',1,'%s','%s','%s','%s','%s','%s','%s','%s',0,'%s','%s')", goods_name, spuID, price, stock_num, title, product_specs, first_image, a, weight, nowTime, nowTime)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

// AdminDeleteSku 删除sku
func AdminDeleteSku(skuID string) error {
	sqlStr := fmt.Sprintf("delete from sku where id='%s'", skuID)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}

//更新sku
func AdminUpdateSku(skuID string, goods_name string, price string, title string, stock_num string, product_specs string, weight string, status bool) error {
	nowTime := time.Now().Format("2006-01-02 15:03:04")
	var a string
	if status {
		a = "1"
	} else {
		a = "0"
	}
	sqlStr := fmt.Sprintf("update sku set goods_name='%s',price='%s',title='%s',stock_num='%s',product_specs='%s',weight='%s',status='%s',update_time='%s' where id='%s'", goods_name, price, title, stock_num, product_specs, weight, a, nowTime, skuID)
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}
