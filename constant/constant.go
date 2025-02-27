package constant

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Params struct {
	Proxy      string
	MainFolder string
	Host       string
	Port       string
	User       string
	Password   string
	RealTime   bool
}

var params Params

func GetParams() Params {
	return params
}
func init() {
	initProxy()
	initDir()
	initFfmpeg()
	initMysql()
}

func initProxy() {
	params.SetProxy(os.Getenv("PROXY"))
	if params.Proxy == "" && runtime.GOOS == "linux" {
		log.Fatalln("容器中未指定外部可用代理")
	}
	if params.Proxy == "" {
		params.Proxy = "http://127.0.0.1:8889"
	}
	if err := ping(params.GetProxy()); err != nil {
		log.Fatalf("指定的代理IP地址不可用,错误信息:%v\n", err)
	}
}

func initDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("无法获取用户的个人文件夹目录:%v", err))
	} else {
		params.SetMainFolder(filepath.Join(home, "Downloads", "media"))
	}
	if tdl := os.Getenv("TDL"); tdl != "" {
		params.SetMainFolder(tdl)
	}
	if runtime.GOARCH == "arm64" && runtime.GOOS == "android" {
		params.SetMainFolder("/sdcard/.telegram")
		log.Printf("在termux上运行,下载位置为%v\n", params.GetMainFolder())
	}
}
func initFfmpeg() {
	if p := os.Getenv("REALTIME"); strings.ToUpper(p) == "TRUE" || strings.ToUpper(p) == "1" || strings.ToUpper(p) == "YES" || strings.ToUpper(p) == "Y" {
		params.SetRealTime(true)
	} else {
		params.SetRealTime(false)
	}
}

const (
	DEFAULT_MYSQL_USER     = "root"
	DEFAULT_MYSQL_PASSWORD = "163453"
	DEFAULT_MYSQL_HOST     = "192.168.2.5"
	DEFAULT_MYSQL_PORT     = "3307"
)

func initMysql() {
	if port := os.Getenv("MYSQL_PORT"); port != "" {
		params.SetPort(port)
	} else {
		params.SetPort(DEFAULT_MYSQL_PORT)
	}
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		params.SetHost(host)
	} else {
		params.SetHost(DEFAULT_MYSQL_HOST)
	}
	if user := os.Getenv("MYSQL_USER"); user != "" {
		params.SetUser(user)
	} else {
		params.SetUser(DEFAULT_MYSQL_USER)
	}
	if password := os.Getenv("MYSQL_PASSWORD"); password != "" {
		params.SetPassword(password)
	} else {
		params.SetPassword(DEFAULT_MYSQL_PASSWORD)
	}

}
func (p *Params) GetRealTime() bool {
	return p.RealTime
}
func (p *Params) SetRealTime(b bool) {
	p.RealTime = b
}
func (p *Params) SetHost(s string) {
	p.Host = s
}

func (p *Params) GetHost() string {
	return p.Host
}

func (p *Params) SetPort(s string) {
	p.Port = s
}

func (p *Params) GetPort() string {
	return p.Port
}

func (p *Params) SetUser(s string) {
	p.User = s
}

func (p *Params) GetUser() string {
	return p.User
}

func (p *Params) SetPassword(s string) {
	p.Password = s
}

func (p *Params) GetPassword() string {
	return p.Password
}

func (p *Params) GetProxy() string {
	return p.Proxy
}

func (p *Params) SetProxy(s string) {
	p.Proxy = s
}

func (p *Params) GetMainFolder() string {
	return p.MainFolder
}

func (p *Params) SetMainFolder(s string) {
	p.MainFolder = s
}

type OneFile struct {
	Channel  string // 频道id
	FileId   int    // 文件id
	Tag      string // 主文件夹名 #后面的文件名
	Subtag   string // 子(二级)文件夹名 &后面的文件名
	FileName string // 手动设置的文件名 @后面的文件名
	Offset   int    // 偏移量 如下载当前媒体之后第n个文件 +后面的数字
	Capacity int    // 下载当前文件和之后的n个文件 %后面的数字
	Success  bool   // 是否下载成功
}

func (f *OneFile) SetChannel(s string) {
	f.Channel = s
}

func (f *OneFile) SetId(i int) {
	f.FileId = i
}

func (f *OneFile) SetTag(s string) {
	f.Tag = s
}

func (f *OneFile) SetSubtag(s string) {
	f.Subtag = s
}

func (f *OneFile) SetFileName(s string) {
	f.FileName = s
}

func (f *OneFile) SetOffset(i int) {
	f.Offset = i
}

func (f *OneFile) SetCapacity(i int) {
	f.Capacity = i
}

func (f *OneFile) SetStatus() {
	f.Success = true
}

func ping(proxy string) error {
	u, err := url.Parse(proxy)
	if err != nil {
		fmt.Println("解析URL失败:", err)
		return err
	}
	ip := u.Hostname()
	port := u.Port()

	address := net.JoinHostPort(ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
