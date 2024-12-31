package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Const ...
const (
	DBCharsetOption string = "DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci"
)

// CustomStr ...
type CustomStr string

// CustomStrs ...
const (
	TRID CustomStr = "trid"
)

// CustomTimes ...
const (
	CtxTimeOut = time.Second * 60
)

// TokenTypes ...
const (
	TokenTypeBaerer string = "Baerer"
)

// TokenAudiences ...
const (
	TokenAudienceAccount string = "account"
)

// GetTRID ...
func GetTRID() string {
	t := time.Now()
	randInt := strconv.Itoa(rand.Intn(8999) + 1000)
	trid := strings.Replace(t.Format("20060102150405.00"), ".", "", -1) + randInt

	return trid
}

// ContextKey ...
const (
	LoginKey = "login"
)

// EnvKey Consts
const (
	OpenAPIKey        = "openAI.apiKey"
	UpbitAccessKey    = "upbitAPI.accessKey"
	UpbitSecretKey    = "upbitAPI.secretKey"
	UpbitFeePercent   = "upbitAPI.feePercent"
	UpbitFeeScale     = "upbitAPI.feeScale"
	UpbitMinBuyAmount = "upbitAPI.minBuyAmount"
	NewsAPIKey        = "newsAPI.apiKey"
)

// Log Const
const (
	LogCtrl = "[New Request]"
	LogSvc  = "[New Service request]"
	LogRepo = "[New Request]"
)

// APIURL ...
type APIURL string

// APIURLConst ...
const (
	APIURLUpbit    APIURL = "https://api.upbit.com"
	AlternativeURL APIURL = "https://api.alternative.me"
	NewsAPIURL     APIURL = "https://newsapi.org"
)
