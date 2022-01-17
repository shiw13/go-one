package main

import (
	"log"

	"github.com/shiw13/go-one/pkg/logger"
	"github.com/shiw13/go-one/pkg/startup"
)

type program struct{}

func main() {
	prg := &program{}

	if err := startup.Run(prg); err != nil {
		log.Fatalf("%s", err)
	}
}

func (p *program) Initialize() error {
	logger.Infof("Initialize done")
	return nil
}

func (p *program) OnStart() error {
	logger.Infof("OnStart done")
	return nil
}

func (p *program) OnStop() error {
	logger.Infof("OnStop done")
	return nil
}
