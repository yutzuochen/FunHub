# README FunHub

<!-- ![專案封面圖](https://fakeimg.pl/500/)

> 此專案是一份 README 的撰寫範本，主要是方便所有人可以快速撰寫 README，讓大家可以更有方向的去寫出 README。

- [線上觀看連結](https://israynotarray.com/) -->

## 功能

測試帳號密碼 

```bash
帳號： m2
密碼： 2222
```

...
<!-- 
## 畫面

> 可提供 1~3 張圖片，讓觀看者透過 README 了解整體畫面

![範例圖片 1](https://fakeimg.pl/500/)
![範例圖片 2](https://fakeimg.pl/500/)
![範例圖片 3](https://fakeimg.pl/500/) -->

## 安裝

Golang 版本建議為：`go 1.19.3` 以上

### 取得專案

```bash
git clone https://github.com/yutzuochen/FunHub.git
```

### 移動到專案內

```bash
cd FunHub
```

### 安裝套件

```bash
go 1.19
```

<!-- ### 環境變數設定

請在終端機輸入 `cp .env.example .env` 來複製 .env.example 檔案，並依據 `.env` 內容調整相關欄位。 -->

### 運行專案

```bash
go run FunHub
```

### 開啟專案

在瀏覽器網址列輸入以下即可看到畫面

```bash
http://localhost:8080/login
http://localhost:8080/play?ante=1
```

## 操作流程

1.先使用URL "http://localhost:8080/login" 進行登入，帳號密碼驗證成功會拿到一組JWT

2.以cookie 攜帶該 JWT，使用 URL "http://localhost:8080/play?ante=1" 進行遊玩



## 資料夾說明

FunHub
 ├── constant/          # 常數定義
 ├── gamecore/          # 遊戲邏輯運算
 ├── db/                # 資料庫相關功能
 ├── doc/               # FunHub 說明文件
 ├── jwtTools/          # token 相關功能
 └── main.go            # entry point
...


-- Demo Video-- 