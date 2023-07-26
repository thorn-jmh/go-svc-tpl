package route

import (
	"github.com/gin-gonic/gin"
	"go-svc-tpl/api/dto"
	"go-svc-tpl/internal/controller"
)

func setupFooController(r *gin.RouterGroup) {
	cw := FooCtlWrapper{
		ctl: controller.NewFooController(),
	}
	p := r.Group("/foo")
	p.GET("/get", cw.GetFoo)
}

type FooCtlWrapper struct {
	ctl controller.IFooController
}

// >>>>>>>>>>>>>>>>>> Controller >>>>>>>>>>>>>>>>>>

// GetFoo godoc
//
//	@Summary		get foo
//	@Description	just get a foo
//	@Tags			foo
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.Resp{data=dto.GetFooResp}
//	@Router			/foo/get [get]
func (w FooCtlWrapper) GetFoo(c *gin.Context) {
	var req dto.GetFooReq
	if err := dto.BindReq(c, &req); err != nil {
		dto.ResponseFail(c, err)
		return
	}
	resp, err := w.ctl.GetFoo(c, &req)
	if err != nil {
		dto.ResponseFail(c, err)
		return
	}
	dto.ResponseSuccess(c, resp)
}
