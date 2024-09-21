package mredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Client *redis.Client
var ctx = context.Background()

func NewRedis(addr string, password string, db int) {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		//fmt.Println("redis 连接失败")
		panic("redis 连接失败")
	}
}

func Get(key string) (val string) {
	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return
}

func GetToStruct(key string) {
	var value interface{}
	err := Client.Get(context.Background(), key).Scan(value)
	if err != nil {
		panic(err)
	}
}

func Set(key string, value interface{}, expiration time.Duration) {
	err := Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
	return
}

// SAdd 集合操作
func SAdd(key string, value interface{}) int64 {
	res, err := Client.SAdd(ctx, key, value).Result()
	if err != nil {
		panic(err)
	}
	return res
}

func SMembers(key string) (val []string) {
	val, err := Client.SMembers(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return
}

func SRem(key string, members ...interface{}) {
	err := Client.SRem(ctx, key, members).Err()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return
}

func SPop(key string) (val string) {
	val, err := Client.SPop(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return
}

func SPopN(key string, count int64) (val []string) {
	val, err := Client.SPopN(ctx, key, count).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return
}

func SScan(key string, cursor uint64, match string, count int64) (val []string) {
	val, _, err := Client.SScan(ctx, key, cursor, match, count).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return
}
