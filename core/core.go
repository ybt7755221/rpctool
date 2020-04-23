package core

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"rpctool/templates"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type SqlField struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

func GetDB(connString string) *sql.DB {
	var err error
	db, err := sql.Open("mysql", connString)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	if err != nil {
		fmt.Printf("connect mysql fail ! [%s]", err)
	} else {
		fmt.Println("connect to mysql success")
	}
	return db
}

//获取mysql结构的slic
func GetMysqlStruct(connString string, tableName string) ([]SqlField, error) {
	var slic = make([]SqlField, 0)
	var sqlField = new(SqlField)
	db := GetDB(connString)
	defer db.Close()
	rows, err := db.Query("desc " + tableName)
	defer rows.Close()
	if err != nil {
		return slic, err
	}
	if err != nil {
		return slic, err
	}
	for rows.Next() {
		//定义变量接收查询数据
		err := rows.Scan(&sqlField.Field, &sqlField.Type, &sqlField.Null, &sqlField.Key, &sqlField.Default, &sqlField.Extra)
		if err != nil {
			return slic, err
		}
		slic = append(slic, *sqlField)
	}
	return slic, err
}

//处理名称
func DealServerName(tableName string) string {
	var slic = make([]string, 0)
	var serverName string
	slic = strings.Split(tableName, "_")
	for i := 0; i < len(slic); i++ {
		serverName += FirstToUpper(slic[i])
	}
	return serverName
}

