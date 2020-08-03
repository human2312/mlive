package dao
// @Time : 2020年01月02日14:08:16
// @Author : Ray
// @Desc: arangodb 数据库

import (
	"context"
	"crypto/tls"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/spf13/viper"
	"mlive/library/logger"
)


//var ArangoDb driver.Database

//func ConnArangoDb()  {
//	var err error
//	ArangoDb,err = connDb()
//	if err != nil {
//		logger.Eprintln("链接失败:",err)
//		log.Fatal("链接失败:",err)
//	}
//}

// 链接db
func NewArangoDb()(driver.Database) {

	userName := viper.GetString("arangodb.userName")
	passWord := viper.GetString("arangodb.passWord")
	dataBaseName := viper.GetString("arangodb.dataBaseName")
	url := viper.GetString("arangodb.url")
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{url},
		TLSConfig: &tls.Config{ /*...*/ },
	})
	if err != nil {
		// Handle error
		logger.Eprintln("conn db err:",err)
		return nil
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(userName, passWord),
	})
	if err != nil {
		// Handle error
		logger.Eprintln("conn new client err:", err)
		return nil
	}
	ctx := context.Background()
	arangoDb, err := c.Database(ctx, dataBaseName)
	if err != nil {
		// handle error
		logger.Eprintln("conn database err:", err)
		return nil
	}
	return arangoDb
}