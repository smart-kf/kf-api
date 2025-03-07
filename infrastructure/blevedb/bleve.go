package blevedb

import (
	"sync"

	"github.com/blevesearch/bleve/v2"

	"github.com/smart-fm/kf-api/config"
)

var (
	index bleve.Index
	o     sync.Once
)

func Load() bleve.Index {
	o.Do(
		func() {
			conf := config.GetConfig()
			idx, err := bleve.New(conf.Bleve.Path, bleve.NewIndexMapping())
			if err != nil {
				panic(err)
			}
			index = idx
		},
	)

	return index
}

func GetIndex() bleve.Index {
	return index
}
