package src

import "time"

type TimeData interface {
	// 获得加入缓存的时间纳秒
	GetCacheTime() int64
}

// 删除最早的缓存
func UpdateCache(cacheMap *map[string]TimeData) (delKey string) {
	// predefine a earliest time
	earliestTime := time.Now().UnixNano()
	for key, value := range *cacheMap {
		if value.GetCacheTime() < earliestTime {
			earliestTime = value.GetCacheTime()
			delKey = key
		}
	}
	delete(*cacheMap, delKey)
	return delKey
}
