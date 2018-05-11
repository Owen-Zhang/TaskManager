package models

import "runtime"

const (
    TempDir  = "Data/Temp" 
    RunDir   = "Data/Run" 
	CNTimeFormat = "2006-01-02 15:04:05"
	SystemWindows = "windows"
	SystemLinux = "linux"
	WindowsShellExt = "bat"
	LinuxShellExt = "sh"
)

var Common *CommonInfo

type CommonInfo struct {
	SystemName string   //系统名称(windows, linux)
}

func init() {
	Common = &CommonInfo{
		SystemName : runtime.GOOS,
	}
}
