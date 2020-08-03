package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/deckarep/golang-set"
	"regexp"
	"strconv"
	"strings"
)

type Response struct {
	 code int
	 msg  string
}

// get raw data params
func GetRawData(r *http.Request,checkParams []interface{})(map[string]interface{},error)  {
	buf 		 := make([]byte, 1024)
	param, _ 	 := r.Body.Read(buf)
	paramsJson   := string(buf[0:param])
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(paramsJson),&dat);err != nil {
		return dat,err
	}else{
		if len(checkParams) > 0 {
			checkSlice := mapset.NewSetFromSlice(checkParams)
			checkData  := mapset.NewSet()
			for i, _ := range dat {
				checkData.Add(i)
			}
			diffData := checkSlice.Difference(checkData)
			num      := len(diffData.ToSlice())
			if num > 0 {
				paramsName := fmt.Sprintf("%s",diffData)
				return dat,errors.New("missing required parameter:"+paramsName)
			}
		}
		return dat,nil
	}
}
// check params format
func ChectIntFloat(p interface{})(bool)  {
	switch p.(type) {
	case float64:
		pStr := strconv.FormatFloat(p.(float64),'f',-1,64)
		pattern := `^(\d+)*$` //反斜杠要转义
		res,_ := regexp.MatchString(pattern,pStr)
		if !res {
			return false
		}
	default:
		return false

	}
	return true
}

// 各种字符串加星
func HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str,"@") {
		res := strings.Split(str,"@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := Substr2(str,0,3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			result =  Substr2(str,0,3) + "****" + Substr2(str,7,11)
		} else {
			nameRune := []rune(str)
			lens  := len(nameRune)

			if  lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result =  string(nameRune[:1]) + "**" + string(nameRune[lens - 1 : lens])
			} else if lens > 4 {
				result =  string(nameRune[:2]) + "***" + string(nameRune[lens - 2 : lens])
			}
		}
		return
	}
}

func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}