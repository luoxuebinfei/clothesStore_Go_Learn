package DB

import "fmt"

func GetIndex() []map[string]string {
	sqlStr := fmt.Sprintf("select * from front_page")
	res, ok := Query(sqlStr)
	if ok {
		return res
	} else {
		return nil
	}
	return nil
}

// UpdateIndex 更新index内容
func UpdateIndex(imageItems string, adImages string, goodItems string) error {
	//fmt.Println(imageItems)
	//fmt.Println(adImages)
	//fmt.Println(goodItems)
	sqlStr := fmt.Sprintf("update front_page set imageItems='%s',adImages='%s',good_items='%s'", imageItems, adImages, goodItems)
	//return nil
	err := Exec(sqlStr)
	if err != nil {
		return err
	}
	return nil
}
