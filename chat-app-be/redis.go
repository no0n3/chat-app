package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisCli struct {
	cli *redis.Client
}

var REDIS_CLI *RedisCli

func NewRedisCli() *RedisCli {
	rcli := &RedisCli{}

	rcli.cli = redis.NewClient(&redis.Options{
		// Addr: "localhost:6378",
		Addr: "redis:6379",
	})

	return rcli
}

func initRedisCli() {
	REDIS_CLI = NewRedisCli()
}

func (cli *RedisCli) UserIpNodeAdd(userId string, ip string) {
	cli.cli.SAdd(ctx, userWsIpsKey(userId), ip)
}

func (cli *RedisCli) SetRemove(userId string, ip string) {
	cli.cli.SRem(ctx, userWsIpsKey(userId), ip)
}

func (cli *RedisCli) GetUserWsIps(userId string) []string {
	return cli.cli.SMembers(ctx, userWsIpsKey(userId)).Val()
}

func userWsIpsKey(userId string) string {
	return userId + "-ws-conn-ips"
}
