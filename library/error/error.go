package error

/**
 * @Author: Lemyhello
 * @Description: 重定义error接口
 * @File:  error
 * @Version: X.X.X
 * @Date: 2020/3/25 上午10:51
 */

import (
	"encoding/json"
)

type Err struct {
	Code int
	Msg   string
}

func (e *Err) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func New(code int, msg string) *Err {
	return &Err{
		Code: code,
		Msg:   msg,
	}
}