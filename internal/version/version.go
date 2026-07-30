package version

import (
	"fmt"
	"runtime"
)

var (
	// Version 是 CLI 的当前版本（在构建时注入）
	Version = "0.1.0"

	// Commit 是 Git 提交哈希（在构建时注入）
	Commit = "unknown"

	// Date 是构建日期（在构建时注入）
	Date = "unknown"
)

// Info 包含版本信息
type Info struct {
	Version string
	Commit  string
	Date    string
	GoOS    string
	GoArch  string
}

// Get 返回版本信息
func Get() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
		GoOS:    runtime.GOOS,
		GoArch:  runtime.GOARCH,
	}
}

// String 返回格式化后的版本字符串
func (i Info) String() string {
	return fmt.Sprintf("hypersku-cli v%s (commit=%s, built=%s, %s/%s)",
		i.Version, i.Commit, i.Date, i.GoOS, i.GoArch)
}
