package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shiw13/go-one/pkg/logger"
	"github.com/shiw13/go-one/pkg/net/httpx"
	"github.com/shiw13/go-one/pkg/startup"
)

type program struct {
	srv *httpx.Server
}

func main() {
	prg := &program{}

	if err := startup.Run(prg); err != nil {
		log.Fatalf("%s", err)
	}
}

func (p *program) Initialize() error {
	p.srv = httpx.NewServer(
		httpx.WithAddress("0.0.0.0:9443"),
	)
	p.srv.Gin().GET("/test", receiveHandler)

	logger.Infof("Initialize done")
	return nil
}

func (p *program) OnStart() error {
	p.srv.Run()

	logger.Infof("OnStart done")
	return nil
}

func (p *program) OnStop() error {
	logger.Infof("OnStop done")
	return nil
}

func receiveHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
}
