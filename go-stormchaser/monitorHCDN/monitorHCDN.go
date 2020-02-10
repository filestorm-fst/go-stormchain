// monitorHCDN project monitorHCDN.go
package monitorHCDN

import (
	"stormchaser/publicFuncHandler"
	"stormchaser/singleLog"
	"strings"
	"sync"
	"time"
)

var countMutex sync.Mutex

var (
	pingNum      int = 0 //ping成功次数
	detectionNum int = 0 //检测HCDN状态成功次数
)

//检测爱奇艺链接是否能ping通
func DetectionPing() {

	var aqiyiArray []string = []string{"pdata.video.qiyi.com", "stat.hcdn.qiyi.com"}

	timer := time.NewTimer(time.Second * 10 * 60)
	for {
		select {
		case <-timer.C:
			var passNum int = 0
			for index, value := range aqiyiArray {
				dataBytes, err := publicFuncHandler.ExecPING(value)
				if err != nil {
					singleLog.GetInstance().Error("ping错误：", index, err)
					continue
				}
				if strings.Count(string(dataBytes), "timeout") <= 2 && strings.Count(string(dataBytes), "超时") <= 2 {
					//通过
					passNum += 1
				} else {
					//未通过
					singleLog.GetInstance().Error("ping超时：", index)
				}
			}
			if passNum == len(aqiyiArray) {
				//正常
				countMutex.Lock()
				pingNum += 1
				countMutex.Unlock()
			} else {
				//异常
				singleLog.GetInstance().Error("网络异常（网络波动较大），请检查网络是否链接！")
			}
		}
		timer.Reset(time.Second * 10 * 60)
	}
}

//检测HCDN状态
func DetectionHcdnState(path string) {

	str := "HAPP"
	lStr := strings.ToLower(str)

	timer := time.NewTimer(time.Second * 5 * 60)
	for {
		select {
		case <-timer.C:
			dataBytes, err := publicFuncHandler.ExecPsAndGrep(lStr)
			if err != nil {
				publicFuncHandler.ExecStartUpHapp(lStr, path)
				singleLog.GetInstance().Error("检测状态错误：", err)
				continue
			}
			strs := strings.Split(string(dataBytes), lStr)
			if len(strs) >= 2 {
				//运行中
				countMutex.Lock()
				detectionNum += 1
				countMutex.Unlock()
			} else {
				//已停止运行，需启动happ
				publicFuncHandler.ExecStartUpHapp(lStr, path)
			}
		}
		timer.Reset(time.Second * 5 * 60)
	}
}

func GetSearchMiningState() string {

	if pingNum >= 14 && detectionNum >= 14 {
		return "1"
	}
	return "0"
}

func ClearCount() {

	countMutex.Lock()
	defer countMutex.Unlock()
	pingNum = 0
	detectionNum = 0
}
