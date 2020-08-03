package arango

import (
	"context"
	"errors"
	"github.com/arangodb/go-driver"
	"log"
	"math"
	d "mlive/dao"
	"sort"
	"strconv"
	"time"
)



type ArangoUser struct {

}

type MyUserInfo struct {
	AId  	 	 string 		`json:"_id"`
	AKey 	 	 string 		`json:"_key"`
	Id 			 int64			`json:"id"`
	InviteUpId	 int64			`json:"invite_up_id"`
	InviteTime   string			`json:"invite_time"`
	UserName	 string			`json:"user_name"`
	HeadImgUrl	 string			`json:"head_img_url"`
	Nickname	 string			`json:"nickname"`
	Name    	 string			`json:"name"`
	Gender		 int64			`json:"gender"`
	Level		 int64			`json:"level"`
	InviteCode	 string			`json:"invite_code"`
	Mobile	 	 string			`json:"mobile"`
	TelPhone	 string			`json:"tel_phone"`
	Email		 string			`json:"email"`
	IsCompany  	 int64			`json:"is_company"`
	Status		 int64			`json:"status"`
	Operator	 string			`json:"operator"`
	AdminUserId  int64			`json:"admin_user_id"`
	CreateTime   string			`json:"create_time"`
	UpdateTime   string			`json:"update_time"`

}

type MyUserList struct {
	Id 			 int64			`json:"id"`
	InviteUpId	 int64			`json:"invite_up_id"`
	InviteTime   string			`json:"invite_time"`
	UserName	 string			`json:"user_name"`
	HeadImgUrl	 string			`json:"head_img_url"`
	Nickname	 string			`json:"nickname"`
	Name    	 string			`json:"name"`
	Gender		 int64			`json:"gender"`
	Level		 int64			`json:"level"`
	InviteCode	 string			`json:"invite_code"`
	Mobile	 	 string			`json:"mobile"`
	TelPhone	 string			`json:"tel_phone"`
	Email		 string			`json:"email"`
	IsCompany  	 int64			`json:"is_company"`
	Status		 int64			`json:"status"`
	Operator	 string			`json:"operator"`
	AdminUserId  int64			`json:"admin_user_id"`
	CreateTime   string			`json:"create_time"`
	UpdateTime   string			`json:"update_time"`
	ChildList 	 []*MyUserList	`json:"childList"`
}

func (arango *ArangoUser)Table()string  {
	table	 := "mlive_user"
	return table
}

