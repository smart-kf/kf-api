package utils

import "encoding/json"

// LeveldbKey 生成leveldb key
func LeveldbKey(key string, keys ...string) []byte {
	var res []byte
	res = append(res, []byte(key)...)
	for _, k := range keys {
		res = append(res, []byte(k)...)
	}
	return res
}

func MustMarshalObject(obj interface{}) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func MustUnMarshalObject(data []byte, obj interface{}) {
	json.Unmarshal(data, obj)
}
