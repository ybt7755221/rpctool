package core

import (
	"testing"
)

func TestGetDB(t *testing.T) {
	db := GetDB("xin:48sdf37EB7@tcp(rm-2ze0q80w59h4uyvx4rw.mysql.rds.aliyuncs.com:3306)/xin?charset=utf8")
	rows, err := db.Query("desc crm_clue")
	if err != nil {
		t.Log(err.Error())
	}
	var slic = make([]SqlField, 0)
	var sqlField = new(SqlField)
	if err != nil {
		t.Logf("query faied, error:[%v]", err.Error())
		return
	}
	for rows.Next() {
		//定义变量接收查询数据
		err := rows.Scan(&sqlField.Field, &sqlField.Type, &sqlField.Null, &sqlField.Key, &sqlField.Default, &sqlField.Extra)
		if err != nil {
			t.Logf("get data failed, error:[%v]", err.Error())
		}
		slic = append(slic, *sqlField)
	}
	t.Log(slic)
	//关闭结果集（释放连接）
	rows.Close()
}

func TestGetMysqlStruct(t *testing.T) {
	slic, err := GetMysqlStruct("xin:48sdf37EB7@tcp(rm-2ze0q80w59h4uyvx4rw.mysql.rds.aliyuncs.com:3306)/xin?charset=utf8", "crm_clue")
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(slic)
}

func TestDealServerName(t *testing.T) {
	t.Log(DealServerName("gin_contest"))
}

func TestGenerator(t *testing.T) {
	Generator("xin:48sdf37EB7@tcp(rm-2ze0q80w59h4uyvx4rw.mysql.rds.aliyuncs.com:3306)/xin?charset=utf8", "gin_user_fields", "")
	t.Log("sucess")
}

func BenchmarkGetDB(b *testing.B) {
	db := GetDB("GinUser:userGin@tcp(127.0.0.1:3306)/gin?charset=utf8")
	rows, err := db.Query("desc gin_contents")
	if err != nil {
		b.Log(err.Error())
	}
	var slic = make([]SqlField, 0)
	var sqlField = new(SqlField)
	if err != nil {
		b.Logf("query faied, error:[%v]", err.Error())
		return
	}
	for rows.Next() {
		//定义变量接收查询数据
		err := rows.Scan(&sqlField.Field, &sqlField.Type, &sqlField.Null, &sqlField.Key, &sqlField.Default, &sqlField.Extra)
		if err != nil {
			b.Logf("get data failed, error:[%v]", err.Error())
		}
		slic = append(slic, *sqlField)
	}
	b.Log(slic)
	//关闭结果集（释放连接）
	rows.Close()
}