//// 获取用户信息
func (arango *ArangoUser)Info(userId int64)(MyUserInfo,error)  {
	arangoDb := d.Db.ArangoDb
	ctx 	 := context.Background()
	quy	 	 := "FOR d in @@collection FILTER d._key == @uid   return d "
	bindVars := map[string]interface{}{
		"@collection":arango.Table(),
		"uid":strconv.FormatInt(userId,10),
	}
	cursor,err := arangoDb.Query(ctx,quy,bindVars)
	if err != nil {
		return MyUserInfo{},err
	}
	defer cursor.Close()
	var doc MyUserInfo
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



//获取分润上级的信息
func (arango *ArangoUser)GetMoneySuperior(userIdInt int64) (*MyUserInfo,error)  {
	arangoDb := d.Db.ArangoDb
	myUserInfo, err := arango.Info(userIdInt)
	if err != nil {
		return &MyUserInfo{}, err
	}
	if myUserInfo.Id <= 0 {
		return &MyUserInfo{}, err
	}
	level := strconv.FormatInt(myUserInfo.Level,10)
	userId := strconv.FormatInt(userIdInt,10)
	ctx := context.Background()
	//  and  v.status==1 and v.deleted==0
	query := "FOR v IN 1..99999999 OUTBOUND '"+arango.Table()+"/" + userId + "' "+arango.TableUserOf()+" FILTER v.level > "+level+" and  v.status==1  Limit 1 RETURN v"
	cursor, err := arangoDb.Query(ctx, query, nil)
	if err != nil {
		// handle error
		return &MyUserInfo{}, err
	}
	defer cursor.Close()
	var data *MyUserInfo
	_, err = cursor.ReadDocument(ctx, &data)
	if driver.IsNoMoreDocuments(err) {
	} else if err != nil {
		// handle other errors
		return &MyUserInfo{}, err
	}
	if data == nil {
		return &MyUserInfo{}, nil
	}
	return data, nil
}

// 邀请同级- 用户id,指定的同级
func (arango *ArangoUser)GetInviteSameLevel(userIdInt int64,iLevel int64) (*MyUserInfo,error)  {
	arangoDb := d.Db.ArangoDb
	myUserInfo, err := arango.Info(userIdInt)
	if err != nil {
		return &MyUserInfo{}, err
	}
	if myUserInfo.Id <= 0 {
		return &MyUserInfo{}, err
	}
	level := strconv.FormatInt(iLevel,10)
	userId := strconv.FormatInt(userIdInt,10)
	ctx := context.Background()
	//  and  v.status==1 and v.deleted==0
	query := "FOR v IN 1..99999999 OUTBOUND '"+arango.Table()+"/" + userId + "' "+arango.TableUserOf()+" FILTER v.level == "+level+" and  v.status==1  Limit 1 RETURN v"
	cursor, err := arangoDb.Query(ctx, query, nil)
	if err != nil {
		// handle error
		return &MyUserInfo{}, err
	}
	defer cursor.Close()
	var data *MyUserInfo
	_, err = cursor.ReadDocument(ctx, &data)
	if driver.IsNoMoreDocuments(err) {
	} else if err != nil {
		// handle other errors
		return &MyUserInfo{}, err
	}
	if data == nil {
		return &MyUserInfo{}, nil
	}
	return data, nil
}


// 同级奖励:  1 代
func  (arango *ArangoUser)GetSameLevelReward(userIdInt int64) (*MyUserInfo, error) {
	arangoDb := d.Db.ArangoDb
	userId 	 := strconv.FormatInt(userIdInt,10)
	myUserInfo, err := arango.Info(userIdInt)
	if err != nil {
		return &MyUserInfo{}, err
	}
	if myUserInfo.Id <= 0 {
		return &MyUserInfo{}, err
	}
	ctx := context.Background()
	// FILTER  v.status==1 and v.deleted==0
	query := "FOR v IN 1..99999999 OUTBOUND '"+arango.Table()+"/" + userId + "' "+arango.TableUserOf()+"  FILTER  v.status==1  RETURN v"
	cursor, err := arangoDb.Query(ctx, query, nil)
	if err != nil {
		// handle error
		return &MyUserInfo{}, err
	}
	defer cursor.Close()

	var lArr = make([]map[string]int64, 0)
	var left1 = make(map[string]int64) // 一代奖励

	for {
		var data *MyUserInfo
		var lMap = make(map[string]int64)
		_, err := cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
			return &MyUserInfo{}, err
		}
		if data != nil {
			lMap["Id"]	  = data.Id
			lMap["Level"] = data.Level
			lArr = append(lArr, lMap)
		}

	}
	if len(lArr) > 0 {
		for i := range lArr {
			if lArr[i]["Level"] == myUserInfo.Level {
					left1["Id"]	   = lArr[i]["Id"]
					left1["Level"] = lArr[i]["Level"]
					break
			} else if lArr[i]["Level"]-myUserInfo.Level > 0 {
				break
			}
		}
	} else {
		return &MyUserInfo{}, err
	}
	leftInfo,_ := arango.Info(left1["Id"])
	return &leftInfo, nil
}


// 越级用户信息(废除)
func (arango *ArangoUser)GetLeapFrogInfo(userIdInt int64) (*MyUserInfo, error) {
	arangoDb := d.Db.ArangoDb
	userId 	 := strconv.FormatInt(userIdInt,10)
	myUserInfo, err := arango.Info(userIdInt)
	if err != nil {
		return &MyUserInfo{}, err
	}
	if myUserInfo.Id <= 0 {
		return &MyUserInfo{}, err
	}
	ctx := context.Background()
	// FILTER  v.status==1 and v.deleted==0
	query := "FOR v IN 1..99999999 OUTBOUND '"+arango.Table()+"/" + userId + "' "+arango.TableUserOf()+"  FILTER  v.status==1  RETURN v"
	cursor, err := arangoDb.Query(ctx, query, nil)
	if err != nil {
		// handle error
		return &MyUserInfo{}, err
	}
	defer cursor.Close()

	var lArr = make([]map[string]int64, 0)
	var left1 = make(map[string]int64) // 一代奖励

	for {
		var data *MyUserInfo
		var lMap = make(map[string]int64)
		_, err := cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
			return &MyUserInfo{}, err
		}
		if data != nil {
			lMap["Id"]	  = data.Id
			lMap["Level"] = data.Level
			lArr = append(lArr, lMap)
		}

	}
	if len(lArr) > 0 {
		for i := range lArr {
			if lArr[i]["Level"] == 0 { //上级等级==普通用户,  : 继续找
				continue
			}else if lArr[i]["Level"]-myUserInfo.Level >= 0{ // 上级等级 >= 当前用户等级  : 退出
				break
			}else if lArr[i]["Level"] < myUserInfo.Level && lArr[i]["Level"]  > 0 { 	// 上级等级<当前用户等级 并且 上级等级>普通用户: 取上级作为越级用户
				left1["Id"]	   = lArr[i]["Id"]
				left1["Level"] = lArr[i]["Level"]
				break
			}else {
				break
			}
		}
	} else {
		return &MyUserInfo{}, err
	}
	leftInfo,_ := arango.Info(left1["Id"])
	return &leftInfo, nil
}



