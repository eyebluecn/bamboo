package rest

import (
	"github.com/eyebluecn/bamboo/code/core"
)

//@Service
type ArticleService struct {
	BaseBean
	articleDao *ArticleDao
	userDao    *UserDao
}

func (this *ArticleService) Init() {
	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.articleDao)
	if b, ok := b.(*ArticleDao); ok {
		this.articleDao = b
	}

	b = core.CONTEXT.GetBean(this.userDao)
	if b, ok := b.(*UserDao); ok {
		this.userDao = b
	}

}

func (this *ArticleService) Detail(uuid string) *Article {

	article := this.articleDao.CheckByUuid(uuid)

	return article
}

func (this *ArticleService) Bootstrap() {

}
