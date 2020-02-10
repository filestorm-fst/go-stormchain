package publicFuncHandler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"stormchaser/singleLog"
	"time"

	"github.com/rs/xid"
)

//获取外网IP
func GetExternal() string {

	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	//s := buf.String()
	return string(content)
}

/*
//局域网测试 - 获取局域网IP
func GetExternal() string {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}*/

//获取时间戳，以上海时区时间为准
func GetTimestamp() (int64, error) {

	return time.Now().Unix(), nil
}

//创建唯一码
func CreateUniqueCode() string {

	tmpXid := xid.New()
	str1 := tmpXid.String()
	str2 := fmt.Sprint(time.Now().UnixNano())

	var reStr string
	for index, value := range str1 {
		if len(str2) > index {
			reStr = fmt.Sprintf("%s%c%c", reStr, value, str2[index])
		} else {
			reStr = fmt.Sprintf("%s%c", reStr, value)
		}
	}
	if len(reStr) > 32 {
		reStr = reStr[:32]
	}
	return reStr
}

//获取当前路径
func GetPath() string {

	dir, aErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if aErr != nil {
		singleLog.GetInstance().Error("Unknown directory path")
		os.Exit(0)
	}
	return dir
}