// 注册用户关系
func (arango *ArangoUser)Create(mapData map[string]interface{})(int,error)  {
	arangoDb 	 := d.Db.ArangoDb

	userId 		 := mapData["id"].(int64)
	uidString 	 := strconv.FormatInt(userId,10)
	userName     := mapData["user_name"].(string)
	headImgUrl   := mapData["head_img_url"].(string)
	name   		 := mapData["name"].(string)
	gender    	 := mapData["gender"].(int64)
	level    	 := mapData["level"].(int64)
	nickname     := mapData["nickname"].(string)
	inviteCode   := mapData["invite_code"].(string)
	mobile   	 := mapData["mobile"].(string)
	inviteUpId   := mapData["invite_up_id"].(int64)
	isCompany    := mapData["is_company"].(int64)
	status    	 := mapData["status"].(int64)
	inviteTime   := mapData["invite_time"].(string)
	operator     := mapData["operator"].(string)
	createTime   := mapData["create_time"].(string)
	updateTime   := mapData["update_time"].(string)

	//nowTime 	:=  time.Now().Format("2006-01-02 15:04:05")
	userDocument := MyUserInfo{
		AId:arango.Table()+"/"+uidString,
		AKey:uidString,
		Id:userId,
		UserName:userName,
		Nickname:nickname,
		HeadImgUrl:headImgUrl,
		Name:name,
		Gender:gender,
		IsCompany:isCompany,
		Level:level,
		InviteUpId:inviteUpId,
		Mobile:mobile,
		InviteCode:inviteCode,
		InviteTime:inviteTime,
		Status:status,
		Operator:operator,
		CreateTime:createTime,
		UpdateTime:updateTime,
	}

	ctx := context.Background()
	_,err := arango.CheckCollection(ctx,1,arango.Table())
	if err != nil {
		return 0,err
	}
	col,err := arangoDb.Collection(ctx,arango.Table())
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

// 更改用户关系
func (arango *ArangoUser)Update(userId int64,patch  map[string]interface{})(string,error)  {
	uidString := strconv.FormatInt(userId,10)
	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	//patch := map[string]interface{}{
	//	"status":1,
	//}
	col,err := arangoDb.Collection(ctx,arango.Table())
	if err != nil {
		return "",errors.New("update collection err:"+err.Error())
	}
	meta,err := col.UpdateDocument(ctx,uidString,patch)
	if err != nil {
		return "",errors.New("update collection err:"+err.Error())
	}
	return meta.Key,nil
}

// 检查是否有col
func (arango *ArangoUser)CheckCollection(ctx context.Context,collectionType int,col string)(bool,error)  {
	arangoDb := d.Db.ArangoDb
	found,err := arangoDb.CollectionExists(ctx,col)
	if !found {
		// handle error
		status,err := arango.CreateCollection(ctx,collectionType,col)
		return status,err
	}
	if err != nil {
		return false,err
	}else{
		return true,err
	}

}

// 创建col
func (arango *ArangoUser)CreateCollection(ctx context.Context,collectionType int,col string)(bool,error)  {
	arangoDb := d.Db.ArangoDb

	options := &driver.CreateCollectionOptions{}
	// type: 1 document,2 edges
	if collectionType == 1 {
		options.Type = driver.CollectionTypeDocument
	}else{
		options.Type = driver.CollectionTypeEdge
	}
	_,err   := arangoDb.CreateCollection(ctx,col,options)
	if err != nil {
		return false,err
	}else{
		return true,nil
	}
}

// 用户的团队数量(包含自己)
func (arango *ArangoUser)GetTeamNum(userId int64,level int64)(int64,error)  {

	arangoDb := d.Db.ArangoDb
	ctx := driver.WithQueryCount(context.Background())
	//query := "FOR v IN 0..99999999999 INBOUND @table_user @table_user_team_of RETURN v"
	// PRUNE:截取(但是还是会返回当前用户),FILTER:过滤
	// POSITION(prunv, v._key): 判断 _key是否在 prunv 中，如果是就是true否就是false
	query := " let plevel = @level " +
		"let userId = @uId " +
		"let table_user = 'mlive_user/' " +
		"let prunv  = ( FOR v,e,p IN 1..99999999999 INBOUND CONCAT(table_user,userId) mlive_user_of " +
		"PRUNE v.level > plevel " +
		"FILTER v.level == plevel " +
		"let firstv = LAST( " +
		"FOR v2,e2,p2 IN 1..99999999 OUTBOUND v._id mlive_user_of " +
		"PRUNE v2.level < plevel OR v2._key == userId " +
		"FILTER v2.level == plevel " +
		"RETURN v2 " +
		") " +
		"FILTER firstv == null || firstv._key != userId " +
		"RETURN v._key " +
		") " +
		" FOR v,e,p IN 1..999999999 INBOUND  CONCAT(table_user, userId) mlive_user_of" +
		" PRUNE (v.level > plevel OR POSITION(prunv, v._key)) AND v.status == 1 OPTIONS {bfs: true}" +
		" FILTER !POSITION(prunv, v._key)" +
		" FILTER v.level <= plevel" +
		" RETURN {u:v, pid:REGEX_REPLACE(e._to, table_user, ''), id:v._key} "

	bindVars := map[string]interface{}{
		"level":level,
		"uId":strconv.FormatInt(userId,10),
	}
	cursor,err := arangoDb.Query(ctx,query,bindVars)

	if err != nil{
		return 0,err
	}
	defer cursor.Close()
	return cursor.Count()+1,nil
}

type AgoTeamList struct {
	Id    				int64 	 	`json:"id"`
	InviteUpId   		int64 		`json:"inviteUpId"`
	Level				int64 		`json:"level"`
}
// 团队列表
func (arango *ArangoUser)GetTeamList(userId int64,plevel int64)([]AgoTeamList,error)  {

	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	//query := "FOR v IN 0..99999999999 INBOUND @table_user @table_user_team_of RETURN v"
	// PRUNE:截取(但是还是会返回当前用户),FILTER:过滤
	// POSITION(prunv, v._key): 判断 _key是否在 prunv 中，如果是就是true否就是false
	query := " let plevel = @level " +
		"let userId = @uId " +
		"let table_user = 'mlive_user/' " +
		"let prunv  = ( FOR v,e,p IN 1..99999999999 INBOUND CONCAT(table_user,userId) mlive_user_of " +
		"PRUNE v.level > plevel " +
		"FILTER v.level == plevel " +
		"let firstv = LAST( " +
		"FOR v2,e2,p2 IN 1..99999999 OUTBOUND v._id mlive_user_of " +
		"PRUNE v2.level < plevel OR v2._key == userId " +
		"FILTER v2.level == plevel " +
		"RETURN v2 " +
		") " +
		"FILTER firstv == null || firstv._key != userId " +
		"RETURN v._key " +
		") " +
		" FOR v,e,p IN 1..999999999 INBOUND  CONCAT(table_user, userId) mlive_user_of" +
		" PRUNE (v.level > plevel OR POSITION(prunv, v._key)) AND v.status == 1 OPTIONS {bfs: true}" +
		" FILTER !POSITION(prunv, v._key)" +
		" FILTER v.level <= plevel" +
		" RETURN {u:v, pid:REGEX_REPLACE(e._to, table_user, ''), id:v.id, inviteUpId:v.invite_up_id, level:v.level} "

	bindVars := map[string]interface{}{
		"level":plevel,
		"uId":strconv.FormatInt(userId,10),
	}

	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil{
		return nil,err
	}
	defer cursor.Close()
	var doc []AgoTeamList
	for {
		var data *AgoTeamList
		_, err = cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
			return nil, err
		}
		if data != nil {
			doc = append(doc,AgoTeamList{
					Id:data.Id,
					InviteUpId:data.InviteUpId,
					Level:data.Level,
			})
		}
	}
	return doc, nil
}

