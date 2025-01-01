package model

import (
	"log"
	"magmar/util"
	"os"
	"strings"
	"time"

	"gorm.io/plugin/soft_delete"
)

var zlog *util.Logger

func init() {
	_, err := util.NewLogger()
	if err != nil {
		log.Fatalf("InitLog module[model] err[%s]", err.Error())
		os.Exit(1)
	}
}

// BaseModel ...
type BaseModel struct {
	ID        uint64                `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time             `json:"createdAt" gorm:"<-:create;autoCreateTime;not null"`
	UpdatedAt time.Time             `json:"updatedAt" gorm:"autoUpdateTime;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt" gorm:"default:0"`
}

func toStr(strs []string, idx int) string {
	if idx == 999999 {
		return strs[len(strs)-1]
	}
	if len(strs) <= idx {
		return strs[0]
	}
	return strs[idx]
}

func unmarshalJSON(data []byte, mapData map[string]int, defaultVal int) (val int) {
	strData := strings.Trim(string(data), "\"")
	return unmarshal(strData, mapData, defaultVal)
}

func unmarshal(data string, mapData map[string]int, defaultVal int) (val int) {
	if data == "" {
		return defaultVal
	}

	return mapData[strings.ToLower(data)]
}
