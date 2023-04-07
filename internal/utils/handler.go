package utils

import (
	errm "tmios/pkg/model/errors"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

type RequestType int
type ReturnType int

var (
	validate *validator.Validate
)

const (
	ReturnTypeJSON ReturnType = iota
	ReturnTypeNone
)

type PreHook func(ctx *ReqContext) error
type PostHook func(ctx *ReqContext, reqArg any, rsp any, err error)

type HandlerAttr struct {
	ReturnType ReturnType
	PreHooks   []PreHook
	PostHooks  []PostHook
}

type HandlerOption func(*HandlerAttr)

func newHandlerAttr(opts ...HandlerOption) *HandlerAttr {
	attr := HandlerAttr{}
	for _, o := range opts {
		o(&attr)
	}

	return &attr
}

type ReqContext struct {
	Gin  *gin.Context
	Data map[string]interface{}
}

func init() {
	validate = validator.New()
}

func isRedisNil(err error) bool {
	return err.Error() == "redis: nil"
}

func Handler[T any](
	handlerFunc func(ctx *ReqContext, req *T) (interface{}, error),
	opts ...HandlerOption,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			attr   = newHandlerAttr(opts...)
			reqArg T
		)
		switch c.Request.Method {
		case "POST":
			if err := c.ShouldBind(&reqArg); err != nil {
				ApiErr(c, errm.ErrParam.SetDetail("%s", err.Error()))
				return
			}
		case "GET":
			if err := c.ShouldBindQuery(&reqArg); err != nil {
				ApiErr(c, errm.ErrParam.SetDetail("%s", err.Error()))
				return
			}
		}

		ctxt := &ReqContext{
			Gin:  c,
			Data: make(map[string]interface{}),
		}

		// Check request arguments
		if c.Request.Method == "POST" {
			if err := validate.Struct(&reqArg); err != nil {
				ApiErr(c, errm.ErrParam.SetDetail("%s", err.Error()))
				return
			}
		}

		for _, hookFunc := range attr.PreHooks {
			if err := hookFunc(ctxt); err != nil {
				ApiErr(c, err)
				return
			}
		}

		rsp, err := handlerFunc(ctxt, &reqArg)

		for _, postHook := range attr.PostHooks {
			postHook(ctxt, &reqArg, rsp, err)
		}

		switch attr.ReturnType {
		case ReturnTypeJSON:
			if err != nil {
				ApiErr(c, err)
			} else {
				ApiRet(c, rsp)
			}

		case ReturnTypeNone:
		}
	}
}
