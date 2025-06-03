package xlog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const RequestIdKey = "request_id"

var LogLevel = "debug"
var LogFile = ""

func Init() {
	err := os.Mkdir("logs", os.ModePerm)
	if err != nil {
		fmt.Printf("无法创建logs目录: %v\n", err)
	}
	if LogFile != "" {
		f, err := os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("open log file err")
		}
		log.SetOutput(f)
	}
}

type FormatOutput struct {
	Level    string    `json:"level"`
	Filename string    `json:"filename"`
	Line     int       `json:"line"`
	Content  string    `json:"content"`
	Now      time.Time `json:"ts"`
}

func NewFormatOutput(level string, filename string, line int, content string) FormatOutput {
	return FormatOutput{
		Now:      time.Now(),
		Level:    level,
		Filename: filename,
		Line:     line,
		Content:  content,
	}
}
func (fo FormatOutput) String() string {
	b, _ := json.Marshal(fo)
	return string(b)
}

func Errorf(format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("ERROR", filename, line, content))
}

func Warnf(format string, args ...interface{}) {
	if LogLevel != "debug" && LogLevel != "warn" {
		return
	}
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("WARN", filename, line, content))
}

func Debugf(format string, args ...interface{}) {
	if LogLevel != "debug" {
		return
	}
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("DEBUG", filename, line, content))
}

func Infof(format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("INFO", filename, line, content))
}

func Error(args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf("%v", args...)
	log.Printf("%s\n", NewFormatOutput("ERROR", filename, line, content))
}

func Debug(args ...interface{}) {
	if LogLevel != "debug" {
		return
	}
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf("%v", args...)
	log.Printf("%s\n", NewFormatOutput("DEBUG", filename, line, content))
}

func Info(args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf("%v", args...)
	log.Printf("%s\n", NewFormatOutput("INFO", filename, line, content))
}

func ErrorfN(n int, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(n)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("ERROR", filename, line, content))
}

func DebugfN(n int, format string, args ...interface{}) {
	if LogLevel != "debug" {
		return
	}
	_, filename, line, _ := runtime.Caller(n)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("DEBUG", filename, line, content))
}

func InfofN(n int, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(n)
	content := fmt.Sprintf(format, args...)
	log.Printf("%s\n", NewFormatOutput("INFO", filename, line, content))
}

func ErrorN(n int, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(n)
	content := fmt.Sprintf("%v", args...)
	log.Printf("%s\n", NewFormatOutput("ERROR", filename, line, content))
}

func DebugN(n int, args ...interface{}) {
	if LogLevel != "debug" {
		return
	}
	_, filename, line, _ := runtime.Caller(n)
	content := fmt.Sprintf("%v", args...)
	log.Printf("%s\n", NewFormatOutput("DEBUG", filename, line, content))
}

func InfoN(n int, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(n)
	content := fmt.Sprintf("%v", args...)
	log.Printf("%s\n", NewFormatOutput("INFO", filename, line, content))
}

func CtxErrorf(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	// function name: runtime.FuncForPC(pc).Name()
	log.Printf("[ERROR] in [%s:%d] [request_id:%s] %v", filename, line, ctx.Value(RequestIdKey), content)
}

func CtxDebugf(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	// function name: runtime.FuncForPC(pc).Name()
	log.Printf("[DEBUG] in [%s:%d][request_id:%s]  %v", filename, line, ctx.Value(RequestIdKey), content)
}

func CtxInfof(ctx context.Context, format string, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	content := fmt.Sprintf(format, args...)
	// function name: runtime.FuncForPC(pc).Name()
	log.Printf("[INFO] in [%s:%d] [request_id:%s] %v", filename, line, ctx.Value(RequestIdKey), content)
}

func CtxDebug(ctx context.Context, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[DEBUG] in [%s:%d] [request_id:%s] %v", filename, line, ctx.Value(RequestIdKey), args)
}

func CtxInfo(ctx context.Context, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[INFO] in [%s:%d] [request_id:%s] %v", filename, line, ctx.Value(RequestIdKey), args)
}

func CtxError(ctx context.Context, args ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[ERROR] in [%s:%d] [request_id:%s] %v", filename, line, ctx.Value(RequestIdKey), args)
}
