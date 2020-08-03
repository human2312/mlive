package dao

import (
	"log"
	"mlive/dao"
	"mlive/library/config"
	"mlive/service/java"
	"mlive/util"
	"testing"
	"mlive/library/logger"
	"time"
)
var(
	u		 User
)

func init()  {

	// 初始化配置
	config.Init()
	// 初始化日志
	logger.Init()
	// 链接数据库
	dao.NewDao()
	// 初始化mq
	util.NewMQ()
}

// 注册用户
func TestRegister(t *testing.T)  {
	//
	//var (
	//	userId	  int64  = 19
	//	isCompany int64  = 1
	//	level	  int64  = 1
	//	inviteId  int64  = 9
	//)
	//userId,err :=  u.Register(userId,isCompany,level,inviteId)
	//t.Log("reg return userId:",userId,err)
}

// 初始化ago
func TestInitAgoData(t *testing.T)  {
	res,err := u.InitAgoData()
	t.Log("res:",res)
	t.Log("err:",err)
}


// java sync time https://blog.csdn.net/mirage003/article/details/86073046
func TestJavaUserTime(t *testing.T)  {

	var (
		userId int64 = 11527
	)
	userInfo,_ := u.GetUserInfo(userId)

	log.Println("userInfo:",userInfo)

	javaInfo,_ :=java.InfoById(userId)
	log.Println("javaInfo:",javaInfo)
	javaUpdateTime,_ := time.ParseInLocation("2006-01-02 15:04:05",javaInfo.UpdateTime,time.Local)
	log.Println("java time :",javaInfo.UpdateTime)
	log.Println("java time ss :",javaUpdateTime)
	log.Println("userInfo time :",userInfo.UpdateTime)
	if javaUpdateTime.Unix() > userInfo.UpdateTime.Unix() {
		log.Println("111:",javaUpdateTime)
	}

	timeTemplate1 := "2006-01-02 15:04:05" //常规类型
	t1 := "2019-01-08 13:50:30" //外部传入的时间字符串
	stamp, _ := time.ParseInLocation(timeTemplate1, t1, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	log.Println(stamp.Unix())  //输出：1546926630


}
