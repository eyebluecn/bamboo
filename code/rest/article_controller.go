package rest

import (
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/tool/result"
	"net/http"
)

type ArticleController struct {
	BaseController
	articleDao     *ArticleDao
	articleService *ArticleService
}

func (this *ArticleController) Init() {
	this.BaseController.Init()

	b := core.CONTEXT.GetBean(this.articleDao)
	if b, ok := b.(*ArticleDao); ok {
		this.articleDao = b
	}

	b = core.CONTEXT.GetBean(this.articleService)
	if b, ok := b.(*ArticleService); ok {
		this.articleService = b
	}

}

func (this *ArticleController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	routeMap["/api/article/create"] = this.Wrap(this.Create, USER_ROLE_USER)

	routeMap["/api/article/delete"] = this.Wrap(this.Delete, USER_ROLE_USER)

	return routeMap
}

// Create an article.
func (this *ArticleController) Create(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	title := request.FormValue("title")
	path := request.FormValue("path")
	author := request.FormValue("author")

	user := this.checkUser(request)
	article := &Article{
		UserUuid: user.Uuid,
		Title:    title,
		Path:     path,
		Author:   author,
	}

	article = this.articleDao.Create(article)

	return this.Success(article)
}

//delete an article.
func (this *ArticleController) Delete(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic(result.BadRequest("uuid cannot be null"))
	}

	article := this.articleDao.CheckByUuid(uuid)

	//only the author can delete the article.
	user := this.checkUser(request)
	if article.UserUuid != user.Uuid {
		panic(result.UNAUTHORIZED)
	}

	this.articleDao.Delete(article)

	return this.Success("OK")
}
