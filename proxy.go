package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
)

// UserAgent 列表
// TODO: 定期更新
var userAgentList = []string{
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
	"Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:46.0) Gecko/20100101 Firefox/46.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:46.0) Gecko/20100101 Firefox/46.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/13.10586",
}

// 对 []string 包装一下
type _MIMEList []string

func (p *_MIMEList) Exist(s string) bool {
	for _, v := range *p {
		if v == s {
			return true
		}
	}
	return false
}

// 有效的 Mime 列表
var validMIMEList = _MIMEList{
	// 图片
	"image/gif",
	"image/png",
	"image/jpeg",
	"image/bmp",
	"image/webp",
	"image/x-icon",
	"image/vnd.microsoft.icon",
	// 音频
	"audio/midi",
	"audio/mpeg",
	"audio/webm",
	"audio/ogg",
	"audio/wav",
	// 视频
	"video/webm",
	"video/ogg",
}

// getRandomUserAgent 获得随机 UA
func getRandomUserAgent() string {
	length := len(userAgentList)
	return userAgentList[rand.Intn(length)]
}

// acceptedHostReg 有效的主机名正则
var acceptedHostReg = []string{
	"(.*).pixiv.(.*)",
	"(.*).pximg.net",
}

// checkHost 检测主机名
func checkHost(URL string) (bool, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return false, err
	}
	var t bool
	for _, r := range acceptedHostReg {
		if t {
			break
		}
		reg := regexp.MustCompile(r)
		t = len(reg.FindAllString(u.Host, -1)) != 0
	}
	return t, nil
}

// fetchPixivResources 用于拉取资源
func fetchPixivResources(URL string) (io.ReadCloser, string, error) {
	// 检测拉取主机名是否有效
	if pass, err := checkHost(URL); err != nil {
		return nil, "", err
	} else if !pass {
		return nil, "", errors.New("host is not a valid pixiv domain")
	}
	// 发起请求
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, "", err
	}
	// 设置请求 Header
	req.Header.Set("Referer", "http://www.pixiv.net/")
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Accept-Encoding", "")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != 200 { // 状态码不为 200
		log.Error("[proxy.fetch] 资源状态码不为 200")
		log.Error(resp)
		return nil, "", errors.New("the status code is not 200")
	}
	// 检测 Mime
	mime := resp.Header.Get("Content-Type")
	if !validMIMEList.Exist(mime) {
		resp.Body.Close()
		err = fmt.Errorf("[proxy] the url(%s) you required is not a valid resource. The MIME that we know is `%s`", URL, mime)
		log.Error(err.Error())
		return nil, mime, err
	}
	return resp.Body, mime, nil
}
