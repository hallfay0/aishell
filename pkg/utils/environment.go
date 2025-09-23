package utils

import (
	"os"
	"runtime"
	"time"
)

// GetEnvironmentInfo 获取当前环境信息
func GetEnvironmentInfo() (currentDir, currentOS, currentArch, currentTime string) {
	currentDir, _ = os.Getwd()
	currentOS = runtime.GOOS
	currentArch = runtime.GOARCH
	currentTime = time.Now().Format("2006-01-02 15:04:05")
	return
}