// 当天用户新增团队数
func (arango *ArangoUser)GetTeamAddNum(userId int64,level int64)(int64,error)  {

	arangoDb := d.Db.ArangoDb
	ctx := driver.WithQueryCount(context.Background())
	//query := "FOR v IN 1..99999999999 INBOUND @table_user @table_user_team_of " +
	//		"filter v.create_time like @map_create_time RETURN v"

	query := " let plevel = @level " +
		"let userId = @uId " +
		"let table_user = 'mlive_user/' " +
		"let prunv  = ( FOR v,e,p IN 1..99999999999 INBOUND CONCAT(table_user,userId) mlive_user_of " +
		"PRUNE v.level > plevel " +
		"FILTER v.level == plevel " +
		"let firstv = LAST( " +
		"FOR v2,e2,p2 IN 1..99999999 OUTBOUND v._id mlive_user_of " +
		"PRUNE v2.level < plevel OR v2._key == userId " +
		"FILTER v2.level == plevel " +
		"RETURN v2 " +
		") " +
		"FILTER firstv == null || firstv._key != userId " +
		"RETURN v._key " +
		") " +
		" FOR v,e,p IN 1..999999999 INBOUND  CONCAT(table_user, userId) mlive_user_of" +
		" PRUNE (v.level > plevel OR POSITION(prunv, v._key)) AND v.status == 1 OPTIONS {bfs: true}" +
		" FILTER !POSITION(prunv, v._key)" +
		" FILTER v.level <= plevel " +
		" FILTER v.create_time like @map_create_time " +
		" RETURN {u:v, pid:REGEX_REPLACE(e._to, table_user, ''), id:v._key} "
		SameDay := time.Now().Format("2006-01-02")
		bindVars := map[string]interface{}{
			"level":level,
			"uId":strconv.FormatInt(userId,10),
			"map_create_time":"%"+SameDay+"%",
		}

	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil{
		return 0,err
	}
	defer cursor.Close()
	return cursor.Count(),nil
}

