package arango

import (
	"github.com/arangodb/go-driver"
	"strconv"
	d "mlive/dao"
	"context"
	"errors"
)



type  MyUserTeamOf struct {
	Key     string  `json:"_key"`
	To		string  `json:"_to"`
	From	string  `json:"_from"`
}

func (arango *ArangoUser)TableTeamOf()string  {
	table := "mlive_user_team_of"
	return table
}

func (arango *ArangoUser)InfoTeamOf(userId int64)(MyUserTeamOf,error)  {
	arangoDb := d.Db.ArangoDb
	ctx 	 := context.Background()
	quy	 	 := "FOR d in @@collection FILTER d._key == @uid  return d "
	bindVars := map[string]interface{}{
		"@collection":arango.TableTeamOf(),
		"uid":strconv.FormatInt(userId,10),
	}
	cursor,err := arangoDb.Query(ctx,quy,bindVars)
	if err != nil {
		return MyUserTeamOf{},err
	}
	defer cursor.Close()
	var doc MyUserTeamOf
	for {
		_,err := cursor.ReadDocument(ctx,&doc)
		if driver.IsNoMoreDocuments(err){
			break
		}else if err != nil {
			//log.Fatal("info read document err:",err)
			return doc,err
		}
	}
	return doc,nil
}


// 创建邀请-团队关联
// @ fromId childId
// @ toId	parentId
func (arango *ArangoUser)CreateTeamOf(fromId int64,toId int64)(int,error)  {
	arangoDb 		:= d.Db.ArangoDb
	userId          := strconv.FormatInt(fromId,10)
	fromIdString 	:= arango.Table()+"/"+userId
	toIdString 		:= arango.Table()+"/"+strconv.FormatInt(toId,10)
	userDocument := MyUserTeamOf{
		Key:userId,
		From:fromIdString,
		To:toIdString,
	}
	ctx := context.Background()
	_,err := arango.CheckCollection(ctx,2,arango.TableTeamOf())
	if err != nil{
		return 0,err
	}
	col,err := arangoDb.Collection(ctx,arango.TableTeamOf())
	if err != nil {
		return 0,err
	}
	meta,err := col.CreateDocument(ctx,userDocument)
	if err != nil{
		return 0,err
	}
	key,_ :=  strconv.Atoi(meta.Key)
	return key,nil
}

// 删除团队关联
func (arango *ArangoUser)DeleteTeamOf(userId int64)(bool,error)  {
	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	col,err := arangoDb.Collection(ctx,arango.TableTeamOf())
	if err != nil {
		return false,err
	}
	key := strconv.FormatInt(userId,10)
	_,err = col.RemoveDocument(ctx, key)
	if err != nil {
		return false,err
	}
	return true,nil
}

// 更改团队关联
func (arango *ArangoUser)UpdateTeamOf(userId int64,patch  map[string]interface{})(string,error)  {
	uidString := strconv.FormatInt(userId,10)
	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	//patch := map[string]interface{}{
	//	"status":1,
	//}
	col,err := arangoDb.Collection(ctx,arango.TableTeamOf())
	if err != nil {
		return "",errors.New("update collection err:"+err.Error())
	}
	meta,err := col.UpdateDocument(ctx,uidString,patch)
	if err != nil {
		return "",errors.New("update collection err:"+err.Error())
	}
	return meta.Key,nil
}