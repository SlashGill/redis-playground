# URL-SHORTENER-ADVANCED - 短網址服務

本專案利用 Redis 做為 Cache 機制，調整來自 URL-SHORTENER-BASIC 



## 雙方主要差異

1. 在 **/internal/service/service.go**  內，要找到短網址應對的原始網址，會先從 Redis 查找。找不到才到 Postgres 找，並且在 Redis 加上一筆資料做為快取

```
unc Resolve(shortCode string) (string, error) {

	// 先查詢 Redis 快取
	val, err := db.Rdb.Get(db.Ctx, shortCode).Result()
	if err == nil {
		log.Printf("Cache hit for %s: %s", shortCode, val)
		return val, nil // 快取命中
	}

	// 如果 Redis 中沒有，則查詢資料庫
	url, err := repository.Find(shortCode)
	if err != nil {
		log.Printf("get from db for %s: %s", shortCode, err)
		return "", errors.New("URL not found")
	}

	// 寫入 Redis
	db.Rdb.Set(db.Ctx, shortCode, url, time.Hour*24)
	time.Sleep(time.Second * 5) // 模擬延遲
	return url, nil
}
```

2. 多了 **/internal/db/redis.go** 對 Redis 進行連線設定
