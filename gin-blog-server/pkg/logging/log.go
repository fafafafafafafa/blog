package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	log.Println("logging Setup...")
	filePath := getLogFilePath()
	fileName := getLogFileName()

	var err error
	F, err = openLogFile(filePath, fileName)
	if err != nil {
		log.Fatalln(err)
	}
	// F 要写入日志的句柄
	// DefaultPrefix 定义每个生成的日志行的开头
	// flag定义了日志记录属性
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

// 根据level 设置前缀
func setPrefix(level Level) {
	// DefaultCallerDepth=0：Caller()会报告Caller()的调用者的信息
	// DefaultCallerDepth=1：Caller()回报告Caller()的调用者的调用者的信息
	// 用于回报 调用的层数
	// 	file：带路径的完整文件名
	// line：该调用在文件中的行号
	// ok：是否可以获得信息
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		// filepath.Base 返回路径的最后一个名称
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}
