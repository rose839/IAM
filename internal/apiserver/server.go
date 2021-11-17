package apiserver

import (
	genericapiserver "github.com/rose839/IAM/internal/pkg/server"
	"github.com/rose839/IAM/pkg/shutdown"
)

type apiserver struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericapiserver.GenericAPIServer
}

type preparedAPIServer struct {
	*apiserver
}

