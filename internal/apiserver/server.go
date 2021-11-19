package apiserver

import (
	"log"

	"github.com/rose839/IAM/internal/apiserver/config"
	genericapiserver "github.com/rose839/IAM/internal/pkg/server"
	"github.com/rose839/IAM/pkg/shutdown"
	"github.com/rose839/IAM/pkg/shutdown/shutdownmanagers/posixsignal"
)

// apiServer represent iam apiserver runtime instance.
type apiServer struct {
	gs               *shutdown.GracefulShutdown         // graceful shutdown instance
	genericAPIServer *genericapiserver.GenericAPIServer // rest api server
	gRPCAPIServer    *grpcAPIServer                     // grpc server
}

// preparedAPIServer represent an iam apiserver runtime instance that is prepared.
type preparedAPIServer struct {
	*apiServer
}

// Create rest api server config from app config.
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
	return
}

// Create iam apiserver instance.
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
	// init rest api server router
	initRouter(s.genericAPIServer.Engine)

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
