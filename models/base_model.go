package models

import (
	"time"

	. "github.com/onestack/cron-room/helper"
	"github.com/labstack/gommon/color"
	"github.com/jinzhu/gorm"
	"fmt"
)

var (
	DB            *gorm.DB
	firstRegModel bool = true
	model         *Model
)

type BaseModel struct {
	CreatedTimeStr string     `json:"createdTime" gorm:"-"`
	CreatedTime    *time.Time `json:"-" gorm:"column:created_time; type:datetime; not null; default:current_timestamp"`
	UpdatedTime    *time.Time `json:"-" gorm:"column:updated_time; type:datetime"`
	DeletedTime    *time.Time `json:"-" gorm:"column:deleted_time; type:datetime"`
	IsDeleted      int64      `json:"-" gorm:"column:is_deleted; type:tinyint(1); not null; default:0"`
}

func init() {
	initModel()
	DB = model.DB
	DB.LogMode(GetEnv().DebugOn)
}

func (this *BaseModel) AfterFind() {
	this.CreatedTimeStr = this.CreatedTime.Format(DateShortLayout)
}

func initModel() {
	// Initial the DB with "github.com/jinzhu/gorm"
	InitEnv()
	env := GetEnv()
	dialect := env.DBDialect
	host := env.DBHost
	port := env.DBPort
	dbname := env.DBName
	user := env.DBUser
	password := env.DBPassword
	model = NewModel()
	if !model.IsOpenDB() {
		err := model.OpenWithConfig(dialect, host, port, dbname, user, password)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := NewModelWithConfig(dialect, host, port, dbname, user, password)
		if err != nil {
			panic(err)
		}
	}
}

// RegisterModels registers models
func RegisterModels(models ...interface{}) {
	if firstRegModel {
		fmt.Println(fmt.Sprintf("%s%s%s",
			color.White(Line2),
			color.Bold(color.Green("Model information")),
			color.White(Line2)))
	}
	model.AddModels(models...)
	firstRegModel = false
}

// Migrate runs migrations on the global
func MigrateAll() {
	model.AutoMigrateAll()
}
