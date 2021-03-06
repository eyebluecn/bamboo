package support

import (
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/rest"
	"github.com/eyebluecn/bamboo/code/tool/cache"
	"github.com/jinzhu/gorm"
	"net/http"
	"reflect"
)

type BambooContext struct {
	//db connection
	db *gorm.DB
	//session cache
	SessionCache *cache.Table
	//bean map.
	BeanMap map[string]core.Bean
	//controller map
	ControllerMap map[string]core.Controller
	//router
	Router *BambooRouter
}

func (this *BambooContext) Init() {

	//create session cache
	this.SessionCache = cache.NewTable()

	//init map
	this.BeanMap = make(map[string]core.Bean)
	this.ControllerMap = make(map[string]core.Controller)

	//register beans. This method will put Controllers to ControllerMap.
	this.registerBeans()

	//init every bean.
	this.initBeans()

	//create and init router.
	this.Router = NewRouter()

	//if the application is installed. Bean's Bootstrap method will be invoked.
	this.InstallOk()

}

func (this *BambooContext) GetDB() *gorm.DB {
	return this.db
}

func (this *BambooContext) GetSessionCache() *cache.Table {
	return this.SessionCache
}

func (this *BambooContext) GetControllerMap() map[string]core.Controller {
	return this.ControllerMap
}

func (this *BambooContext) Cleanup() {
	for _, bean := range this.BeanMap {
		bean.Cleanup()
	}
}

//can serve as http server.
func (this *BambooContext) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	this.Router.ServeHTTP(writer, request)
}

func (this *BambooContext) OpenDb() {

	var err error = nil
	this.db, err = gorm.Open("mysql", core.CONFIG.MysqlUrl())

	if err != nil {
		core.LOGGER.Panic("failed to connect mysql database")
	}

	//whether open the db log. (only true when debug)
	this.db.LogMode(false)
}

func (this *BambooContext) CloseDb() {

	if this.db != nil {
		err := this.db.Close()
		if err != nil {
			core.LOGGER.Error("occur error when closing db %s", err.Error())
		}
	}
}

func (this *BambooContext) registerBean(bean core.Bean) {

	typeOf := reflect.TypeOf(bean)
	typeName := typeOf.String()

	if element, ok := bean.(core.Bean); ok {

		if _, ok := this.BeanMap[typeName]; ok {
			core.LOGGER.Error("%s has been registerd, skip", typeName)
		} else {
			this.BeanMap[typeName] = element

			//if is controller type, put into ControllerMap
			if controller, ok1 := bean.(core.Controller); ok1 {
				this.ControllerMap[typeName] = controller
			}

		}

	} else {
		core.LOGGER.Panic("%s is not the Bean type", typeName)
	}

}

func (this *BambooContext) registerBeans() {

	//article
	this.registerBean(new(rest.ArticleController))
	this.registerBean(new(rest.ArticleDao))
	this.registerBean(new(rest.ArticleService))

	//install
	this.registerBean(new(rest.InstallController))

	//preference
	this.registerBean(new(rest.PreferenceController))
	this.registerBean(new(rest.PreferenceDao))
	this.registerBean(new(rest.PreferenceService))

	//footprint
	this.registerBean(new(rest.FootprintController))
	this.registerBean(new(rest.FootprintDao))
	this.registerBean(new(rest.FootprintService))

	//session
	this.registerBean(new(rest.SessionDao))
	this.registerBean(new(rest.SessionService))

	//user
	this.registerBean(new(rest.UserController))
	this.registerBean(new(rest.UserDao))
	this.registerBean(new(rest.UserService))

}

func (this *BambooContext) GetBean(bean core.Bean) core.Bean {

	typeOf := reflect.TypeOf(bean)
	typeName := typeOf.String()

	if val, ok := this.BeanMap[typeName]; ok {
		return val
	} else {
		core.LOGGER.Panic("%s not registered", typeName)
		return nil
	}
}

func (this *BambooContext) initBeans() {

	for _, bean := range this.BeanMap {
		bean.Init()
	}
}

//if application installed. invoke this method.
func (this *BambooContext) InstallOk() {

	if core.CONFIG.Installed() {
		this.OpenDb()

		for _, bean := range this.BeanMap {
			bean.Bootstrap()
		}
	}

}

func (this *BambooContext) Destroy() {
	this.CloseDb()
}
