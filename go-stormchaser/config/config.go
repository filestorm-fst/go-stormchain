// config project config.go
package config

import (
	"os"
	"stormchaser/publicFuncHandler"
	"stormchaser/singleLog"
	"strings"

	"github.com/Han-Ya-Jun/qrcode2console"
	"github.com/Unknwon/goconfig"
)

const (
	configName string = "/stormchaser_config"
)

var (
	cfg        *goconfig.ConfigFile = nil
	configPath string               //配置文件路径

	UniqueCode string //唯一码
)

func init() {

	configPath = publicFuncHandler.GetPath() + configName
	_, osErr := os.Stat(configPath)
	if osErr != nil && os.IsNotExist(osErr) {
		file, fErr := os.Create(configPath)
		if fErr != nil {
			if strings.Contains(fErr.Error(), "Permission denied") {
				singleLog.GetInstance().Error("Configuration file create fail! The first run needs to be run by the root user!")
			} else {
				singleLog.GetInstance().Error("Configuration file create fail!", fErr)
			}
			os.Exit(0)
		}
		file.Close()
		if cErr := os.Chmod(configPath, 0666); cErr != nil {
			os.Remove(configPath)
			singleLog.GetInstance().Error("file attributes modify fail!")
			os.Exit(0)
		}
	}

	var lErr error
	cfg, lErr = goconfig.LoadConfigFile(configPath)
	if lErr != nil {
		singleLog.GetInstance().Error("Configuration file load fail!", lErr)
		os.Exit(0)
	}
}

func LoadConfig() string {

	var gErr error
	UniqueCode, gErr = GetConfig("UniqueCode")
	if gErr != nil && strings.Contains(gErr.Error(), "not found") {
		UniqueCode = publicFuncHandler.CreateUniqueCode()
		if sErr := SaveConfig("UniqueCode", UniqueCode); sErr != nil {
			singleLog.GetInstance().Error("UniqueCode save fail!", sErr)
			os.Exit(0)
		}
	}
	reStr := UniqueCode + "@#" + publicFuncHandler.GetExternal()
	publicFuncHandler.ShowQRcode(qrcode2console.NewQRCode2ConsoleWithUrl(reStr, true))

	return reStr
}

//获取配置文件key对应的value
func GetConfig(key string) (string, error) {

	return cfg.GetValue("DEFAULT", key)
}

//保存key-value到配置文件
func SaveConfig(key, value string) error {

	cfg.SetValue("DEFAULT", key, value)
	return goconfig.SaveConfigFile(cfg, configPath)
}
