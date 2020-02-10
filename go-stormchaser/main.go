// stormchaser project main.go
package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"stormchaser/config"
	"stormchaser/monitorHCDN"
	"stormchaser/netWork"
	"stormchaser/publicFuncHandler"
	"stormchaser/singleLog"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {

	timer := time.NewTimer(time.Second * 3 * 60 * 60)
	for {
		select {
		case <-timer.C:
			netWork.CheckSearchNodeNetRequest()
		}
		timer.Reset(time.Second * 3 * 60 * 60)
	}
}

func (p *program) Stop(s service.Service) error {
	return nil
}

var searchId string

func main() {

	var path string
	flag.StringVar(&path, "path", "", "检索矿工缓存目录")
	flag.Parse()
	if path == "" {
		path = publicFuncHandler.GetPath() //当前目录
	}

	searchId = config.LoadConfig()

	go publicFuncHandler.CuttingLogFile("storm.out") //日志切割
	go monitorHCDN.DetectionPing()                   //检测网络（ping）
	go monitorHCDN.DetectionHcdnState(path)          //检测HCDN状态

	hardDiskDetection(path) //硬盘检测

	svcConfig := &service.Config{
		Name:        "stormchaser",
		DisplayName: "filestorm 检索矿工",
		Description: "filestorm 检索矿工",
	}
	prg := &program{}
	sys := service.ChosenSystem()
	s, nErr := sys.New(prg, svcConfig)
	if nErr != nil {
		singleLog.GetInstance().Error("Server start fail", nErr)
		os.Exit(0)
	}
	go s.Run()

	http.HandleFunc("/search/getSearchMiningIdAndIp", getSearchMiningIdAndIpHandler)
	http.ListenAndServe(":52530", nil)
}

//硬盘检测
func hardDiskDetection(path string) {

	_, osErr := os.Stat(path + "/hdata")
	if osErr != nil && os.IsNotExist(osErr) {

		//下载happ安装包
		dErr := downloadHapp("https://file.filestorm.info/management", path)
		if dErr != nil {
			singleLog.GetInstance().Error("File download operating fail", dErr)
			os.Exit(0)
		}

		//hdata文件夹不存在，需检测硬盘空间是否够（200GB）
		size, _, dErr := publicFuncHandler.DiskUsage(path)
		if dErr != nil {
			singleLog.GetInstance().Error("路径错误：", dErr)
			os.Exit(0)
		}
		if size < 200*1024*1024*1024 {
			//硬盘存储空间足，至少可用200G
			singleLog.GetInstance().Error("硬盘存储空间足，至少可用200G")
			os.Exit(0)
		}
	}
}

func downloadHapp(urlStr, pathStr string) (dErr error) {

	if publicFuncHandler.Platform == "linux" {
		filePath := pathStr + "/happ"
		_, dErr = netWork.DownloadFile(urlStr+"/linux/happ", filePath)
		if dErr == nil {
			_, dErr = publicFuncHandler.AddExecutePermission(filePath)
		}
	} else if publicFuncHandler.Platform == "windows" {
		_, dErr = netWork.DownloadFile(urlStr+"/windows/happ.exe", pathStr+"/happ.exe")
		if dErr == nil {
			_, dErr = netWork.DownloadFile(urlStr+"/windows/happnet.dll", pathStr+"/happnet.dll")
		}
	} else {
		dErr = errors.New("Paltform error")
	}
	return dErr
}

func getSearchMiningIdAndIpHandler(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte(searchId))
}
