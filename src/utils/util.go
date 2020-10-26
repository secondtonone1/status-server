package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"status-server/logging"
	"strconv"
	"time"
)

func Zipmap(a, b []string) (map[string]string, error) {

	if len(a) != len(b) {
		return nil, errors.New("zip: arguments must be of same length")
	}

	r := make(map[string]string, len(a))

	for i, e := range a {
		r[e] = b[i]
	}

	return r, nil
}

//随机数
func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

//生成32位md5值
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

//格式化时间
func FormatTime(t time.Time) string {
	f := "2006-01-02 15:04:05"
	return t.Format(f)
}

//字符串转换为时间戳
func Stime2TimeStamp(stime string) int64 {
	f := "2006-01-02 15:04:05"
	t, e := time.ParseInLocation(f, stime, time.Local)
	if e != nil {
		logging.Logger.Errorf("string to time stamp stime=[%s],error=[%s]", stime, e)
	}
	return t.UnixNano() / 1e6
}

//获取当前时间戳，为毫秒
func TimeStamp() int64 {
	return time.Now().UnixNano() / 1e6
}

//创建token
func CreateToken(code, accounts, uid, time string) string {
	return Md5(code + accounts + uid + time)
}

func StrToUint(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseUint(strNumber, 10, 64)
	switch v := number.(type) {
	case uint64:
		switch d := value.(type) {
		case *uint64:
			*d = v
		case *uint:
			*d = uint(v)
		case *uint16:
			*d = uint16(v)
		case *uint32:
			*d = uint32(v)
		case *uint8:
			*d = uint8(v)
		}
	}
	return
}
