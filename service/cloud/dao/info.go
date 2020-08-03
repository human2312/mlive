package dao

import (
	d "mlive/dao"
	"mlive/service/cloud/model"
	"errors"
)

type  CloudStorage struct {

}

type CloudStorageSurplusNum struct {
	Number 		int64	`json:"number"`
}
// 总云仓-剩余数量
func (cloud *CloudStorage)GetCloudStorageSurplusNum(userId int64)(int64,error)  {
	if userId <= 0 {
		return 0,errors.New("Please enter user ID")
	}
	var Results CloudStorageSurplusNum
	d.Db.DB.Table(model.MliveCloudStorage{}.TableName()).Where("user_id = ?",userId).Where("status >= ?",0).
		Select("sum(number) as number").Scan(&Results)
	return Results.Number,nil
}

// 云仓使用-总数量-(过滤无效订单)
func (cloud *CloudStorage)GetCloudStoragUseCount(useUserId int64)(int64,error)  {
	if useUserId <= 0 {
		return 0,errors.New("Please enter user ID")
	}
	var total int64 = 0
	if err := d.Db.DB.Model(&model.MliveCloudStorage{}).Where("use_user_id = ?",useUserId).Where("status >= ?",0).Count(&total).Error; err != nil {
		return 0,err
	}
	return total,nil
}

// 云仓使用-待处理数量
func (cloud *CloudStorage)GetCloudStoragStay(useUserId int64)(int64,error)  {
	if useUserId <= 0 {
		return 0,errors.New("Please enter user ID")
	}
	var total int64 = 0
	if err := d.Db.DB.Model(&model.MliveCloudStorage{}).Where("use_user_id = ?",useUserId).Where("status = ?",0).Count(&total).Error; err != nil {
		return 0,err
	}
	return total,nil
}

// 首页云仓返回-数据列表
func (cloud *CloudStorage)GetCloudStoragList(userId int64,page int64,row int64)([]*model.MliveCloudStorage,error)  {
	var data = []*model.MliveCloudStorage{}
	u := d.Db.DB.Model(&model.MliveCloudStorage{}).Order("id desc").Limit(row).Offset((page-1)*row)
	u = u.Where("user_id = ? ",userId).Where("type=?",2).Where("channel <= ?",2).Where("status=?",1)
	if err := u.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}
// 首页云仓返回-总数量
func (cloud *CloudStorage)GetCloudStoragCount(userId int64)(int64,error)  {
	var total int64 = 0
	u := d.Db.DB.Model(&model.MliveCloudStorage{}).Order("id desc")
	u = u.Where("user_id = ? ",userId).Where("type=?",2).Where("channel <= ?",2).Where("status=?",1)
	if err := u.Count(&total).Error; err != nil{
		return 0,err
	}
	return  total,nil
}


//获取云仓-订单信息
func (cloud *CloudStorage)GetCloudStoragOrderInfo(userId int64,No int64)(*model.MliveCloudStorage,error)  {

	var data = model.MliveCloudStorage{}
	db     :=  d.Db.DB
	if userId > 0 {
		db = db.Where("user_id = ?",userId)
	}
	if No > 0 {
		db = db.Where("no = ?",No)
	}
	if  No <= 0 {
		return nil,errors.New("【no】 Parameter error ")
	}
	if err := db.First(&data).Error; err != nil{
		return &model.MliveCloudStorage{},err
	}
	return  &data,nil
}

// 订单号-查云仓
func (cloud *CloudStorage)GetCloudStoragOrderNoInfo(orderNo string,types int64)(*model.MliveCloudStorage,error)  {

	var data = model.MliveCloudStorage{}
	db     :=  d.Db.DB
	db = db.Where("order_no = ?",orderNo)
	if types > 0 {
		db = db.Where("type = ?",types)
	}
	if err := db.First(&data).Error; err != nil{
		return &model.MliveCloudStorage{},err
	}
	return  &data,nil
}

// 订单号-修改
func (cloud *CloudStorage)SaveCloudLogByOrderNo(orderNo string,mapData map[string]interface{})(bool,error) {
	//mapData["quantity"] = gorm.Expr("quantity-?",1)
	db := d.Db.DB.Model(&model.MliveCloudStorage{})
	db = db.Where("order_no = ?",orderNo)
	if err := db.UpdateColumn(mapData).Error; err != nil {
		return false, err
	} else {
		return true, nil
	}
}



