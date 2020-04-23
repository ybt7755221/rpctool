package templates

var EntityTpl = `package entity

type {{ServerName}} struct {
	{{TableSchema}}
}

type {{ServerName}}Query struct {
	Conditions {{ServerName}} 'json:"conditions"'
	PageNum    int      	  'json:"page_num"'
	PageSize   int      	  'json:"page_size"'
}

type {{ServerName}}UpdateForm struct {
	Conditions {{ServerName}}	'json:"conditions"'
	Modifies   {{ServerName}}	'json:"modifies"'
}`
