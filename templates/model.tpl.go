package templates

var ModelTpl = `package model

import (
	"errors"
	. "grpc-server/entity"
	DB "grpc-server/library/database"
)

type {{ServerName}}Model struct {
}

//查找多条数据
func (u *{{ServerName}}Model) Find({{EntityName}}Query {{ServerName}}Query) ([]{{ServerName}}, error) {
	dbConn := DB.GetDB({{DbName}})
	defer dbConn.Close()
	{{EntityName}} := make([]{{ServerName}}, 0)
	//limit
	if {{EntityName}}Query.PageNum > 0 && {{EntityName}}Query.PageSize > 0 {
		limitSlic := getLimit({{EntityName}}Query.PageNum, {{EntityName}}Query.PageSize)
		dbConn.Limit(limitSlic[0], limitSlic[1])
	}
	err := dbConn.Find(&{{EntityName}}, {{EntityName}}Query.Conditions)
	return {{EntityName}}, err
}

//根据id查找单条数据
func (u *{{ServerName}}Model) Get({{EntityName}} {{ServerName}}) ({{ServerName}}, error) {
	dbConn := DB.GetDB({{DbName}})
	_, err := dbConn.Get(&{{EntityName}})
	defer dbConn.Close()
	return {{EntityName}}, err
}

//插入
func (u *{{ServerName}}Model) Insert({{EntityName}} {{ServerName}}) ({{ServerName}}, error) {
	dbConn := DB.GetDB({{DbName}})
	defer dbConn.Close()
	affected, err := dbConn.Insert(&{{EntityName}})
	if err != nil {
		return {{EntityName}}, err
	}
	if affected < 1 {
		err = errors.New("插入影响行数: 0")
		return {{EntityName}}, err
	}
	return {{EntityName}}, err
}

//更新
func (u *{{ServerName}}Model) Update(conditions {{ServerName}}, {{EntityName}} {{ServerName}}) (affected int64, err error) {
	dbConn := DB.GetDB({{DbName}})
	affected, err = dbConn.Update({{EntityName}}, conditions)
	defer dbConn.Close()
	return
}
`