//首字母小写
func FirstToLower(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

//首字母大写其他小写
func FirstToUpper(s string) string {
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

//生成model/service文件
func GeneFile(dbName string, tableName string, fileDir string, ftype string) {
	var fileString string
	var valMap map[string]string
	var ServerName = DealServerName(tableName)
	valMap = map[string]string{
		"{{ServerName}}": ServerName,
		"{{EntityName}}": FirstToLower(ServerName),
	}
	switch ftype {
	case "model":
		fileString = templates.ModelTpl
		valMap["{{DbName}}"] = FirstToUpper(dbName)
	case "service":
		fileString = templates.ServiceTpl
		valMap["{{TableName}}"] = tableName
	}
	//替换关键字
	for k, v := range valMap {
		fileString = strings.ReplaceAll(fileString, k, v)
	}
	//生成文件
	fileDirName, fileName := GetFileInfo(fileDir, tableName, ftype)
	WriteFile(fileDirName, fileName, fileString, 0755)
}

//生成entity文件
func GeneratorEntity(connString string, tableName string, fileDir string) {
	var fileString = templates.EntityTpl
	fieldSlic, err := GetMysqlStruct(connString, tableName)
	if err != nil {
		fmt.Println(err.Error())
	}
	format := map[string]string{
		"{{ServerName}}":  DealServerName(tableName),
		"{{TableSchema}}": convertMysqlToEntity(fieldSlic),
	}
	//替换关键字
	for k, v := range format {
		fileString = strings.ReplaceAll(fileString, k, v)
	}
	fileString = strings.ReplaceAll(fileString, "'", "`")
	//生成文件
	fileDirName, fileName := GetFileInfo(fileDir, tableName, "entity")
	WriteFile(fileDirName, fileName, fileString, 0755)
}

//生成proto文件
func Generator(connString string, tableName string, fileDir string) {
	var fileString = templates.ProtoTpl
	//获取mysql结构
	fieldSlic, err := GetMysqlStruct(connString, tableName)
	if err != nil {
		fmt.Println(err.Error())
	}
	//整理参数
	format := map[string]string{
		"{{TableSchema}}": ConvertMysqlTypeToProtoType(fieldSlic),
		"{{TableName}}":   tableName,
		"{{ServerName}}":  DealServerName(tableName),
	}
	//替换关键字
	for k, v := range format {
		fileString = strings.ReplaceAll(fileString, k, v)
	}
	//生成文件
	fileDirName, fileName := GetFileInfo(fileDir, tableName, "proto")
	WriteFile(fileDirName, fileName, fileString, 0755)
}

func convertMysqlToEntity(fieldSlic []SqlField) string {
	var schema string
	//处理变量
	for i := 0; i < len(fieldSlic); i++ {
		if i > 0 {
			schema += "    "
		}
		schema += DealServerName(fieldSlic[i].Field)
		if strings.Index(fieldSlic[i].Type, "bigint") > -1 {
			schema += "\tint64 "
		} else if strings.Index(fieldSlic[i].Type, "tinyint") > -1 {
			schema += "\tint8 "
		} else if strings.Index(fieldSlic[i].Type, "smallint") > -1 {
			schema += "\tint16 "
		} else if strings.Index(fieldSlic[i].Type, "int") > -1 {
			schema += "\tint "
		} else if strings.Index(fieldSlic[i].Type, "text") > -1 {
			schema += "\tstring "
		} else if strings.Index(fieldSlic[i].Type, "char") > -1 {
			schema += "\tstring "
		} else if strings.Index(fieldSlic[i].Type, "enum") > -1 {
			schema += "\tstring "
		} else if strings.Index(fieldSlic[i].Type, "blob") > -1 {
			schema += "\tstring "
		} else if strings.Index(fieldSlic[i].Type, "float") > -1 {
			schema += "\tfloat32 "
		} else if strings.Index(fieldSlic[i].Type, "double") > -1 {
			schema += "\tfloat64 "
		} else if strings.Index(fieldSlic[i].Type, "date") > -1 {
			schema += "\tstring "
		} else if strings.Index(fieldSlic[i].Type, "time") > -1 {
			schema += "\tstring "
		} else {
			schema += "\tstring "
		}
		schema += "`json:\"" + fieldSlic[i].Field + "\"`"
		if i < len(fieldSlic)-1 {
			schema += "\n"
		}
	}
	return schema
}

//获取文件路径和文件名
func GetFileInfo(fileDir string, tableName string, ftype string) (dirName string, fileName string) {
	if ftype == "proto" {
		if fileDir != "" {
			dirName = filepath.Join(fileDir, ftype+"s", tableName)
		} else {
			dirName = filepath.Join("output", ftype+"s", tableName)
		}
		fileName = tableName + ".proto"
	} else if ftype == "entity" {
		if fileDir != "" {
			dirName = filepath.Join(fileDir, ftype)
		} else {
			dirName = filepath.Join("output", ftype)
		}
		fileName = tableName + ".go"
	} else {
		if fileDir != "" {
			dirName = filepath.Join(fileDir, ftype)
		} else {
			dirName = filepath.Join("output", ftype)
		}
		fileName = DealServerName(tableName) + FirstToUpper(ftype) + ".go"
		fileName = FirstToLower(fileName)
	}
	return
}

//convert mysql type to proto type
func ConvertMysqlTypeToProtoType(fieldSlic []SqlField) string {
	var schema string
	//处理变量
	for i := 0; i < len(fieldSlic); i++ {
		if i > 0 {
			schema += "    "
		}
		numStr := strconv.Itoa(i + 1)
		if strings.Index(fieldSlic[i].Type, "bigint") > -1 {
			schema += "int64 "
		} else if strings.Index(fieldSlic[i].Type, "int") > -1 {
			schema += "int32 "
		} else if strings.Index(fieldSlic[i].Type, "text") > -1 {
			schema += "string "
		} else if strings.Index(fieldSlic[i].Type, "char") > -1 {
			schema += "string "
		} else if strings.Index(fieldSlic[i].Type, "enum") > -1 {
			schema += "string "
		} else if strings.Index(fieldSlic[i].Type, "blob") > -1 {
			schema += "string "
		} else if strings.Index(fieldSlic[i].Type, "float") > -1 {
			schema += "float "
		} else if strings.Index(fieldSlic[i].Type, "double") > -1 {
			schema += "double "
		} else if strings.Index(fieldSlic[i].Type, "date") > -1 {
			schema += "string "
		} else if strings.Index(fieldSlic[i].Type, "time") > -1 {
			schema += "string "
		} else {
			schema += "string "
		}
		schema += fieldSlic[i].Field + " = " + numStr + ";"
		if i < len(fieldSlic)-1 {
			schema += "\n"
		}
	}
	return schema
}

//写入文件
func WriteFile(fileDir string, fileName string, file string, mode os.FileMode) error {
	_, err := os.Stat(fileDir)
	if err != nil {
		err = os.MkdirAll(fileDir, mode)
		if err != nil {
			log.Fatalln(err.Error() + ": " + fileDir)
		}
	}
	fn := filepath.Join(fileDir, fileName)
	err = ioutil.WriteFile(fn, []byte(file), mode)
	if err != nil {
		log.Fatalln(err.Error() + ": " + fn)
	}
	fmt.Println("success create :" + fn)
	return err
}
