package app

import (
	"google.golang.org/grpc/status"
)

func grpcGetStatusCode(err error) string {
	s, ok := status.FromError(err)
	if ok {
		return s.Code().String()
	}
	return ""
}
