package apiserver

import (
	"context"
	"fmt"
	"log"

	pb "github.com/rose839/IAM/api/proto/apiserver/v1"
	"github.com/rose839/IAM/internal/apiserver/config"
	cachev1 "github.com/rose839/IAM/internal/apiserver/controller/v1/cache"
	"github.com/rose839/IAM/internal/apiserver/store"
	"github.com/rose839/IAM/internal/apiserver/store/mysql"
	genericoptions "github.com/rose839/IAM/internal/pkg/options"
	genericapiserver "github.com/rose839/IAM/internal/pkg/server"
	"github.com/rose839/IAM/pkg/shutdown"
	"github.com/rose839/IAM/pkg/shutdown/shutdownmanagers/posixsignal"
	"github.com/rose839/IAM/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// apiServer represent iam apiserver runtime instance.
type apiServer struct {
	gs               *shutdown.GracefulShutdown         // graceful shutdown instance
	redisOptions     *genericoptions.RedisOptions       // redis options
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
	MaxMsgSize   int    // grpc max message size
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
		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
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

	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:               gs,
		redisOptions:     cfg.RedisOptions,
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
	}

	return server, nil
}

// do some prepare work on apiserver object.
func (s *apiServer) PrepareRun() preparedAPIServer {
	// init rest api server router
	initRouter(s.genericAPIServer.Engine)

	// init redis connection
	s.initRedisStore()

	// add graceful shutdown callback
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		// close mysql connection
		mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
		if mysqlStore != nil {
			return mysqlStore.Close()
		}

		// close rest api server and grpc server
		s.genericAPIServer.Close()
		s.gRPCAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

// start api server
func (s preparedAPIServer) Run() error {
	// start rpc server
	go s.gRPCAPIServer.Run()

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
	// create grpc encryption communication
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials: %s", err.Error())
	}

	// create grpc server
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	// create mysql store instance
	storeIns, _ := mysql.GetMySQLFactoryOr(c.MySQLOptions)
	store.SetClient(storeIns)

	// create grpc instance
	cacheIns, err := cachev1.GetCacheInsOr(storeIns)
	if err != nil {

	}

	// register grpc server
	pb.RegisterCacheServer(grpcServer, cacheIns)

	// register grpc reflcetion server to enable grpcurl and grpc Swagger test
	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
}

func (s *apiServer) initRedisStore() {
	ctx, cancle := context.WithCancel(context.Background())
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		cancle()

		return nil
	}))

	config := &storage.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}

	// try to connect to redis
	go storage.ConnectToRedis(ctx, config)
}
