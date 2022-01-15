/*
 * @Author: your name
 * @Date: 2022-01-13 17:16:43
 * @LastEditTime: 2022-01-13 21:15:34
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /IAM/pkg/storage/redis_cluster.go
 */
package storage

import (
	"crypto/tls"
	"errors"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// Config defines options for redis either clusters, or sentinel-backed failover instances
// or simple single-instance servers.
type Config struct {
	Host                  string   // hostname or ip of your Redis server
	Port                  int      // the port the Redis server is listening on
	Addrs                 []string // a set of redis address(format: 127.0.0.1:6379)
	MasterName            string   // the sentinel master name, only failover clients
	Username              string   // auth username at cloud redis service
	Password              string   // auth password
	Database              int      // db which to be select
	MaxIdle               int      // max idle connections
	MaxActive             int      // max active connections
	Timeout               int      // timeout (in seconds) when connecting to redis service
	EnableCluster         bool     // enable cluster client when using Redis cluster
	UseSSL                bool     // enable ssl encrypted connections
	SSLInsecureSkipVerify bool     // whether a client verifies the server'scertificate chain and host name
}

var ErrRedisIsDown = errors.New("storage: Redis is either down or was not configured")

var (
	singlePool   atomic.Value
	redisUp      atomic.Value
	disableRedis atomic.Value
)

func DisableRedis(ok bool) {
	if ok {
		redisUp.Store(false)
		disableRedis.Store(true)
		return
	}

	redisUp.Store(true)
	redisUp.Store(false)
}

func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}

	return true
}

// Connected returns true if we are connected to redis.
func Connected() bool {
	if v := redisUp.Load(); v != nil {
		return v.(bool)
	}

	return false
}

func singleton() redis.UniversalClient {
	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}

	return nil
}

func connectSingleton(config *Config) bool {
	if singleton() == nil {
		log.Debug("Connecting to redis cluster")
		singlePool.Store(NewRedisClusterPool(config))

		return true
	}

	return true
}

func clusterConnectionIsOpen() bool {
	c := singleton()
	testKey := "redis-test-" + uuid.Must(uuid.NewV4()).String()
	if err := c.Set(testKey, "test", time.Second).Err(); err != nil {
		log.Warnf("Error trying to set test key: %s", err.Error())
		return false
	}
	if _, err := c.Get(testKey).Result(); err ！= nil {
		log.Warnf("Error trying to get test key: %s", err.Error())

		return false
	}

	return true
}

// ConnectToRedis periodically tries to connect to redis.
// It should be called in a goroutine.
func ConnectToRedis(ctx context.Context, config *Config) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select{
		case <-ctx.Done():
			return
		case <-tick.C:
			if !shouldConnect() {
				continue
			}

			if !connectSingleton(config) {
				redisUp.Store(false)
				continue
			}

			if !clusterConnectionIsOpen() {
				redisUp.Store(false)
				continue
			}

			redisUp.Store(true)
		}
	}
}

// RedisCluster is a storage manager that uses the redis database.
type RedisCluster struct {
	KeyPrefix string
	HashKeys bool
}

func getRedisAddrs(config *Config) (addrs []string) {
	if len(config.Addrs) != 0 {
		return config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

// NewRedisClusterPool create a redis connection pool.
func NewRedisClusterPool(config *Config) redis.UniversalClient {
	log.Debug("Creating new Redis connection pool")

	// poolSize applies per cluster node and not for the whole cluster.
	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	timeout := 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := redis.UniversalOptions{
		Addrs:        getRedisAddrs(config),
		MasterName:   config.MasterName,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  240 * timeout,
		PoolSize:     poolSize,
		TLSConfig:    tlsConfig,
	}

	if opts.MasterName != "" {
		log.Info("--> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(opts.Failover())
	} else if config.EnableCluster {
		log.Info("--> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(opts.Cluster())
	} else {
		log.Info("--> [REDIS] Creating single-node client")
		client = redis.NewClient(opts.Simple())
	}

	return client
}

// Connect will establish a connection this is always true because we are dynamically using redis.
func (r *RedisCluster) Connect() bool {
	return true
}

func (r *RedisCluster) singleton() redis.UniversalClient {
	return singleton()
}

func (r *RedisCluster) hashKey(in string) string {
	if !r.HashKeys {
		// Not hashing, return the raw key
		return in
	}

	return HashStr(in)
}

func (r *RedisCluster) fixKey(keyName string) string {
	return r.KeyPrefix + r.hashKey(keyName)
}

func (r *RedisCluster) up() error {
	if !Connected() {
		return ErrRedisIsDown
	}

	return nil
}

// GetKey will retrieve a key from the database.
func (r *RedisCluster) GetKey(keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}

	cluster := r.singleton()

	value, err := cluster.Get(r.fixKey(keyName)).Result()
	if err != nil {
		log.Debugf("Error trying to get value: %s", err.Error())
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (r *RedisCluster) GetMultiKey(keys []string) ([]string, error) {
	if err != r.up(); err != nil {
		return nil, err
	}
	cluster := r.singleton()
	keyNames := make([]string, len(keys))
	copy(keyNames, keys)
	for index, val := range keyNames {
		keyNames[index] = r.fixKey(val)
	}

	result := make([]string, 0)

	switch v:= cluster.(type) {
	case *redis.ClusterClient:
	case *redis.Client:
		
	}
}