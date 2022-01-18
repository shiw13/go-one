package httpx

import (
	"crypto/tls"

	"github.com/gin-gonic/gin"
	"github.com/shiw13/go-one/pkg/logger"
)

type ServerOption func(o *Server)

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithTLS(certPath, keyPath string) ServerOption {
	return func(s *Server) {
		s.tlsEnabled = true
		s.certPath = certPath
		s.keyPath = keyPath
	}
}

type Server struct {
	addr       string
	tlsEnabled bool
	certPath   string
	keyPath    string
	engine     *gin.Engine
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{}

	for _, o := range opts {
		o(srv)
	}

	gin.SetMode(gin.ReleaseMode)
	srv.engine = gin.New()

	return srv
}

func (s *Server) Run() {
	if s.tlsEnabled {
		if _, err := tls.LoadX509KeyPair(s.certPath, s.keyPath); err != nil {
			logger.Fatalf("HTTPS server: invalid cert")
		}

		go func() {
			if err := s.engine.RunTLS(s.addr, s.certPath, s.keyPath); err != nil {
				logger.Fatalf("HTTPS server: run failed. error: %v", err)
			}
			logger.Infof("HTTPS server: running on %s", s.addr)
		}()
		return
	}

	go func() {
		if err := s.engine.Run(s.addr); err != nil {
			logger.Fatalf("HTTP server: run failed. error: %v", err)
		}
		logger.Infof("HTTPS server: running on %s", s.addr)
	}()
}

func (s *Server) Gin() *gin.Engine {
	return s.engine
}
