package main

import (
	"fmt"
	"mlive/app/test"
	"mlive/service/java"
)

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  userId
 * @Version: X.X.X
 * @Date: 2020/3/24 上午10:58
 */


func main()  {
	test.ConfigInit()

	userId,e := java.Token2Id("5hwipp8nhb0moh4qs13p119f50tulp6o1")
	//userId,e := admin.GetAdminUserInfo("6hfhn3xdli2ufpofi86k20j0yj60wowj1",2)
	if e != nil {
		//success...
	}
	fmt.Println(userId,e.Msg)
}
