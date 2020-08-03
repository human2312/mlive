package dao

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"

	"mlive/library/config"
)

func init() {
	viper.AddConfigPath("../")
	config.Init()
}

func Test_NewRedis(t *testing.T) {

	redisClient := NewRedis()

	rc := redisClient.Get()

	defer rc.Close()

	t.Log(rc)

	v, err := redis.Int64(rc.Do("incr", "mykey1"))

	t.Log(v, err)

}