/**
* 检查是否在邀请下级中 true:是,false:否
 */
func (arango *ArangoUser)CheckIsInviteSubordinate(userId int64,inviteOutId int64)(bool,error)  {
	arangoDb := d.Db.ArangoDb
	ctx := driver.WithQueryCount(context.Background())
	query := "FOR v IN 1..99999999999 INBOUND @table_user @table_user_of filter v._key == @keyId RETURN v"
	bindVars := map[string]interface{}{
		"table_user":arango.Table()+"/"+strconv.FormatInt(userId,10),
		"table_user_of":arango.TableUserOf(),
		"keyId":strconv.FormatInt(inviteOutId,10),
	}
	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil{
		return false,err
	}
	defer cursor.Close()
	if cursor.Count() > 0 {
		return true,nil
	}else{
		return false,nil
	}
}

type UserId1TeamList struct {
	UserId				int64 		`json:"userId"`
	InviteUpId			int64 		`json:"inviteUpId"`
	Level				int64 		`json:"level"`
	Nickname			string 		`json:"nickname"`
	InviteCode			string 		`json:"inviteCode"`
	Status				int64 		`json:"status"`
	CreateTime			string 		`json:"createTime"`
	UpdateTime			string 		`json:"updateTime"`
}

// 用户id==1,直属的邀请下级和联创(分公司)
func (arango *ArangoUser)GetUserId1TeamList()([]UserId1TeamList,error)  {
	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	query := " FOR d in @@table_user FILTER d.invite_up_id == @map_invite_up_id or d.level == @map_level " +
		"return {userId:d.id,inviteUpId:d.invite_up_id," +
		"level:d.level,nickname:d.nickname,inviteCode:d.inviteCode," +
		"status:d.Status,createTime:d.createTime,updateTime:d.updateTime}"
	bindVars := map[string]interface{}{
		"@table_user":arango.Table(),
		"map_invite_up_id":1,
		"map_level":5,
	}
	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil {
		return nil,err
	}
	defer cursor.Close()
	var doc []UserId1TeamList
	for  {
		var data *UserId1TeamList
		_,err := cursor.ReadDocument(ctx,&data)
		if driver.IsNoMoreDocuments(err){
			break
		}else if err != nil {
			//log.Fatal("info read document err:",err)
			return nil,err
		}
		if data != nil {
			doc = append(doc,UserId1TeamList{
				UserId:data.UserId,
				InviteUpId:data.InviteUpId,
				Level:data.Level,
				Nickname:data.Nickname,
				InviteCode:data.InviteCode,
				Status:data.Status,
				CreateTime:data.CreateTime,
				UpdateTime:data.UpdateTime,
			})
		}
	}
	return doc,nil
}



