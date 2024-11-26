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
	CtxTimeOut = time.Second * 10
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

// ContestKey ...
const (
	LoginKey = "login"
)
