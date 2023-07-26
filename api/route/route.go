package route

import (
	"github.com/gin-gonic/gin"
	"go-svc-tpl/api/dto"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Msg:  "success",
		Data: "pong~",
	})
}

func SetupRouter(r *gin.RouterGroup) {
	r.GET("/ping", Ping)
	setupFooController(r)
}