// 获取邀请列表列表
func (arango *ArangoUser) GetInviteList(userId int64,inviteUpId int64)([]MyUserList,error)  {

	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	query := "for d in @@table_user   "
	bindVars := map[string]interface{}{
		"@table_user":arango.Table(),
	}
	if  userId  >0 || inviteUpId > 0   {
		query = query + " FILTER "
		if userId > 0 {
			bindVars["userId"] = strconv.FormatInt(userId, 10)
			query = query + "   d._key == @userId "
		}else
		if inviteUpId >= 0 {
			query = query + "   d.invite_up_id == @inviteUpId "
			bindVars["inviteUpId"] = inviteUpId
		}
	}else{
		bindVars["inviteUpId"] = 1
		query = query + "  FILTER d.invite_up_id == @inviteUpId "
	}
	query = query+"  SORT d.id ASC  return d "
	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil {
		return nil,err
	}
	defer cursor.Close()
	var doc []MyUserList
	var initArr  []*MyUserList
	for {
		var data *MyUserInfo
		_,err := cursor.ReadDocument(ctx,&data)
		if driver.IsNoMoreDocuments(err){
			break
		}else if err != nil {
			//log.Fatal("info read document err:",err)
			return nil,err
		}
		if data != nil {

			doc = append(doc,MyUserList{
				Id 			 :data.Id,
				InviteUpId	 :data.InviteUpId,
				InviteTime   :data.InviteTime,
				UserName	 :data.UserName,
				HeadImgUrl	 :data.HeadImgUrl,
				Nickname	 :data.Nickname,
				Name    	 :data.Name,
				Gender		 :data.Gender,
				Level		 :data.Level,
				InviteCode	 :data.InviteCode,
				Mobile	 	 :data.Mobile,
				TelPhone	 :data.TelPhone,
				Email		 :data.Email,
				IsCompany  	 :data.IsCompany,
				Status		 :data.Status,
				Operator	 :data.Operator,
				AdminUserId  :data.AdminUserId,
				CreateTime   :data.CreateTime,
				UpdateTime   :data.UpdateTime,
				ChildList    :initArr,
			})
		}
	}

	return doc,nil
}


// 向父上层寻找所有列表
func (arango *ArangoUser)GetInviteOutBoundList(userId int64)([]MyUserList,error)  {

	arangoDb := d.Db.ArangoDb
	ctx := context.Background()

	query := "FOR v IN 1..99999999999 OUTBOUND @table_user @table_user_of filter v.level <= 5 RETURN v"
	bindVars := map[string]interface{}{
		"table_user":arango.Table()+"/"+strconv.FormatInt(userId,10),
		"table_user_of":arango.TableUserOf(),
	}
	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil {
		return nil,err
	}
	defer cursor.Close()
	var doc []MyUserList
	var initArr []*MyUserList
	for {
		var data *MyUserInfo
		_,err := cursor.ReadDocument(ctx,&data)
		if driver.IsNoMoreDocuments(err){
			break
		}else if err != nil {
			//log.Fatal("info read document err:",err)
			return nil,err
		}
		if data != nil {
			doc = append(doc,MyUserList{
				Id 			 :data.Id,
				InviteUpId	 :data.InviteUpId,
				InviteTime   :data.InviteTime,
				UserName	 :data.UserName,
				HeadImgUrl	 :data.HeadImgUrl,
				Nickname	 :data.Nickname,
				Name    	 :data.Name,
				Gender		 :data.Gender,
				Level		 :data.Level,
				InviteCode	 :data.InviteCode,
				Mobile	 	 :data.Mobile,
				TelPhone	 :data.TelPhone,
				Email		 :data.Email,
				IsCompany  	 :data.IsCompany,
				Status		 :data.Status,
				Operator	 :data.Operator,
				AdminUserId  :data.AdminUserId,
				CreateTime   :data.CreateTime,
				UpdateTime   :data.UpdateTime,
				ChildList   :initArr,
			})
		}
	}
	return doc,nil
}


type UserListStruct struct {
	Id  				int64		`json:"id"`
	Nickname  			string		`json:"nickname"`
	InviteCode  		string		`json:"invite_code"`
	Mobile				string		`json:"mobile"`
	Level				int64		`json:"level"`
	InviteUpId			int64		`json:"invite_up_id"`
	CreateTime			string		`json:"create_time"`
	UpdateTime			string		`json:"update_time"`
	Status				int64		`json:"status"`
	Operator			string		`json:"operator"`
}

