package graceful

/**
 * @Author: Lemyhello
 * @Description: 平滑重启
 * @File:  graceful
 * @Version: X.X.X
 * @Date: 2020/3/25 下午3:57
 */

var (
	grace = make(chan bool,1)
)

func Pointer()  *chan bool{
	return &grace
}

func Get() (bool) {
	return <- grace
}

func Put(isValue bool) {
	grace <- isValue
}

func GetChan() (chan bool) {
	return grace
}