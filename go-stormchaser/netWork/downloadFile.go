// netWork project fileDownload.go
package netWork

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

const (
	bufSize int = 1024 * 1024 //buf大小
)

//下载文件，返回文件长度
func DownloadFile(fileUrl, filePath string) (int64, error) {

	res, gErr := http.Get(fileUrl)
	if gErr != nil {
		return 0, gErr
	}
	defer res.Body.Close()

	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, bufSize)

	file, cErr := os.Create(filePath)
	if cErr != nil {
		return 0, cErr
	}
	defer file.Close()
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, coErr := io.Copy(writer, reader)
	if coErr != nil {
		return 0, coErr
	}

	return written, writer.Flush() //刷新后内容写入
}
