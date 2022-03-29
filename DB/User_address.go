package DB

import (
	"fmt"
	"time"
)

//用户收货地址相关数据库操作

// GetAddress 获取用户所有的收货地址
func GetAddress(uid string) []map[string]string {
	sqlStr := fmt.Sprintf("select id, name, phone, address, is_default from user_address where uid='%s'", uid)
	res, ok := Query(sqlStr)
	if ok {
		return res
	} else {
		return nil
	}
}

// AddNewAddress 新增地址
func AddNewAddress(uid string, name string, phone string, address string, isDefault string) bool {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	ok := OnlyOneDefault("", uid, false)
	if ok == -1 {
		//新增唯一一个地址时，设置为默认地址
		isDefault = "1"
	} else {
		if isDefault == "1" {
			//不是唯一地址且新增的地址为默认地址时，将原本的默认地址取消
			OnlyOneDefault("", uid, true)
		}
	}
	sqlStr := fmt.Sprintf("insert into user_address (uid, name, phone, address, is_default, created_time, update_time) value ('%s','%s','%s','%s','%s','%s','%s')", uid, name, phone, address, isDefault, nowTime, nowTime)
	err := Exec(sqlStr)
	if err != nil {
		return false
	}
	return true
}

// DeleteAddress 删除地址
func DeleteAddress(id string, uid string) bool {
	sqlStr := fmt.Sprintf("delete from user_address where id='%s' and uid='%s'", id, uid)
	err := Exec(sqlStr)
	if err != nil {
		return false
	}
	sqlStr = fmt.Sprintf("select * from user_address where uid='%s'", uid)
	res, ok := Query(sqlStr)
	if ok {
		if len(res) == 1 {
			sqlStr = fmt.Sprintf("update user_address set is_default=1 where id='%s'", res[0]["id"])
			Exec(sqlStr)
		}
	}
	return true
}

// UpdateAddress 更新地址
func UpdateAddress(id string, uid string, name string, phone string, address string, isDefault string) bool {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	if isDefault == "1" {
		//不是唯一地址且新增的地址为默认地址时，将原本的默认地址取消
		if OnlyOneDefault(id, uid, true) == -1 {
			return false
		}
	} else if isDefault == "0" {
		ok := OnlyOneDefault(id, uid, false)
		if ok == 1 {
			isDefault = "1"
		} else if ok == -1 {
			return false
		}
	}
	sqlStr := fmt.Sprintf("update user_address set name='%s',phone='%s',address='%s',is_default='%s',update_time='%s' where id='%s' and uid='%s'", name, phone, address, isDefault, nowTime, id, uid)
	err := Exec(sqlStr)
	if err != nil {
		//fmt.Errorf("更新用户收货地址时发生错误：%s",err)
		fmt.Println(err)
		return false
	}
	return true
}

// OnlyOneDefault 检测同uid默认选项的唯一性
func OnlyOneDefault(id string, uid string, isDefault bool) int {
	if id != "" {
		sqlStr := fmt.Sprintf("select * from user_address where id='%s' and uid='%s'", id, uid)
		_, ok := Query(sqlStr)
		if !ok {
			return -1
		}
	}
	sqlStr := fmt.Sprintf("select id,is_default from user_address where uid='%s'", uid)
	res, ok := Query(sqlStr)
	if ok {
		if isDefault {
			sqlStr = fmt.Sprintf("update user_address set is_default=0 where uid='%s'", uid)
			err := Exec(sqlStr)
			if err != nil {
				return 0
			}
		}
		if len(res) == 1 {
			//如果只有一个地址
			return 1
		}
		return 0
	} else {
		return -1
	}
}

// ChangeDefault 点击切换默认地址时使用
func ChangeDefault(id string, uid string) bool {
	ok := OnlyOneDefault(id, uid, true)
	if ok == -1 {
		return false
	}
	sqlStr := fmt.Sprintf("update user_address set is_default=1 where id='%s' and uid='%s'", id, uid)
	err := Exec(sqlStr)
	if err != nil {
		//fmt.Errorf("切换默认地址时发生错误：%s",err)
		//fmt.Println(err)
		return false
	}
	return true
}
