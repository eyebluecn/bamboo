package rest

import (
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/tool/result"
	"github.com/eyebluecn/bamboo/code/tool/util"
	"net/http"
)

type PreferenceController struct {
	BaseController
	preferenceDao     *PreferenceDao
	preferenceService *PreferenceService
}

func (this *PreferenceController) Init() {
	this.BaseController.Init()

	b := core.CONTEXT.GetBean(this.preferenceDao)
	if b, ok := b.(*PreferenceDao); ok {
		this.preferenceDao = b
	}

	b = core.CONTEXT.GetBean(this.preferenceService)
	if b, ok := b.(*PreferenceService); ok {
		this.preferenceService = b
	}

}

func (this *PreferenceController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	routeMap["/api/preference/ping"] = this.Wrap(this.Ping, USER_ROLE_GUEST)
	routeMap["/api/preference/fetch"] = this.Wrap(this.Fetch, USER_ROLE_GUEST)
	routeMap["/api/preference/edit"] = this.Wrap(this.Edit, USER_ROLE_ADMINISTRATOR)
	routeMap["/api/preference/system/cleanup"] = this.Wrap(this.SystemCleanup, USER_ROLE_ADMINISTRATOR)

	return routeMap
}

//ping the application. Return current version.
func (this *PreferenceController) Ping(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	return this.Success(core.VERSION)

}

func (this *PreferenceController) Fetch(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	preference := this.preferenceService.Fetch()

	return this.Success(preference)
}

func (this *PreferenceController) Edit(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	name := request.FormValue("name")

	logoUrl := request.FormValue("logoUrl")
	faviconUrl := request.FormValue("faviconUrl")
	copyright := request.FormValue("copyright")
	record := request.FormValue("record")
	allowRegisterStr := request.FormValue("allowRegister")

	if name == "" {
		panic(result.BadRequest("name cannot be null"))
	}

	var allowRegister = false
	if allowRegisterStr == TRUE {
		allowRegister = true
	}

	preference := this.preferenceDao.Fetch()
	preference.Name = name
	preference.LogoUrl = logoUrl
	preference.FaviconUrl = faviconUrl
	preference.Copyright = copyright
	preference.Record = record
	preference.AllowRegister = allowRegister

	preference = this.preferenceDao.Save(preference)

	//reset the preference cache
	this.preferenceService.Reset()

	return this.Success(preference)
}

//cleanup system data.
func (this *PreferenceController) SystemCleanup(writer http.ResponseWriter, request *http.Request) *result.WebResult {

	user := this.checkUser(request)
	password := request.FormValue("password")

	if !util.MatchBcrypt(password, user.Password) {
		panic(result.BadRequest("password error"))
	}

	//this will trigger every bean to cleanup.
	core.CONTEXT.Cleanup()

	return this.Success("OK")
}
