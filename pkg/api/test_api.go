package api

import (
	"github.com/gin-gonic/gin"
	http2 "net/http"
	"tmios/internal/http"
)

func WithTest() http.Option {
	return func(api *http.Api) {
		group := api.Router.Group("api/v1/test")
		group.GET("/get", func(ctx *gin.Context) {
			ctx.JSON(http2.StatusOK, "success")
		})
	}
}
