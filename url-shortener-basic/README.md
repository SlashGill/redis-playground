# URL-SHORTENER-BASIC - 短網址服務

本專案是一個使用 Go 語言構建的輕量級短網址服務，旨在提供基本的 URL 縮短與重定向功能。


## 技術棧

* **程式語言:** Go (golang)
* **網路處理:** `net/http` (Go 標準庫)
* **資料庫:** PostgreSQL
* **PostgreSQL 驅動:** `github.com/lib/pq`
* **JSON 處理:** `encoding/json` (Go 標準庫)

## 部署與運行

### 前置條件

* 已安裝並成功啟動 PostgreSQL 資料庫伺服器。
* 已在 PostgreSQL 中建立好儲存 URL 映射關係的資料表。
* 已安裝 Go 語言開發環境。

### 編譯與執行

1.  **初始化專案**
    ```bash
    go mod init url-shortener-basic
    ```

2.  **安裝依賴**
    ```bash
    go mod tidy
    ```

3.  **編譯程式:**
    ```bash
    go build cmd/main.go
    ```
    這會在當前目錄下產生一個可執行檔 `main` 

4.  **啟動服務:**
    ```bash
    ./main
    ```
    執行編譯後的程式以啟動短網址服務。服務預設監聽 `http://localhost:9999`。

## API 介面說明

### 1. 建立短網址 (Shorten URL)

* **HTTP 方法:** `POST`
* **請求路徑:** `/shorten`
* **請求體 (Request Body):** `application/json`
    ```json
    {
      "url": "https://example.com/gill-page"
    }
    ```
* **成功回應 (Success Response):** HTTP 200 OK, `application/json`
    ```json
    {
      "short_url": "http://localhost:9999/abc123"
    }
    ```
    * `short_url`: 生成的短網址。

### 2. 使用短碼重定向至原始網址 (Resolve Short Code)

* **HTTP 方法:** `GET`
* **請求路徑:** `/yourShortCode` (將 `yourShortCode` 替換為實際的短代碼)
* **成功回應 (Success Response):** HTTP 302 Found (臨時重定向)
* **失敗回應 (Failure Response):** HTTP 404 Not Found
    * 當提供的短代碼在資料庫中找不到對應的原始網址時返回。