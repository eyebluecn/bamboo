package rest

import (
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/tool/builder"
	"github.com/eyebluecn/bamboo/code/tool/result"
	"github.com/jinzhu/gorm"

	"github.com/nu7hatch/gouuid"
	"time"
)

type ArticleDao struct {
	BaseDao
}

//find by uuid. if not found return nil.
func (this *ArticleDao) FindByUuid(uuid string) *Article {
	var entity = &Article{}
	db := core.CONTEXT.GetDB().Where("uuid = ?", uuid).First(entity)
	if db.Error != nil {
		if db.Error.Error() == result.DB_ERROR_NOT_FOUND {
			return nil
		} else {
			panic(db.Error)
		}
	}
	return entity
}

//find by uuid. if not found panic NotFound error
func (this *ArticleDao) CheckByUuid(uuid string) *Article {
	entity := this.FindByUuid(uuid)
	if entity == nil {
		panic(result.NotFound("not found record with uuid = %s", uuid))
	}
	return entity
}

func (this *ArticleDao) Page(page int, pageSize int, userUuid string, title string, path string, author string, sortArray []builder.OrderPair) *Pager {

	var wp = &builder.WherePair{}

	if userUuid != "" {
		wp = wp.And(&builder.WherePair{Query: "user_uuid = ?", Args: []interface{}{userUuid}})
	}

	if title != "" {
		wp = wp.And(&builder.WherePair{Query: "title LIKE ?", Args: []interface{}{"%" + title + "%"}})
	}

	if path != "" {
		wp = wp.And(&builder.WherePair{Query: "path LIKE ?", Args: []interface{}{"%" + path + "%"}})
	}

	if author != "" {
		wp = wp.And(&builder.WherePair{Query: "author LIKE ?", Args: []interface{}{"%" + author + "%"}})
	}

	var conditionDB *gorm.DB
	conditionDB = core.CONTEXT.GetDB().Model(&Article{}).Where(wp.Query, wp.Args...)

	count := 0
	db := conditionDB.Count(&count)
	this.PanicError(db.Error)

	var articles []*Article
	db = conditionDB.Order(this.GetSortString(sortArray)).Offset(page * pageSize).Limit(pageSize).Find(&articles)
	this.PanicError(db.Error)
	pager := NewPager(page, pageSize, count, articles)

	return pager
}

func (this *ArticleDao) Create(article *Article) *Article {

	timeUUID, _ := uuid.NewV4()
	article.Uuid = string(timeUUID.String())
	article.CreateTime = time.Now()
	article.UpdateTime = time.Now()
	article.Sort = time.Now().UnixNano() / 1e6
	db := core.CONTEXT.GetDB().Create(article)
	this.PanicError(db.Error)

	return article
}

func (this *ArticleDao) Save(article *Article) *Article {

	article.UpdateTime = time.Now()
	db := core.CONTEXT.GetDB().Save(article)
	this.PanicError(db.Error)

	return article
}

func (this *ArticleDao) Delete(article *Article) {

	db := core.CONTEXT.GetDB().Delete(&article)
	this.PanicError(db.Error)
}

func (this *ArticleDao) CountBetweenTime(startTime time.Time, endTime time.Time) int64 {
	var count int64
	db := core.CONTEXT.GetDB().Model(&Article{}).Where("create_time >= ? AND create_time <= ?", startTime, endTime).Count(&count)
	this.PanicError(db.Error)
	return count
}

//System cleanup.
func (this *ArticleDao) Cleanup() {
	this.logger.Info("[ArticleDao][DownloadTokenDao] clean up. Delete all Article")
	db := core.CONTEXT.GetDB().Where("uuid is not null").Delete(Article{})
	this.PanicError(db.Error)
}
