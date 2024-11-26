package model

import (
	"log"
	"magmar/util"
	"os"
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
