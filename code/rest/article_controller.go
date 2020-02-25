package rest

import (
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/tool/builder"
	"github.com/eyebluecn/bamboo/code/tool/result"
	"net/http"
	"strconv"
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

	routeMap["/api/article/edit"] = this.Wrap(this.Edit, USER_ROLE_USER)

	routeMap["/api/article/delete"] = this.Wrap(this.Delete, USER_ROLE_USER)

	routeMap["/api/article/list"] = this.Wrap(this.List, USER_ROLE_USER)

	routeMap["/api/article/detail"] = this.Wrap(this.Detail, USER_ROLE_USER)

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

// Edit an article.
func (this *ArticleController) Edit(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	title := request.FormValue("title")
	path := request.FormValue("path")
	author := request.FormValue("author")

	article := this.articleDao.CheckByUuid(uuid)

	article.Title = title
	article.Path = path
	article.Author = author

	article = this.articleDao.Save(article)

	return this.Success(article)
}

//delete an article.
func (this *ArticleController) Delete(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic(result.BadRequest("uuid cannot be null"))
	}

	article := this.articleDao.CheckByUuid(uuid)

	user := this.checkUser(request)
	if article.UserUuid != user.Uuid {
		//TODO: only the author can delete the article.
		//panic(result.UNAUTHORIZED)
	}

	this.articleDao.Delete(article)

	return this.Success("OK")
}

func (this *ArticleController) List(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	//use 0 base.
	pageStr := request.FormValue("page")
	pageSizeStr := request.FormValue("pageSize")
	orderCreateTime := request.FormValue("orderCreateTime")
	orderUpdateTime := request.FormValue("orderUpdateTime")
	orderSort := request.FormValue("orderSort")

	userUuid := request.FormValue("userUuid")
	title := request.FormValue("title")
	path := request.FormValue("path")
	author := request.FormValue("author")

	var page int
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	pageSize := 20
	if pageSizeStr != "" {
		tmp, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			pageSize = tmp
		}
	}

	sortArray := []builder.OrderPair{
		{
			Key:   "create_time",
			Value: orderCreateTime,
		},
		{
			Key:   "update_time",
			Value: orderUpdateTime,
		},
		{
			Key:   "sort",
			Value: orderSort,
		},
	}

	pager := this.articleDao.Page(page, pageSize, userUuid, title, path, author, sortArray)

	return this.Success(pager)
}

func (this *ArticleController) Detail(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	uuid := request.FormValue("uuid")

	article := this.articleDao.CheckByUuid(uuid)

	return this.Success(article)

}
