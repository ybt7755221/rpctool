package templates

var ServiceTpl = `package service

import (
	"encoding/json"
	"fmt"
	"grpc-server/entity"
	"grpc-server/library/gutil"
	"grpc-server/model"
	{{TableName}}_proto "grpc-server/protos/{{TableName}}"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type {{ServerName}}Server struct {
	{{EntityName}}Model *model.{{ServerName}}Model
}

//注册服务
func (s *{{ServerName}}Server) Register(gs *grpc.Server) {
	{{TableName}}_proto.Register{{ServerName}}Server(gs, s)
}

//查询服务-带分页
func (s *{{ServerName}}Server) FindByPagination(ctx context.Context, in *{{TableName}}_proto.QuerySchema) (*{{TableName}}_proto.FindRes, error) {
	{{EntityName}} := new(entity.{{ServerName}})
	gutil.BeanUtil({{EntityName}}, in.Conditions)
	{{EntityName}}Query := entity.{{ServerName}}Query{}
	{{EntityName}}Query.Conditions = *{{EntityName}}
	{{EntityName}}Query.PageNum = int(in.PageNum)
	{{EntityName}}Query.PageSize = int(in.PageSize)
	fmt.Println({{EntityName}}Query)
	{{EntityName}}List, err := s.{{EntityName}}Model.Find({{EntityName}}Query)
	result := new({{TableName}}_proto.FindRes)
	if err != nil {
		result.Code = entity.ERROR
		result.Msg = err.Error()
		result.Data = ""
		return result, err
	}
	byteData, err := json.Marshal({{EntityName}}List)
	if err != nil {
		result.Code = entity.ERROR
		result.Msg = err.Error()
		result.Data = ""
		return result, err
	}
	result.Code = entity.SUCCESS
	result.Msg = entity.GetResultInfo(entity.SUCCESS)
	result.Data = string(byteData)
	return result, nil
}

//查询单条
func (s *{{ServerName}}Server) FindOne(ctx context.Context, in *{{TableName}}_proto.{{ServerName}}Schema) (*{{TableName}}_proto.FindOneRes, error) {
	{{EntityName}} := new(entity.{{ServerName}})
	gutil.BeanUtil({{EntityName}}, in)
	{{EntityName}}Res, err := s.{{EntityName}}Model.Get(*{{EntityName}})
	gutil.BeanUtil(in, &{{EntityName}}Res)
	res := new({{TableName}}_proto.FindOneRes)
	res.Code = 1000
	res.Msg = "ok"
	res.Data = in
	return res, err
}

//创建
func (s *{{ServerName}}Server) Create(ctx context.Context, in *{{TableName}}_proto.{{ServerName}}Schema) (res *{{TableName}}_proto.FindOneRes, err error) {
	fmt.Println("CREATE")
	{{EntityName}} := new(entity.{{ServerName}})
	gutil.BeanUtil({{EntityName}}, in)
	{{EntityName}}Res, err := s.{{EntityName}}Model.Insert(*{{EntityName}})
	res = new({{TableName}}_proto.FindOneRes)
	if err != nil {
		res.Code = entity.ERROR
		res.Msg = err.Error()
	} else {
		res.Code = 1000
		res.Msg = "success"
	}
	gutil.BeanUtil(in, &{{EntityName}}Res)
	res.Data = in
	return res, err
}

//更新
func (s *{{ServerName}}Server) Update(ctx context.Context, in *{{TableName}}_proto.UpdateSchema) (*{{TableName}}_proto.FindRes, error) {
	updateForm := new(entity.{{ServerName}}UpdateForm)
	gutil.BeanUtil(&updateForm.Conditions, in.Conditions)
	gutil.BeanUtil(&updateForm.Modifies, in.Modifies)
	aff, err := s.{{EntityName}}Model.Update(updateForm.Conditions, updateForm.Modifies)
	result := new({{TableName}}_proto.FindRes)
	if err != nil {
		result.Code = entity.ERROR
		result.Msg = err.Error()
		result.Data = ""
		return result, err
	}
	result.Code = entity.SUCCESS
	result.Msg = entity.GetResultInfo(entity.SUCCESS)
	result.Data = fmt.Sprintf("affect lines: %d", aff)
	return result, err
}`
