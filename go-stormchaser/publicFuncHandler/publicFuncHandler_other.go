// +build !windows

package publicFuncHandler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"stormchaser/singleLog"
	"strings"
	"syscall"
	"time"

	"github.com/Han-Ya-Jun/qrcode2console"
)

const (
	Platform string = "linux"
)

//命令 (ping)
func ExecPING(argsStr string) ([]byte, error) {

	return ExecLinuxCommand("ping -c 4 " + argsStr)
}

//命令 (ps -aux | grep happ)
func ExecPsAndGrep(argsStr string) ([]byte, error) {

	return ExecLinuxCommand("ps -x | grep -i " + argsStr + " | grep -v grep")
}

//命令 (启动happ)
func ExecStartUpHapp(argsStr, path string) ([]byte, error) {

	return ExecLinuxCommand(path + "/" + argsStr + " " + path)
}

//添加执行权限
func AddExecutePermission(argsStr string) ([]byte, error) {

	return ExecLinuxCommand("chmod 0755 " + argsStr)
}

//添加扩展名
func AddExtensionName() string {

	return ""
}

//执行Linux命令
func ExecLinuxCommand(cmdStr string) ([]byte, error) {

	cmd := exec.Command("sh", "-c", cmdStr)
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
			_, err := ExecLinuxCommand("cp -f " + logFileName + " ./logs/" + getYMD() + ".out")
			if err == nil {
				ExecLinuxCommand("echo updateLog >" + logFileName)
			}
			timer.Reset(time.Second * 3 * 24 * 60 * 60)
		}
	}
}

//获取目录剩余空间 返回单位：byte
func DiskUsage(path string) (uint64, string, error) {

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(strings.TrimSpace(path), &fs)
	if err != nil {
		return 0, "", err
	}

	return fs.Bfree * uint64(fs.Bsize), path, nil
}

//在终端显示二维码
func ShowQRcode(qr *qrcode2console.QRCode2Console) {

	qr.Output()
}
