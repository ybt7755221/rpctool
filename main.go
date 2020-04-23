package main

import (
	"fmt"
	"rpctool/core"

	"github.com/droundy/goopt"
)

var (
	sqlTable   = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
	dbName     = goopt.String([]string{"-d", "--dbName"}, "", "database to build struct from")
	connString = goopt.String([]string{"-m", "--mysql"}, "", "mysql config")
	fileDir    = goopt.String([]string{"-o", "--out-file"}, "", "file dir,if not use 'output' ")
)

func init() {
	goopt.Description = func() string {
		return "rpctool is tool that can automaticlly generate proto file."
	}
	goopt.Version = "0.1"
	goopt.Summary = `rpctool --mysql user:password@tcp\(host:port\)/database\?charset=utf8 --dbName database --table tableName --out-file ./`
	goopt.Parse(nil)
}

func main() {
	if *connString == "" {
		fmt.Println("--mysql can not is empty")
		return
	}
	if *dbName == "" {
		fmt.Println("--dbName can not is empty")
		return
	}
	if *sqlTable == "" {
		fmt.Println("--table can not is empty")
		return
	}
	//生成proto文件
	core.Generator(*connString, *sqlTable, *fileDir)
	core.GeneratorEntity(*connString, *sqlTable, *fileDir)
	//生成model
	core.GeneFile(*dbName, *sqlTable, *fileDir, "model")
	//生成service
	core.GeneFile(*dbName, *sqlTable, *fileDir, "service")
}
