package clear

import (
	d "mlive/dao"
	"errors"
	"context"
)

type ClearDb struct {

}

// 清除 用户数据
func (c *ClearDb)ToClear()(bool,error)  {

	if err := d.Db.DB.Exec("truncate mpm_user").Error; err != nil{
		return false,errors.New("清空mysql user表失败")
	}
	status,err := c.clearArangoDb("mpm_user")
	if !status || err != nil{
		return false,errors.New("清空arangodb user表 Truncate 失败")
	}
	status1,err := c.clearArangoDb("mpm_user_layer_of")
	if !status1 || err != nil{
		return false,errors.New("清空arangodb mpm_user_layer_of表 Truncate 失败")
	}
	return true,nil
}

func (c *ClearDb)clearArangoDb(TableName string)(bool,error)  {
	ctx := context.Background()
	col,err := d.Db.ArangoDb.Collection(ctx,TableName)
	if err != nil{
		return false,errors.New("清空arangodb user表失败")
	}
	err = col.Truncate(ctx)
	if err  != nil{
		return false,errors.New("清空arangodb user表 Truncate 失败")
	}
	return true,nil
}