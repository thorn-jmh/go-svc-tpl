package dto

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/sirupsen/logrus"
    "go-svc-tpl/utils/stacktrace"
    "net/http"
)

// >>>>>>>>>>>>>>> Response >>>>>>>>>>>>>>>>>>

// Resp
//	@name			Response
//	@Description	Resp is the general response struct
type Resp struct {
    Code int    `json:"code" example:"200"`     // status code
    Msg  string `json:"msg"  example:"success"` // message
    Data any    `json:"data"`                   // data payload
}

func Response(ctx *gin.Context, code int, msg string, data any) {
    ctx.JSON(http.StatusOK, Resp{
        Code: code,
        Msg:  msg,
        Data: data,
    })
}

func ResponseSuccess(ctx *gin.Context, data any) {
    Response(ctx, http.StatusOK, "success", data)
}

func ResponseFail(ctx *gin.Context, err error) {
    logrus.Error(err)

    code := stacktrace.GetCode(err)
    if code == stacktrace.NoCode {
        code = http.StatusInternalServerError
    }
    msg := stacktrace.Current(err).Error()
    Response(ctx, int(code), msg, nil)

}

// <<<<<<<<<<<<<<<<< Response >>>>>>>>>>>>>>>>>

// >>>>>>>>>>>>>>>> Request >>>>>>>>>>>>>>>>>>

// BindReq bind request, and log info&err.
// Support json & query & header,
// if there is a conflict, the former will overwrite the latter.
func BindReq[T any](c *gin.Context, req *T) error {
    logrus.Debugf("Req.Url: %s, Req.Body: %+v", c.Request.URL, c.Request.Body)
    if err := c.ShouldBindWith(req, GENERAL_BINDER); err != nil {
        return stacktrace.PropagateWithCode(err, http.StatusBadRequest, "Failed to bind request")
    }
    return nil
}

var (
    GENERAL_BINDER = General{}
)

type General struct{}

func (General) Name() string {
    return "general"
}

func (General) Bind(r *http.Request, obj any) error {
    // general binder
    // 1. bind header
    // 2. bind query or form
    // 3. bind body (json)
    _ = binding.Header.Bind(r, obj)
    _ = binding.Query.Bind(r, obj)
    _ = binding.JSON.Bind(r, obj)
    return binding.Validator.ValidateStruct(obj)
}

// <<<<<<<<<<<<<<<<< Request <<<<<<<<<<<<<<<<<<
