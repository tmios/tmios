package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"tmios/internal/config"
)

var defaultListen string

// Api api
type Api struct {
	Router     *gin.Engine
	cnf        *config.Config
	listenAddr string
}
type Option func(*Api)

func NewHttp(opts ...Option) *Api {
	g := gin.Default()
	g.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.New(cors.Config{
			AllowMethods: []string{"OPTIONS", "POST", "GET"},
			AllowHeaders: []string{"Origin", "X-Requested-With",
				"Content-Type", "Accept", "X-TOKEN"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		}))

	http := &Api{
		Router:     g,
		cnf:        config.NewConfig(),
		listenAddr: defaultListen,
	}
	for _, opt := range opts {
		opt(http)
	}
	return http
}

func (a *Api) Run() error {
	if a.listenAddr == defaultListen {
		a.listenAddr = a.cnf.Conf.API.ListenAddr
		err := a.Router.Run(a.listenAddr)
		if err != nil {
			return err
		}
	} else {
		go func() {
			err := a.Router.Run(a.listenAddr)
			if err != nil {
				logrus.Fatal(err)
			}
		}()
	}

	return nil
}

type PageReq struct {
	PageIndex int                    `json:"page_index" validate:"required"`
	PageSize  int                    `json:"page_size" validate:"required"`
	Query     map[string]interface{} `json:"query"`
}

type PageResp struct {
	Page  int   `json:"page"`
	Total int64 `json:"total"`
	Data  any   `json:"data"`
}
