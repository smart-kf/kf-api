package config

import "fmt"

type Config struct {
	Debug         bool          `json:"debug"`
	Web           Web           `json:"web"`
	Log           Log           `json:"log"`
	DB            Db            `json:"db"`
	LevelDBConfig LevelDBConfig `json:"levelDB"`
	JwtKey        string        `json:"jwtKey"`
	BillConfig    BillConfig    `json:"billConfig"`
	RedisConfig   RedisConfig   `json:"redis"`
	Kafka         Kafka         `json:"kafka"`
	HttpClient    HttpClient    `json:"httpClient"`
}

type LevelDBConfig struct {
	Path string `json:"path"`
}

type Web struct {
	Addr string `json:"addr" default:"127.0.0.1"`
	Port int    `json:"port" default:"8081"`
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

type Kafka struct {
	Addrs          []string `json:"addrs"`
	ImMessageTopic string   `json:"imMessageTopic"`
	ImMessageGroup string   `json:"imMessageGroup"`
}

type HttpClient struct {
	LogicAddress string `json:"logicAddress"` // logic 服务 http 地址
	Timeout      int    `json:"timeout"`
	Proxy        string `json:"proxy"`
}
