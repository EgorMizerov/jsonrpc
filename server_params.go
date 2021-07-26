package jsonrpc

import "crypto/tls"

type ServerParam struct {
	fn func(s *Server) error
}

func Addr(a string) ServerParam {
	return ServerParam{
		func(s *Server) error {
			s.server.Addr = a
			return nil
		},
	}
}

func TLSConfig(t *tls.Config) ServerParam {
	return ServerParam{
		func(s *Server) error {
			s.server.TLSConfig = t
			return nil
		},
	}
}
