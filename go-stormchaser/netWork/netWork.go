// netWork project netWork.go
package netWork

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"stormchaser/aes"
	"stormchaser/config"
	"stormchaser/monitorHCDN"
	"stormchaser/publicFuncHandler"
	"strings"
)

const (
	apiHost         string = "https://verify.filestorm.info"
	checkSearchNode string = "/moacipfsVerify/verify/checkSearchNode"
)

//获取签名字符串
func getSignStr(argsDict map[string]string, keys []string) string {

	buf := bytes.Buffer{}
	buf.WriteString("@#filestorm@#stormcatcher")
	sort.Strings(keys)
	for _, value := range keys {
		if value != "sign" {
			if argsDict[value] != "" {
				buf.WriteString(argsDict[value])
				buf.WriteString("#$")
			}
		}
	}

	hashValue := sha1.New()
	hashValue.Write(buf.Bytes())
	hashStr := fmt.Sprintf("%x", hashValue.Sum(nil))

	return hashStr
}

func postNetRequest(methodName string, param map[string]string, paramKeys []string) ([]byte, error) {

	param["sign"] = getSignStr(param, paramKeys)

	data, mErr := json.Marshal(param)
	if mErr != nil {
		return nil, mErr
	}
	resp, rErr := http.Post(apiHost+methodName,
		"application/json",
		strings.NewReader(string(data)))
	if rErr != nil {
		return nil, rErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []byte(""), errors.New("网络请求错误")
	}
	rDatas, rErr1 := ioutil.ReadAll(resp.Body)
	if rErr1 != nil {
		return nil, rErr1
	}
	return rDatas, nil
}

/*
 *searchNodeParams - 唯一码@#是否正常（0/1）@#IP
 *formal - 1（正式）
 */
func CheckSearchNodeNetRequest() ([]byte, error) {

	formal := "1"
	uCode := config.UniqueCode
	state := monitorHCDN.GetSearchMiningState()
	ip := strings.Replace(publicFuncHandler.GetExternal(), "\n", "", -1)
	timestamp, _ := publicFuncHandler.GetTimestamp()

	aesStr, aErr := aes.AesOperate(uCode+"@#"+state+"@#"+ip+"@#"+fmt.Sprint(timestamp), true)
	if aErr != nil {
		return nil, aErr
	}

	return postNetRequest(checkSearchNode, map[string]string{"searchNodeParams": aesStr, "formal": formal}, []string{"searchNodeParams", "formal"})
}
