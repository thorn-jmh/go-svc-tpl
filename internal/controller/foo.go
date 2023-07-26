package controller

import (
	"github.com/gin-gonic/gin"
	"go-svc-tpl/api/dto"
	"go-svc-tpl/internal/dao"
	"go-svc-tpl/internal/dao/model"
)

// >>>>>>>>>>>>>>>>>> Interface  >>>>>>>>>>>>>>>>>>

type IFooController interface {
	GetFoo(*gin.Context, *dto.GetFooReq) (*dto.GetFooResp, error)
}

// >>>>>>>>>>>>>>>>>> Controller >>>>>>>>>>>>>>>>>>

// check interface implementation
var _ IFooController = (*FooController)(nil)

var NewFooController = func() *FooController {
	return &FooController{}
}

type FooController struct {
	// maybe some logic config to read from viper
	// or a service dependency
}

// ---------------------- GetFoo ----------------------

func (c *FooController) GetFoo(ctx *gin.Context, req *dto.GetFooReq) (*dto.GetFooResp, error) {
	var resp dto.GetFooResp
	dao.DB(ctx).Model(model.Foo{Name: req.Name}).First(&resp)
	return &resp, nil
}
