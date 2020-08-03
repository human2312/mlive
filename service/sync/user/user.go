package user

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"log"
	d "mlive/dao"
	"mlive/library/logger"
	"mlive/service/sync/model"
	dx "mlive/service/sync/mysql"
	"strconv"
)

var (
	DbX *dx.Dao
)


func initNew()  {
	dSN := viper.GetString("mysqlXfyl.dsn")
	active, _ := strconv.Atoi(viper.GetString("mysqlXfyl.active"))
	idle, _ := strconv.Atoi(viper.GetString("mysqlXfyl.idle"))
	idleTimeout, _ := strconv.Atoi(viper.GetString("mysqlXfyl.idleTimeout"))
	conf := dx.DbConf{DSN:dSN,Active:active,Idle:idle,IdleTimeout:idleTimeout}
	var err error
	DbX,err = dx.New(conf)
	if err != nil {
		logger.Eprintf("db dsn(%s) error:(%v) ", dSN, err)
	}
}



var (
	userCh 	 	 chan []model.ShopUser
	userWait	 chan int
	userStatus	 chan bool
)

// 用户基础表
func User()  {
	// 初始化
	initNew()
	//uch := &userChan{closeChan:make(chan byte,1)}
	userCh  	= make(chan []model.ShopUser,1000)
	userWait 	= make(chan int,1)

	var initId int64 = 0
	go func() {
		for  {
			data ,xid := XfylUser(initId)
			if xid > 0 {
				initId = xid
				userCh <- data
			}else{
				log.Println("data is empty end ")
				close(userCh)
				break
			}
		}
	}()


	var i = 0

  	 go func() {
  	 	var (
  	 		e  []model.ShopUser
  	 		ok bool = true
		)
  	 	for{
			select {
  	 			case e,ok = <- userCh:
				if !ok {
					log.Println("End ok.")
					break
				}
  	 			i = i+len(e)
				//log.Println("ch data:-->:",e)
				//log.Println("ch num:-->:",i)
				ZbUser(e)
			}
			if !ok { // 关闭for
				userWait <- 0
				log.Println("End userWait.")
				break
			}
		}

	 }()
	<- userWait
}

func XfylUser(id int64)([]model.ShopUser,int64)  {

	var initId int64 = 0
	var data  []model.ShopUser
	user  := DbX.DB.Model(&model.ShopUser{})
	if id  > 0 {
		user = user.Where("id > ?",id)
	}
	if err := user.Order("id asc").Limit(300).Find(&data).Error; err != nil{
		log.Println("获取失败 shop user 失败:",err)
		return nil,initId
	}
	num := len(data)
	if num > 1 {
		for i,k := range data{
			var inviteData model.ShopUserInvite
			if err :=  DbX.DB.Model(&model.ShopUserInvite{}).Where("id = ? ",k.ID).First(&inviteData).Error; err != nil{
				log.Println("获取用户关系失败:", k)
			}else {
				data[i].Invite = inviteData
			}

		}

		id = data[num-1].ID
		return data,id
	}
	return nil,initId
}



