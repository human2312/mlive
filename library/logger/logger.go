package logger

import (
	"github.com/spf13/viper"
	"log"
	"mlive/library/logger/lumberjack"
	"os"
	"strings"
)

var (
	info   *log.Logger
	notice *log.Logger
	warn   *log.Logger
	error  *log.Logger

	lReq *log.Logger

	file    string
	dir     string
	fileReq string
	bakDir  string
)

func Init() {
	file = viper.GetString("log.file")
	fileReq = viper.GetString("log.fileRequest")
	dir = viper.GetString("log.dir")
	bakDir = viper.GetString("log.bakDir")

	if err := os.MkdirAll(bakDir, 0777); err != nil {
		panic(err)
	}

	// 递归创建文件夹
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(err)
	}
	// 文件路径
	var (
		path    string
		pathReq string
	)
	if strings.HasSuffix(dir, "/") {
		path = dir + file
		pathReq = dir + fileReq
	} else {
		path = dir + "/" + file
		pathReq = dir + "/" + fileReq
	}
	// f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	f := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    viper.GetInt("log.maxSize"), // megabytes
		MaxBackups: viper.GetInt("log.maxBackups"),
		MaxAge:     0,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
		BakDir:     bakDir,
	}

	// fReq, err := os.OpenFile(pathReq, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	fReq := &lumberjack.Logger{
		Filename:   pathReq,
		MaxSize:    viper.GetInt("log.maxSize"), // megabytes
		MaxBackups: viper.GetInt("log.maxBackups"),
		MaxAge:     0,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
		BakDir:     bakDir,
	}

	info = log.New(f, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	notice = log.New(f, "[NOTICE] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	warn = log.New(f, "[WARN] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	error = log.New(f, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds)

	lReq = log.New(fReq, "[GIN] ", log.Ldate|log.Ltime|log.Lmicroseconds)

}

func Rprintln(v ...interface{}) {
	lReq.Println(v...)
}

func Iprintln(v ...interface{}) {
	info.Println(v...)
}

func Nprintln(v ...interface{}) {
	notice.Println(v...)
}

func Wprintln(v ...interface{}) {
	warn.Println(v...)
}

func Eprintln(v ...interface{}) {
	error.Println(v...)
}

func Eprintf(format string, v ...interface{}) {
	error.Printf(format, v...)
}
