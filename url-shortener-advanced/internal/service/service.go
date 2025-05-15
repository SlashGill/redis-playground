package service

import (
	"errors"
	"log"
	"math/rand"
	"time"
	"url-shortener-basic/internal/db"
	"url-shortener-basic/internal/repository"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(n int) string {
	// 創建一個新的隨機數生成器，使用當前時間的納秒作為種子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))] // 使用新的生成器 r
	}
	return string(b)
}

func Shorten(url string) (string, error) {
	shortCode := generateShortCode(6)
	err := repository.Save(shortCode, url)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func Resolve(shortCode string) (string, error) {

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
