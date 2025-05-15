package repository

import (
	"url-shortener-basic/internal/db"
)

func Save(shortCode, url string) error {
	_, err := db.DB.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, url)
	return err
}

func Find(shortCode string) (string, error) {
	var url string
	err := db.DB.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&url)
	return url, err
}
