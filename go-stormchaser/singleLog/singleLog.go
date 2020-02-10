package singleLog

import (
	"os"
	//	"runtime"
	"sync"

	logging "github.com/op/go-logging"
)

var log *logging.Logger
var once sync.Once

func GetInstance() *logging.Logger {

	//	_, file, line, ok := runtime.Caller(1)
	//	fmt.Println(file, line, ok)

	once.Do(func() {
		//指定是否控制台打印，默认为true
		log = setupLogging(true)
	})
	return log
}

func setupLogging(hasColor bool) *logging.Logger {

	//setup logging
	log, _ := logging.GetLogger("Storm Catcher")
	// `%{time:2006/01/02 15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{id:03x} %{message}`,
	format, _ := logging.NewStringFormatter(
		`%{time:2006/01/02 15:04:05.000} ▶ %{level:.5s} %{id:03x} %{message}`,
	)
	if hasColor {
		format, _ = logging.NewStringFormatter(
			`%{color}%{time:2006/01/02 15:04:05.000} ▶ %{level:.5s} %{id:03x}%{color:reset} %{message}`,
		)
	}

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.CRITICAL, "")
	logging.SetBackend(backendLeveled, backendFormatter)

	return log
}
