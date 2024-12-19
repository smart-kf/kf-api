package config

import "fmt"

type Config struct {
	Debug  bool   `json:"debug"`
	Web    Web    `json:"web"`
	Log    Log    `json:"log"`
	DB     Db     `json:"db"`
	JwtKey string `json:"jwtKey"`

	BillConfig BillConfig `json:"billConfig"`
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
	Accounts []BillAccount `json:"accounts"`
}

type BillAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