func ZbUser(data []model.ShopUser)  {

	// 批量写入
	var (
		batchDataList	 []*model.MliveUser
		batchInviteList  []*model.MliveUserInvite
	)
	// 组装
	for _,k := range data{
		batchDataList = append(batchDataList,
			&model.MliveUser{
				ID:k.ID,
				XfylId:k.ID,
				UserName:k.UserName,
				Password:k.Password,
				HeadImgUrl:k.HeadImgUrl,
				Nickname:k.Name,
				Name:k.Name,
				Level:1,
				Gender:k.Gender,
				Mobile:k.Mobile,
				TelPhone:k.TelPhone,
				Status:k.Status,
				Deleted:k.Deleted,
				Shield:k.Shield,
				CreateTime:k.AddTime,
				UpdateTime:k.UpdateTime,
			})
		if k.Invite.ID > 1  && k.Invite.Pid > 0 {
			batchInviteList = append(batchInviteList,
				&model.MliveUserInvite{
					ID:         k.Invite.ID,
					Pid:        k.Invite.Pid,
					CreateTime: k.Invite.AddTime,
					UpdateTime: k.Invite.UpdateTime,
				})
		}
	}
	userErr := BatchUserSave(batchDataList)
	log.Println("userErr:",userErr)
	if len(batchInviteList) > 0 {
		userInviteErr := BatchUserInviteSave(batchInviteList)
		log.Println("userInviteErr:",userInviteErr)
	}

	// 循环一条一条判断写入
	//for _,k := range data{
	//	userId  := k.ID
	//	var info = model.MliveUser{}
	//	if err := d.Db.DB.First(&info,userId).Error; err != nil{
	//		//log.Println("error-->read",err)
	//		//return
	//	}
	//	if info.ID > 0 {
	//		// 更新
	//		log.Println("更新:",info.ID)
	//		if err := d.Db.DB.Model(&model.MliveUser{}).Updates(model.MliveUser{ID:userId,UserName:k.UserName}).Error; err != nil {
	//			log.Println("error--> update write",err)
	//			return
	//		}
	//	}else{
	//		log.Println("写入:",info.ID)
	//	 	// 写入
	//	 	if err := d.Db.DB.Create(model.MliveUser{ID:userId,UserName:k.UserName}).Error; err != nil {
	//			log.Println("error-->create write",err)
	//			return
	//		}
	//	}
	//}
}

func BatchUserSave(emp []*model.MliveUser) error  {
	var buffer bytes.Buffer
	sql := " insert into `mlive_user` " +
		"(`id`,`xfyl_id`,`user_name`,`password`," +
		"`head_img_url`,`nickname`,`name`,`level`,`gender`," +
		"`mobile`,`tel_phone`,`status`,`deleted`,`create_time`,`update_time`) values "
	if _,err := buffer.WriteString(sql); err != nil {
		return  err
	}
	for i,e := range emp {
		if i == len(emp)-1 {
			buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s','%s','%s','%s','%d','%d','%s','%s','%d','%d','%s','%s');",
				e.ID,e.XfylId,e.UserName,e.Password,e.HeadImgUrl,e.Nickname,e.Name,e.Level,e.Gender,e.Mobile,e.TelPhone,e.Status,e.Deleted,e.CreateTime.Format("2006-01-02 15:04:05"),e.UpdateTime.Format("2006-01-02 15:04:05")))
		}else{
			buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s','%s','%s','%s','%d','%d','%s','%s','%d','%d','%s','%s'),",
				e.ID,e.XfylId,e.UserName,e.Password,e.HeadImgUrl,e.Nickname,e.Name,e.Level,e.Gender,e.Mobile,e.TelPhone,e.Status,e.Deleted,e.CreateTime.Format("2006-01-02 15:04:05"),e.UpdateTime.Format("2006-01-02 15:04:05")))

		}
	}
	return d.Db.DB.Exec(buffer.String()).Error
}
func BatchUserInviteSave(emp []*model.MliveUserInvite) error  {
	var bufferInvite bytes.Buffer
	sql := " insert into `mlive_user_invite` (`id`,`pid`,`create_time`,`update_time`) values"
	if _,err := bufferInvite.WriteString(sql); err != nil {
		return  err
	}
	for i,e := range emp {
		if i == len(emp)-1 {
			bufferInvite.WriteString(fmt.Sprintf("('%d','%d','%s','%s');",
				e.ID,e.Pid,e.CreateTime.Format("2006-01-02 15:04:05"),e.UpdateTime.Format("2006-01-02 15:04:05")))
		}else{
			bufferInvite.WriteString(fmt.Sprintf("('%d','%d','%s','%s'),",
				e.ID,e.Pid,e.CreateTime.Format("2006-01-02 15:04:05"),e.UpdateTime.Format("2006-01-02 15:04:05")))
		}
	}
	return d.Db.DB.Exec(bufferInvite.String()).Error
}