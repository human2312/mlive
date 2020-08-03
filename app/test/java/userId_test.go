package main

import (
	"fmt"
	"mlive/app/test"
	"mlive/service/java"
	"testing"
)

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  userId
 * @Version: X.X.X
 * @Date: 2020/3/24 上午10:58
 */


func Benchmark_Run(t *testing.B)  {
	test.ConfigInit()
	userId,_ := java.Token2Id("5hwipp8nhb0moh4qs13p119f50tulp6o1")
	fmt.Println(userId)
}