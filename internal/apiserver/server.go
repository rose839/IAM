package apiserver

import (
	"log"

	"github.com/rose839/IAM/internal/apiserver/config"
	genericapiserver "github.com/rose839/IAM/internal/pkg/server"
	"github.com/rose839/IAM/pkg/shutdown"
	"github.com/rose839/IAM/pkg/shutdown/shutdownmanagers/posixsignal"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericapiserver.GenericAPIServer
	gRPCAPIServer    *grpcAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, err error) {
	genericConfig = genericapiserver.NewConfig()

	if err = cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return
	}

	if err = cfg.FeatureOptions.ApplyTo(genericConfig); err != nil {
		return
	}

	if err = cfg.SecureServing.ApplyTo(genericConfig); err != nil {
		return
	}

	if err = cfg.InsecureServing.ApplyTo(genericConfig); err != nil {
		return
	}
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
	}

	return server, nil
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {

		// close rest api server
		s.genericAPIServer.Close()

		return nil
	}))
	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// this will block until api server close
	s.genericAPIServer.Run()

	return nil
}
