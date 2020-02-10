// +build windows

package publicFuncHandler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"stormchaser/singleLog"
	"syscall"
	"time"
	"unsafe"

	"github.com/Han-Ya-Jun/qrcode2console"
)

const (
	Platform string = "windows"
)

//命令 (ping)
func ExecPING(argsStr string) ([]byte, error) {

	return ExecLinuxCommand("ping -n 4 " + argsStr)
}

//命令 (tasklist | findstr happ)
func ExecPsAndGrep(argsStr string) ([]byte, error) {

	return ExecLinuxCommand("tasklist | findstr -i " + argsStr)
}

//命令 (启动happ)
func ExecStartUpHapp(argsStr, path string) ([]byte, error) {

	return ExecLinuxCommand(path + "\\" + argsStr + ".exe " + path)
}

//添加执行权限
func AddExecutePermission(argsStr string) ([]byte, error) {

	return nil, nil
}

//添加扩展名
func AddExtensionName() string {

	return ".exe"
}

//执行windows命令
func ExecLinuxCommand(cmdStr string) ([]byte, error) {

	cmd := exec.Command("cmd", "/C", cmdStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		singleLog.GetInstance().Error("ExecLinuxCommand : ", err.Error(), stderr.String())
	}
	return out.Bytes(), err
}

//日志切割
func CuttingLogFile(logFileName string) {

	getYMD := func() string {

		t := time.Now()
		return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	}

	//创建logs目录
	os.Mkdir("logs", os.ModeDir)

	timer := time.NewTimer(time.Second * 3 * 24 * 60 * 60)
	for {
		select {
		case <-timer.C:
			_, err := ExecLinuxCommand("copy " + logFileName + " logs\\" + getYMD() + ".out")
			if err == nil {
				ExecLinuxCommand("echo updateLog >" + logFileName)
			}
			timer.Reset(time.Second * 3 * 24 * 60 * 60)
		}
	}
}

//获取目录剩余空间 返回单位：byte
func DiskUsage(path string) (uint64, string, error) {

	type diskusage struct {
		Path  string `json:"path"`
		Total uint64 `json:"total"`
		Free  uint64 `json:"free"`
	}

	var kernel = syscall.NewLazyDLL("Kernel32.dll")
	var getDiskFreeSpaceExW = kernel.NewProc("GetDiskFreeSpaceExW")

	usage := func(getDiskFreeSpaceExW *syscall.LazyProc, path string) (diskusage, error) {
		lpFreeBytesAvailable := int64(0)
		var info = diskusage{Path: path}
		diskret, _, err := getDiskFreeSpaceExW.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(info.Path))),
			uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
			uintptr(unsafe.Pointer(&(info.Total))),
			uintptr(unsafe.Pointer(&(info.Free))))
		if diskret != 0 {
			err = nil
		}
		return info, err
	}

	info, err := usage(getDiskFreeSpaceExW, path)
	if err != nil {
		return 0, "", err
	}

	return info.Free, path, nil
}

//在终端显示二维码
func ShowQRcode(qr *qrcode2console.QRCode2Console) {

	//windows还没有找到在终端显示二维码的方式
}
