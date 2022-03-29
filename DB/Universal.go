package DB

//数据库文件的通用函数

// Exec 增、删、改
func Exec(SQL string) error {
	ret, err := Db.Exec(SQL) //增、删、改就靠这一条命令就够了，很简单
	if err != nil {
		//fmt.Println(err)
		return err
	}
	_, _ = ret.LastInsertId()
	//fmt.Println(insID)
	//if strconv.FormatInt(insID,10) == "0"{
	//	return errors.New("数据库中无此数据")
	//}
	return nil
}

// Query 通用查询
func Query(SQL string) ([]map[string]string, bool) {
	rows, err := Db.Query(SQL) //执行SQL语句，比如select * from users
	if err != nil {
		panic(err)
	}
	columns, _ := rows.Columns()            //获取列的信息
	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i, _ := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}
	ret := make([]map[string]string, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		err := rows.Scan(values...)  //开始读行，Scan函数只接受指针变量
		m := make(map[string]string) //用于存放1列的 [键/值] 对
		if err != nil {
			panic(err)
		}
		for i, colName := range columns {
			var raw_value = *(values[i].(*interface{})) //读出raw数据，类型为byte
			b, _ := raw_value.([]byte)
			v := string(b) //将raw数据转换成字符串
			m[colName] = v //colName是键，v是值
		}
		ret = append(ret, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	defer rows.Close()

	if len(ret) != 0 {
		return ret, true
	}
	return nil, false
}
