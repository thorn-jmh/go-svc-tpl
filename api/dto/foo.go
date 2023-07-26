package dto

import "go-svc-tpl/internal/dao/model"

type GetFooReq struct {
	Name string `json:"name"`
}

type GetFooResp struct {
	model.Foo
}
