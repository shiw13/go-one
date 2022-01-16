package startup

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Service interface {
	Initialize() error
	OnStart() error
	OnStop() error
}

func Run(srv Service, sig ...os.Signal) error {
	if err := srv.Initialize(); err != nil {
		return err
	}

	if err := srv.OnStart(); err != nil {
		return err
	}

	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	s := <-c
	log.Printf("startup: Receive Signal %s", s.String())
	return srv.OnStop()
}
