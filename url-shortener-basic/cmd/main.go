package main

import (
	"log"
	"net/http"
	"url-shortener-basic/internal/db"
	"url-shortener-basic/internal/handler"
)

func main() {
	db.Init() // 初始化資料庫連線

	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/", handler.ResolveURL)

	log.Println("Server is running on :9999")
	log.Fatal(http.ListenAndServe(":9999", nil))

}
