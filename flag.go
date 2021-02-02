package main

import (
	"fmt"
	pflag "github.com/spf13/pflag"
	"os"
	"runtime"
)

var (
	v bool
	h bool
)

// parseFlag 用于解析 Flag
func parseFlag() {
	pflag.BoolVarP(&Debug, "debug", "D", false, "启动调试模式")
	pflag.BoolVarP(&v, "version", "v", false, "查看版本信息")
	pflag.BoolVarP(&h, "help", "h", false, "查看程序帮助")
	pflag.Parse()
	if h {
		fmt.Printf(`Pixiv Fetcher v%s
使用: pixivFetcher [-Dhv]
选项：
`, Version)
		pflag.PrintDefaults()
		os.Exit(0)
	}
	if v {
		fmt.Printf("Pixiv Fetcher, A lightweight pixiv resource proxy. Authored by a632079\n版本: %s\n提交哈希: %s\n提交时间: %s 编译时间: %s\n编译环境: %s\n", Version, CommitHash, CommitTime, BuildTime, runtime.Version())
		os.Exit(0)
	}
}
