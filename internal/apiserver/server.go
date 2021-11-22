package apiserver

import (
	"fmt"
	"log"

	"github.com/rose839/IAM/internal/apiserver/config"
	genericoptions "github.com/rose839/IAM/internal/pkg/options"
	genericapiserver "github.com/rose839/IAM/internal/pkg/server"
	"github.com/rose839/IAM/pkg/shutdown"
	"github.com/rose839/IAM/pkg/shutdown/shutdownmanagers/posixsignal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

// ExtraConfig defines extra configuration for the iam-apiserver.
type ExtraConfig struct {
	Addr         string // grpc address
	MaxMsgSize   int    // max message size
	ServerCert   genericoptions.GeneratableKeyCert
	MySQLOptions *genericoptions.MySQLOptions
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

func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         fmt.Sprint("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
		ServerCert:   cfg.SecureServing.ServerCert,
		MySQLOptions: cfg.MySQLOptions,
	}, nil
}

// Create iam apiserver instance.
func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	extraConfig.

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

type completedExtraConfig struct {
	*ExtraConfig
}

func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

// New create a grpcAPIServer instance.
func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials: %s", err.Error())
	}

	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
}