// 管理后台查询用户
func (arango *ArangoUser)GetAdminUserList(dataBind map[string]interface{},row int,page int) ([]*UserListStruct,int64,error) {

	arangoDb := d.Db.ArangoDb
	ctx := context.Background()
	query := "for u in @@table_user   "
	bindVars := map[string]interface{}{
		"@table_user":arango.Table(),
	}

	if len(dataBind) > 0 {
		query = query+" FILTER "
		var userId  int64 = 0
		if dataBind["userId"] != nil {
			userId = dataBind["userId"].(int64)
			if userId > 0 {
				bindVars["userId"] = strconv.FormatInt(userId, 10)
				query = query + "   u._key == @userId "
			}
		}

		var  mobile string = ""
		if dataBind["mobile"] != nil {
			mobile = dataBind["mobile"].(string)

			if mobile != "" {
				if userId > 0 {
					query = query+" and "
				}
				bindVars["mobile"] = mobile
				query = query + "   u.mobile == @mobile "
			}
		}
		var nickname string = ""
		if dataBind["nickname"] != nil {
			nickname = dataBind["nickname"].(string)
			if nickname != "" {
				if mobile != "" {
					query = query + " and "
				}
				bindVars["nickname"] = "%"+nickname+"%"
				query = query + "   u.nickname like @nickname "
			}
		}

		var level int64 = -1
		if dataBind["level"] != nil {
			level = dataBind["level"].(int64)
			if level >= 0 {
				if nickname !=  "" {
					query = query + " and "
				}
				bindVars["level"] = level
				query = query + "   u.level == @level "
			}
		}
		var inviteUpId int64 = 0
		if dataBind["inviteUpId"] != nil {
			inviteUpId = dataBind["inviteUpId"].(int64)
			if inviteUpId > 0 {
				if level >= 0 {
					query = query + " and "
				}
				bindVars["inviteUpId"] = inviteUpId
				query = query + "   u.invite_up_id == @inviteUpId "
			}
		}
	}
	offset := (page-1)*row
	offsetString := strconv.Itoa(offset)
	rowString := strconv.Itoa(row)
	query = query+"  SORT u.update_time DESC  "
	queryCount := query
	query = query+"  limit "+offsetString+","+rowString+" return u "

	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil {
		return nil,0,err
	}
	defer cursor.Close()
	var doc []*UserListStruct
	for {
		var data *UserListStruct
		_,err := cursor.ReadDocument(ctx,&data)
		if driver.IsNoMoreDocuments(err){
			break
		}else if err != nil {
			//log.Fatal("info read document err:",err)
			return nil,0,err
		}
		if data != nil {
			log.Println("data result:",data)
			doc = append(doc,&UserListStruct{
				Id:data.Id,
				Nickname:data.Nickname,
				InviteCode:data.InviteCode,
				Mobile:data.Mobile,
				Level:data.Level,
				InviteUpId:data.InviteUpId,
				CreateTime:data.CreateTime,
				UpdateTime:data.UpdateTime,
				Status:data.Status,
				Operator:data.Operator,
			})
		}
	}

	// 统计数量
	ctx1 := driver.WithQueryCount(context.Background())
	queryCount = queryCount+" return u  "
	log.Println("queryCount:",queryCount)
	cursorCount,err := arangoDb.Query(ctx1,queryCount,bindVars)
	if err != nil {
		return nil,0,err
	}
	defer cursorCount.Close()
	return doc,cursorCount.Count(),nil
}


type MoneySort struct {
	UpdateTime string
	Id  int64
}

// Persons a set of person
type MoneySorts []MoneySort

// Len return count
func (p MoneySorts) Len() int {
	return len(p)
}

// Less return bigger true  1、优先排序时间desc,再id desc
func (p MoneySorts) Less(i, j int) bool {
	if p[i].UpdateTime == p[j].UpdateTime {
		return p[i].Id > p[j].Id
	}else{
		return p[i].UpdateTime > p[j].UpdateTime
	}

}

