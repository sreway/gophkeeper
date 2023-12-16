package app

import (
	"context"
	"net"
)

func (s *server) Run(ctx context.Context) error {
	ctxServer, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	listen, err := net.Listen("tcp", s.config.Host)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		s.gRPCServer.GracefulStop()
		stopServer()
	}()

	err = s.gRPCServer.Serve(listen)
	if err != nil {
		return err
	}

	<-ctxServer.Done()
	return nil
}
