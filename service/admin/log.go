package admin

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  log
 * @Version: X.X.X
 * @Date: 2020/4/3 上午10:52
 */

var (
	rabbitmqLog = RabbitmqLog{}
)

//MqAdd RabbitMQ日志
func MqAdd(messageId string, queue string , types string , message string,result int) bool {
	rabbitmqLog.Uid = messageId
	rabbitmqLog.Msg = message
	rabbitmqLog.Queue = queue
	rabbitmqLog.Types = types
	rabbitmqLog.Result = result
	return rabbitmqLog.Add(rabbitmqLog)
}

//MqUpdate RabbitMQ处理结果更新
func MqUpdate(messageId string,result int) bool {
	rabbitmqLog.Uid = messageId
	rabbitmqLog.Result = result
	return rabbitmqLog.Updates(rabbitmqLog)
}