// Swap swap items
func (p MoneySorts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// 获取-分润上级-所有的下级
func (arango *ArangoUser)GetAdminMoneyList(userId int64,dataBind map[string]interface{},row int,page int) ([]*UserListStruct,int64,error) {

	arangoDb := d.Db.ArangoDb
	ctx 	 := context.Background()

	// 1、获取所有的下级
	query := "FOR v IN 1..99999999999 INBOUND @table_user @table_user_of RETURN v"
	bindVars := map[string]interface{}{
		"table_user":arango.Table()+"/"+strconv.FormatInt(userId,10),
		"table_user_of":arango.TableUserOf(),
	}
	cursor,err := arangoDb.Query(ctx,query,bindVars)
	if err != nil {
		return nil,0,err
	}
	defer cursor.Close()
	var child []MyUserList
	for {
		var data *MyUserInfo
		_,err := cursor.ReadDocument(ctx,&data)
		if driver.IsNoMoreDocuments(err){
			break
		}else if err != nil {
			//log.Fatal("info read document err:",err)
			return nil,0,err
		}
		if data != nil {
			child = append(child,MyUserList{
				Id 			 :data.Id,
				InviteUpId	 :data.InviteUpId,
				InviteTime   :data.InviteTime,
				UserName	 :data.UserName,
				HeadImgUrl	 :data.HeadImgUrl,
				Nickname	 :data.Nickname,
				Name    	 :data.Name,
				Gender		 :data.Gender,
				Level		 :data.Level,
				InviteCode	 :data.InviteCode,
				Mobile	 	 :data.Mobile,
				TelPhone	 :data.TelPhone,
				Email		 :data.Email,
				IsCompany  	 :data.IsCompany,
				Status		 :data.Status,
				Operator	 :data.Operator,
				AdminUserId  :data.AdminUserId,
				CreateTime   :data.CreateTime,
				UpdateTime   :data.UpdateTime,
			})
		}
	}
	if len(child) > 0 {
		var moneySorts = MoneySorts{}
		for _,v := range child {

			info,err := arango.GetMoneySuperior(v.Id)
			if err != nil {
				return []*UserListStruct{},0,err
			}
			if info.Id == userId {
				moneySorts = append(moneySorts,MoneySort{Id:v.Id,UpdateTime:v.UpdateTime})
			}
		}
		if len(moneySorts) > 0 {
			sort.Sort(moneySorts)
			var dataArr = []*UserListStruct{}
			for _,v := range moneySorts  {

				ainfo,_ :=arango.Info(v.Id)
				if ainfo.Id > 0 {
					var dataMap = &UserListStruct{
						Id:ainfo.Id,
						Nickname:ainfo.Nickname,
						InviteCode:ainfo.InviteCode,
						Mobile:ainfo.Mobile,
						Level:ainfo.Level,
						InviteUpId:ainfo.InviteUpId,
						CreateTime:ainfo.CreateTime,
						UpdateTime:ainfo.UpdateTime,
						Status:ainfo.Status,
						Operator:ainfo.Operator,
					}
					dataArr = append(dataArr,dataMap)
				}
			}
			if len(dataArr) <= 0 {
				return []*UserListStruct{},0,nil
			}
			var count  = int64(moneySorts.Len())
			startIndex := (page-1)*row
			endIndex   := startIndex+row
			var totalPage   = math.Ceil(float64(count)/float64(row))
			if totalPage < float64(page) {
				return []*UserListStruct{},count,nil
			}
			if int(count)-(endIndex) < 0 {
				endIndex = int(count)
			}
			return dataArr[startIndex:endIndex],count,nil
		}
		return []*UserListStruct{},0,nil
	}else{
		return []*UserListStruct{},0,nil
	}

}

func (arango *ArangoUser)CreateIndex()  {
	arangoDb := d.Db.ArangoDb
	ctx 	 := context.Background()

	log.Println("aaa:开始")
	index,_ := arangoDb.Collection(ctx,arango.Table())
	//fields := []string{"level","status"}
	//
	// op := &driver.EnsureHashIndexOptions{}
	//in1,t1,e1 := index.EnsureHashIndex(ctx,fields,op)
	//log.Println("in1:",in1,"t1:",t1,"e1:",e1)
	//
	//fields1 := []string{"level","status"}
	//op1 := &driver.EnsureSkipListIndexOptions{}
	//in2,t2,e2 := index.EnsureSkipListIndex(ctx,fields1,op1)
	//log.Println("in2:",in2,"t2:",t2,"e2:",e2)

	fields2 := []string{"level","invite_up_id","nickname","create_time","status"}
	op2  := &driver.EnsureHashIndexOptions{}
	in3,t3,e3 := index.EnsureHashIndex(ctx,fields2,op2)
	log.Println("in3:",in3,"t3:",t3,"e3:",e3)



}