// 获取运营后台最后操作的云仓数据
func (cloud *CloudStorage)GetAdminCloudLastInfo(userId int64)(*model.MliveCloudStorage,error)  {
	var data = model.MliveCloudStorage{}
	db     :=  d.Db.DB.Order("id desc").Where("user_id = ?",userId).Where("channel=?",3)
	if err := db.First(&data).Error; err != nil{
		return &model.MliveCloudStorage{},err
	}
	return  &data,nil
}

// 首页云仓返回-数据列表
func (cloud *CloudStorage)GetCloudStoragTimingList()([]*model.MliveCloudStorage,error)  {
	var data = []*model.MliveCloudStorage{}
	u := d.Db.DB.Model(&model.MliveCloudStorage{}).Order("id desc")
	u = u.Where("status=?",0)
	if err := u.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}

//云仓交易记录-展示 充值vip、店长、总监、合伙人、联仓的记录 (type=1&order_no!=""&cloud_type>=1)
func (cloud *CloudStorage)GetAdminPayCloudStoragList(page int,row int,mapBind  map[string]interface{})([]*model.MliveCloudStorage,error)  {
	var data = []*model.MliveCloudStorage{}
	u := d.Db.DB.Model(&model.MliveCloudStorage{}).Order("id desc").Limit(row).Offset((page-1)*row)
	if len(mapBind) > 0 {
		if mapBind["userId"] != nil {
			u = u.Where("user_id =?",mapBind["userId"])
		}
	}
	u = u.Where("type =?",1).Where("cloud_type >= ?",1).Where("order_no !=?"," ")
	if err := u.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}

func (cloud *CloudStorage)GetAdminPayCloudStoragCount(mapBind  map[string]interface{})(int64,error)  {
	var count int64 = 0
	u := d.Db.DB.Model(&model.MliveCloudStorage{}).Order("id desc")
	if len(mapBind) > 0 {
		if mapBind["userId"] != nil {
			u = u.Where("user_id =?",mapBind["userId"])
		}
	}
	u = u.Where("type =?",1).Where("cloud_type >= ?",1).Where("order_no !=?"," ")
	if err := u.Count(&count).Error; err != nil{
		return count,err
	}
	return  count,nil
}



// 获取云仓所有的数据列表(待处理+已经处理)
func (cloud *CloudStorage)GetAdminAllCloudStoragList(page int,row int,mapBind  map[string]interface{})([]*model.MliveCloudStorage,error)  {
	var data = []*model.MliveCloudStorage{}
	u := d.Db.DB.Model(&model.MliveCloudStorage{}).Order("id desc").Limit(row).Offset((page-1)*row)
	if len(mapBind) > 0 {
		if mapBind["userId"] != nil {
			u = u.Where("user_id =?",mapBind["userId"])
		}
		if mapBind["type"] != nil {
			u = u.Where("type =?",mapBind["type"])
		}
		if mapBind["useUserId"] != nil {
			u = u.Where("use_user_id =?",mapBind["useUserId"])
		}
		if mapBind["orderNo"] != nil {
			u = u.Where("order_no =?",mapBind["orderNo"])
		}
	}
	u = u.Where("status >=?",0)
	if err := u.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}
// 获取云仓所有的数据总数(待处理+已经处理)
func (cloud *CloudStorage)GetAdminAllCloudStoragCount(mapBind map[string]interface{})(int64,error)  {
	var count int64 = 0
	u := d.Db.DB.Model(&model.MliveCloudStorage{})
	if len(mapBind) > 0 {
		if mapBind["userId"] != nil {
			u = u.Where("user_id =?",mapBind["userId"])
		}
		if mapBind["type"] != nil {
			u = u.Where("type =?",mapBind["type"])
		}
		if mapBind["useUserId"] != nil {
			u = u.Where("use_user_id =?",mapBind["useUserId"])
		}
		if mapBind["orderNo"] != nil {
			u = u.Where("order_no =?",mapBind["orderNo"])
		}
	}
	u = u.Where("status >=?",0)
	if err := u.Count(&count).Error; err != nil{
		return count,err
	}
	return  count,nil
}