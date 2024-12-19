package config

type Config struct {
	Debug bool `json:"debug"`
	Web   Web  `json:"web"`
	Log   Log  `json:"log"`
	DB    Db   `json:"db"`
}

type Web struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
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
