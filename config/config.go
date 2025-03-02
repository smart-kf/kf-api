package config

import (
	"fmt"
)

type Config struct {
	Debug         bool          `json:"debug"`
	Web           Web           `json:"web"`
	Log           Log           `json:"log"`
	DB            Db            `json:"db"`
	LevelDBConfig LevelDBConfig `json:"levelDB"`
	JwtKey        string        `json:"jwtKey"`
	BillConfig    BillConfig    `json:"billConfig"`
	RedisConfig   RedisConfig   `json:"redis"`
	NSQ           NSQ           `json:"nsq"`
	HttpClient    HttpClient    `json:"httpClient"`
	CardPackages  []CardPackage `json:"cardPackages"`
	Payment       Payment       `json:"payment"`
	DomainPrice   float64       `json:"domainPrice"`
	Ip2Region     Ip2Region     `json:"ip2Region"`
	SocketIO      SocketIO      `json:"socketIo"`
}

type Ip2Region struct {
	XDBPath        string `json:"xdbPath"`
	RegistryApiKey string `json:"registryApiKey"`
	Proxy          string `json:"proxy"`
	Timeout        int    `json:"timeout"`
}

type LevelDBConfig struct {
	Path string `json:"path"`
}

type Payment struct {
	Host     string `json:"host"`
	Token    string `json:"token"`
	AppId    string `json:"appId"`
	FromMail string `json:"fromMail"` // 发邮件的来源地址
}

type Web struct {
	Addr             string `json:"addr" default:"127.0.0.1"`
	Port             int    `json:"port" default:"8081"`
	StaticDir        string `json:"staticDir" default:"static"`
	CdnHost          string `json:"cdnHost"`          // 用户上传的文件cdn域名
	UploadDir        string `json:"uploadDir"`        // 上传文件夹
	KfManagerAddress string `json:"kfManagerAddress"` // 客服登录地址
}

func (w Web) String() string {
	return fmt.Sprintf("%s:%d", w.Addr, w.Port)
}

type Db struct {
	Dsn    string `json:"dsn"`    // 连接
	Driver string `json:"driver"` // 默认 sqlite3
}

type Log struct {
	Level  string `json:"level" default:"info"`
	Format string `json:"format" default:"json"`
	File   string `json:"file"`
}

type BillConfig struct {
	OrderExpireTime int64         `json:"orderExpireTime" default:"600"` // 默认10分钟过期
	Accounts        []BillAccount `json:"accounts"`
}

type BillAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisConfig struct {
	DB       int    `json:"db"`
	Address  string `json:"address"`
	Password string `json:"password"`
}

type NSQ struct {
	Addrs             []string `json:"addrs"`
	Timeout           int      `json:"timeout" default:"60"`
	MessageTopic      string   `json:"messageTopic"`
	MessageTopicGroup string   `json:"messageTopicGroup"`
	OrderExpireTopic  string   `json:"orderExpireTopic"`
}

type HttpClient struct {
	SocketServerClient string `json:"socketServerAddress"`
	Timeout            int    `json:"timeout"`
	Proxy              string `json:"proxy"`
}

type CardPackage struct {
	Id    string `json:"id"`
	Price int64  `json:"price"` // 精确到具体多少U，数据库存储  会带上 4个0
	Day   int    `json:"day"`
	Name  string `json:"name"`
}

func (c Config) GetPackageByID(id string) (CardPackage, bool) {
	for _, item := range c.CardPackages {
		if item.Id == id {
			return item, true
		}
	}
	return CardPackage{}, false
}

type SocketIO struct {
	Host     string `json:"host"`     // 域名:
	FullHost string `json:"fullHost"` // 完整域名带端口
}
