package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
)

type Configuration struct {
	SCSService             bool
	ShowToPublic           bool
	VnodeServiceCfg        *string
	VnodeBeneficialAddress *string
	VnodeIp                *string
	ForceSubnetP2P         []string
}

// GetConfiguration: read config from .json file
// 1) No config file, using default value, don't create new file;
// 2) has config file, error in reading config, stop and display correct info;
// 3) has config file, no error in reading config, go ahead load and run;
func GetConfiguration(configFilePath string) (*Configuration, error) {
	//check the file path
	filePath, _ := filepath.Abs(configFilePath)
	log.Debugf("load vnode config file from %v", filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		//If no config file exists, return nil
		log.Infof("%v not exists\nUse default settings", configFilePath)
		return nil, nil
	}

	if _, err := os.Stat(configFilePath); err != nil {

		log.Errorf("Open %v error: \n%v\n", configFilePath, err)

		return nil, err
	}

	return GetUserConfig(configFilePath)
}

//GetUserConfig -- user configuration in json file
func GetUserConfig(filepath string) (*Configuration, error) {
	file, ferr := os.Open(filepath)
	defer file.Close()
	conf := Configuration{}

	if ferr != nil {
		log.Errorf("%v open error:%v\n", filepath, ferr)
		return nil, ferr
	}

	decoder := json.NewDecoder(file)

	err := decoder.Decode(&conf)
	if err != nil {
		log.Infof("\nConfig file JSON decode Error:%v\n", err)
		hint := "Config file should have following variables:\n{\n  \"SCSService\":false,\n  " +
			"\"ShowToPublic\":false,\n  \"VnodeServiceCfg\":\":50062\",\n  \"VnodeIp\":\"\",\n  \"VnodeBeneficialAddress\":\"\",\n" +
			"  \"ForceSubnetP2P\": []\n}"
		log.Info(hint)
		return nil, err
	}

	return &conf, nil
}

//SaveUserConifg -- Save the input configuration in the json file
//discarded -
func SaveUserConifg(outconfig *Configuration) bool {

	filepath, _ := filepath.Abs("./vnodeconfig.json")

	outfile, ferr := os.Create(filepath)
	defer outfile.Close()

	if ferr != nil {
		//Cannot open the vnodeconfig.json file using the input filepath, save to default
		log.Infof("Cannot open %v configuration file!", filepath)
		return false
	}
	log.Infof("Save VNODE configuration to %v!", filepath)
	outjson := ""
	if outconfig.VnodeBeneficialAddress == nil {
		outjson = fmt.Sprintf("{\n  \"SCSService\":%v,\n  \"ShowToPublic\":%v,\n  \"VnodeServiceCfg\":\"%v\",\n  \"VnodeIp\":\"%v\",\n  \"VnodeBeneficialAddress\":\"\"\n}", outconfig.SCSService, outconfig.ShowToPublic, *outconfig.VnodeServiceCfg, *outconfig.VnodeIp)
	} else {
		outjson = fmt.Sprintf("{\n  \"SCSService\":%v,\n  \"ShowToPublic\":%v,\n  \"VnodeServiceCfg\":\"%v\",\n  \"VnodeIp\":\"%v\",\n  \"VnodeBeneficialAddress\":\"%v\"\n}", outconfig.SCSService, outconfig.ShowToPublic, *outconfig.VnodeServiceCfg, *outconfig.VnodeIp, *outconfig.VnodeBeneficialAddress)
	}

	_, err := outfile.WriteString(outjson)
	if err != nil {
		log.Info("Write output configuration file Error:", err)
		return false
	}
	return true
}
