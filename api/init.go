package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-svc-tpl/api/route"
	"net/http"
	"time"
)

//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag fmt -d ../ -g api/init.go
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag init -d ../ -g api/init.go --ot yaml -o ../docs

//	@title			go-svc-tpl API
//	@version		0.0.1
//	@description	A simple go service template, which is used to build a go service quickly.
//	@BasePath		/api

type WebServerCfg struct {
	Port         int `mapstructure:"Port"`
	WriteTimeout int `mapstructure:"WriteTimeout"`
	ReadTimeout  int `mapstructure:"ReadTimeout"`
}

func StartServer() error {
	var cfg WebServerCfg
	if err := viper.Sub("WebServer").UnmarshalExact(&cfg); err != nil {
		return err
	}

	e := gin.Default()
	route.SetupRouter(e.Group("/api"))

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        e,
		ReadTimeout:    time.Second * time.Duration(cfg.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(cfg.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}
