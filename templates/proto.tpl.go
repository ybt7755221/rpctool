package templates

var ProtoTpl = `
syntax = "proto3";

package {{TableName}}_proto;

service {{ServerName}} {
    rpc FindByPagination(QuerySchema) returns(FindRes);
    rpc FindOne({{ServerName}}Schema) returns(FindOneRes);
    rpc Create({{ServerName}}Schema) returns(FindOneRes);
    rpc Update(UpdateSchema) returns(FindRes);
}
//数据库结构
message {{ServerName}}Schema{
	{{TableSchema}}
}
//更新结构
message UpdateSchema {
    {{ServerName}}Schema conditions = 1;
    {{ServerName}}Schema modifies = 2;
}
//查询结构
message QuerySchema{
    {{ServerName}}Schema conditions = 1;
    int32 page_num = 2;
    int32 page_size = 3;
}
//查询返回对象
message FindOneRes{
    int32 code = 1;
    string msg = 2;
    {{ServerName}}Schema data = 3;
}
//查询返回string
message FindRes{
    int32 code = 1;
    string msg = 2;
    string data= 3;
}`
