package ipinfo

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	o        sync.Once
	searcher *xdb.Searcher
)

type Location struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Net      string
}

func GetLocation(dbPath string, ip string) Location {
	o.Do(
		func() {
			// 1、从 dbPath 加载 VectorIndex 缓存，把下述 vIndex 变量全局到内存里面。
			vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
			if err != nil {
				fmt.Printf("failed to load vector index from `%s`: %s\n", dbPath, err)
				return
			}

			// 2、用全局的 vIndex 创建带 VectorIndex 缓存的查询对象。
			searcher, err = xdb.NewWithVectorIndex(dbPath, vIndex)
			if err != nil {
				fmt.Printf("failed to create searcher with vector index: %s\n", err)
				return
			}
		},
	)

	if searcher == nil {
		return Location{}
	}

	// 3、查询 IP 地址的信息。
	res, err := searcher.SearchByStr(ip)
	if err != nil {
		return Location{}
	}
	parts := strings.Split(res, "|")

	loc := Location{}

	// 中国|0|四川省|成都市|电信
	if len(parts) > 0 {
		loc.Country = parts[0]
	}
	if len(parts) > 2 {
		loc.Province = parts[2]
	}
	if len(parts) > 3 {
		loc.City = parts[3]
	}
	if len(parts) > 4 {
		loc.Net = parts[4]
	}

	return Location{}
}
