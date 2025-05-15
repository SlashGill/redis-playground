# redis-playground 
這個專案用短網址服務搭配 Redis 快取機制做為展示。重點不在於程式開發細節，也不會深入介紹 Redis 的各項功能。這個年代想要細節就去問 LLM，絕對比我能提供的好上 100 倍。我只希望你在進行 系統設計（System Design）規劃的初期，能意識到***系統是否需要導入快取機制來提升效能與穩定性***，本專案就算功德圓滿了。

## 火種
你是個計算機，只做除法運算，最後答案 四捨五入。每天都有一堆人來問你一些問題，比如 10 / 4 或是 8 / 5 之類的。你發現每天都有超多人問你 100 / 17 這個問題。這個問題讓你算到火大，而你每次都是告訴他們相同一句話「答案是該死的 6」 

有一次，你發現如果先把 100 / 17 = 6 寫在一張紙上，當有人問起這個問題，就先來這張紙看看正確答案，不但不用動腦計算，還可以很快的回覆他們「答案是該死的 6」。老天，這方法真是太棒了!

這招就是 ***快取***

## 前置技能要求
1. golang 開發
2. RDB 操作
3. Docker 操作 (含 docker compose)

## 安裝套件
如果和我一樣懶得一個一個裝，就用 docker compose 吧
```
docker compose up -d
```

如果你是正港的男子漢，不妨自己參考官網手動安裝!

1. [postgres](https://hub.docker.com/_/postgres) : 資料庫
2. [redis](https://hub.docker.com/_/redis) : 快取資料庫
3. [redisinsight](https://hub.docker.com/r/redis/redisinsight) : 查看 redis 儲存的內容
4. [cloudbeaver](https://hub.docker.com/r/dbeaver/cloudbeaver) : 操作 postgres DB
5. [cAdvisor](https://hub.docker.com/r/google/cadvisor) : 查看 container 資源消耗
6. [Postman](https://www.postman.com/) 沒有 Docker Image : 測試 Web API 

## 環境設定
有關預設 ip / port / 帳號 / 密碼, 都可以參考 docker-compose.yaml

### 透過 cloudbeaver 建立 postgresql Table
1. 進入 cloudbeaver 
2. 設定 postgresql 連線
3. 建立 Table 的 SQL 語法
```
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
```
### redisinsight
1. 進入 redisinsight
2. 建立 Redis 連線方式


## 資料夾說明
1. **url-shortener-basic**
主要功能 : 從 DB 取得短網址對應的原始網址

2. **url-shortener-advanced**
主要功能 : 從 Redis 取得短網址對應的原始網址，取不到改由 DB 再取一次

# 學習流程
## 沒有 Redis Cache 的情境

1. 到 url-shortener-basic 的 bin 底下，
```
./compile.sh
./main &
```
2. 用 postman 或是 cmd 指令
```
curl --location 'http://localhost:9999/shorten' \
--header 'Content-Type: application/json' \
--data '{
  "url": "http://127.0.0.1:8080"
}'
```

3. 得到類似這樣的短網址
```
{"short_url":"http://localhost:9999/nPsaI7"}
```

4. 將這個短網址填在 /cmd/test/concurrent.go 裡
```
const (
	apiURL      = "http://localhost:9999/nPsaI7"
	numRequests = 1000
)
```
5. 回到 bin 底下
```
./compile.sh
./concurrent
```
6. 應該可以看到類似 (總共花了 6s)
```
Request 879 successful. Status Code: 200 (Duration: 6.439564634s)
Request 64 successful. Status Code: 200 (Duration: 6.4564304s)
All requests completed. total duration: 6.462918696s(base) 
```
7.進入 http://127.0.0.1:8080/docker/ 找到 postgres container 進入
應該會看到 CPU 和 memory 使用量被往上拉的很高


## 有 Redis Cache 的情境

1. 到 url-shortener-advanced , 將剛才得到的短網址填在 /cmd/test/concurrent.go 裡
```
const (
	apiURL      = "http://localhost:9999/nPsaI7"
	numRequests = 1000
)
```
2. 回到 bin 底下，
```
./compile.sh
./main &
./concurrent
```

3. 應該可以看到類似 (總共花了 6s)
```
Request 879 successful. Status Code: 200 (Duration: 6.439564634s)
Request 64 successful. Status Code: 200 (Duration: 6.4564304s)
All requests completed. total duration: 6.462918696s(base) 
```

4. 進入 http://127.0.0.1:8080/docker/ 找到 postgres container 進入
會看到 CPU 和 memory 被往上拉的幅度非常小! 

## Cache  用了之後要注意什麼
引用一個新系統進入環境後，不會只有好處沒有壞處。衡量利害關係是一定要做的，下面幾個要注意的事項，這可能是會造成你額外負擔! 當然所有問題的解法都很多種，就不一一列舉了

**資料一致性問題**

快取資料可能會與資料庫資料不同步，尤其是在資料頻繁變動的情況下。需要設計合適的快取失效策略（如 TTL、主動更新快取）來避免讀取過期或錯誤資料。

**快取穿透（Cache Penetration）**

當查詢的資料在 Redis 和 Postgres 資料庫都不存在時，等於一個查詢造成 2 個系統壓力。可以考慮對不存在的資料也快取一個空結果，這樣無效查詢在 Redis 就被處理掉了。

**快取雪崩（Cache Avalanche）**

Redis 上大量快取同時過期，此時剛好有大量請求，就會改向後方 Postgres 資料庫查詢，可能引發系統崩潰。最簡單的方式是設定每筆資料有不同的過期時間。

**快取擊穿（Cache Breakdown）**

某個熱點資料過期時，又剛好有大量請求這筆資料，此時會加大資料庫的運作壓力。可使用互斥鎖，比如 Redis 快取沒拿到先等5秒後再去 Redis 取第二次 (也許上一筆已經寫上新的快取)。第二次還是沒取到，才去 Postgres 資料庫取資料

**容量與記憶體管理**

Redis 是基於記憶體的快取，容量有限。需合理設定快取大小及淘汰策略，避免因記憶體不足導致 Redis 效能下降。
