package rest

import "github.com/eyebluecn/bamboo/code/core"

type Preference struct {
	Base
	Name          string `json:"name" gorm:"type:varchar(45)"`
	LogoUrl       string `json:"logoUrl" gorm:"type:varchar(255)"`
	FaviconUrl    string `json:"faviconUrl" gorm:"type:varchar(255)"`
	Copyright     string `json:"copyright" gorm:"type:varchar(1024)"`
	Record        string `json:"record" gorm:"type:varchar(1024)"`
	AllowRegister bool   `json:"allowRegister" gorm:"type:tinyint(1) not null;default:0"`
	Version       string `json:"version" gorm:"-"`
}

// set File's table name to be `profiles`
func (this *Preference) TableName() string {
	return core.TABLE_PREFIX + "preference"
}
