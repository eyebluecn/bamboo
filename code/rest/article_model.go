package rest

import "github.com/eyebluecn/bamboo/code/core"

/**
 * article
 */
type Article struct {
	Base
	UserUuid       string `json:"userUuid" gorm:"type:char(36)"`
	Title          string `json:"title" gorm:"type:varchar(255) not null"`
	Path           string `json:"path" gorm:"type:varchar(45) not null"`
	Tags           string `json:"tags" gorm:"type:varchar(1024) not null"`
	PosterTankUuid string `json:"posterTankUuid" gorm:"type:char(36)"`
	PosterUrl      string `json:"posterUrl" gorm:"type:varchar(512) not null"`
	Author         string `json:"author" gorm:"type:varchar(45) not null"`
	Digest         string `json:"digest" gorm:"type:varchar(512) not null"`
	IsMarkdown     bool   `json:"isMarkdown" gorm:"type:tinyint(1) not null;default:1"`
	Markdown       string `json:"markdown" gorm:"type:mediumtext"`
	Html           string `json:"html" gorm:"type:mediumtext"`
	Privacy        bool   `json:"privacy" gorm:"type:tinyint(1) not null;default:0"`
	Top            bool   `json:"top" gorm:"type:tinyint(1) not null;default:0"`
	Agree          int64  `json:"agree" gorm:"type:bigint(20) not null;default:0"`
	Words          int64  `json:"words" gorm:"type:bigint(20) not null;default:0"`
	Hit            int64  `json:"hit" gorm:"type:bigint(20) not null;default:0"`
	CommentNum     int64  `json:"commentNum" gorm:"type:bigint(20) not null;default:0"`
}

// set File's table name to be `profiles`
func (this *Article) TableName() string {
	return core.TABLE_PREFIX + "article"
}
