package src

type KfPerson struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	IdCard string `db:"idcard"`
}

type QueryResult struct {
	//
	Value []KfPerson
	// time for adding into cache
	CacheTime int64
	// Count was queried
	Count int
}

func (qr *QueryResult) GetCacheTime() int64 {
	return qr.CacheTime
}
