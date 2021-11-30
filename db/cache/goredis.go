package cache

import (
	"context"
	"fmt"
	"github.com/ehwjh2010/cobra/client"
	"github.com/ehwjh2010/cobra/log"
	"github.com/go-redis/redis/v8"
	"time"
)

func InitCacheWithGoRedis(redisConfig *client.Cache) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Network:            network,
		Addr:               fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Username:           redisConfig.User,
		Password:           redisConfig.Pwd,
		DB:                 redisConfig.Database,
		DialTimeout:        time.Duration(redisConfig.ConnectTimeout) * time.Second,
		ReadTimeout:        time.Duration(redisConfig.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(redisConfig.WriteTimeout) * time.Second,
		PoolFIFO:           true,
		PoolSize:           redisConfig.MaxOpenConnCount,
		MinIdleConns:       redisConfig.MinFreeConnCount,
		IdleTimeout:        time.Duration(redisConfig.FreeMaxLifetime) * time.Minute,
	})

	ctx := context.Background()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	log.Info("connect redis success")
	return redisClient, nil
}
