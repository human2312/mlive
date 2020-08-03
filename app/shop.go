package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mlive/service/shop"
)

func ShopList(c *gin.Context)  {
	res := shop.List()
	fmt.Println(res)
	return
}

func ShopAdd(c *gin.Context)  {
	shop.Add()
}

func ShopEdit(c *gin.Context)  {
	shop.Edit()
}

func ShopDel(c *gin.Context)  {
	shop.Del()
}
