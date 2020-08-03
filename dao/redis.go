package dao

import (
    "time"

    "github.com/garyburd/redigo/redis"
    "github.com/spf13/viper"
)

func NewRedis() *redis.Pool {
    redisClient := &redis.Pool{
        MaxIdle:     viper.GetInt("redis.maxIdle"),
        MaxActive:   viper.GetInt("redis.maxActive"),
        IdleTimeout: viper.GetDuration("redis.idleTimeout") * time.Second,
        Wait:        true,
        Dial: func() (redis.Conn, error) {
            con, err := redis.Dial("tcp", viper.GetString("redis.host")+":"+viper.GetString("redis.port"),
                redis.DialPassword(viper.GetString("redis.passWord")),
                redis.DialConnectTimeout(viper.GetDuration("redis.connectTimeout")*time.Second),
                redis.DialReadTimeout(viper.GetDuration("redis.readTimeout")*time.Second),
                redis.DialWriteTimeout(viper.GetDuration("redis.writeTimeout")*time.Second))
            if err != nil {
                panic("redis connect failed" + err.Error())
            }
            _, err = con.Do("ping")
            if err != nil {
                panic("redis ping failed" + err.Error())
            }
            return con, nil
        },
    }
    return redisClient
}
