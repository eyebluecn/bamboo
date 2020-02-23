package i18n

import (
	"golang.org/x/text/language"
	"net/http"
)

const (
	LANG_KEY = "_lang"
)

var matcher = language.NewMatcher([]language.Tag{
	// The first language is used as fallback.
	language.English,
	language.Chinese,
})

type Item struct {
	English string
	Chinese string
}

var (
	UsernameOrPasswordCannotNull = &Item{English: `username or password cannot be null`, Chinese: `用户名或密码不能为空`}
	UsernameOrPasswordError      = &Item{English: `username or password error`, Chinese: `用户名或密码错误`}
	UsernameExist                = &Item{English: `username "%s" exists`, Chinese: `用户名"%s"已存在`}
	UsernameNotExist             = &Item{English: `username "%s" not exists`, Chinese: `用户名"%s"不存在`}
	UsernameIsNotAdmin           = &Item{English: `username "%s" is not admin user`, Chinese: `用户名"%s"不是管理员账号`}
	UsernameError                = &Item{English: `username can only be lowercase letters, numbers or _`, Chinese: `用户名必填，且只能包含小写字母，数字和'_'`}
	UserRegisterNotAllowd        = &Item{English: `admin has banned register`, Chinese: `管理员已禁用自主注册`}
	UserPasswordLengthError      = &Item{English: `password at least 6 chars`, Chinese: `密码长度至少为6位`}
	UserOldPasswordError         = &Item{English: `old password error`, Chinese: `旧密码不正确`}
	UserDisabled                 = &Item{English: `user has been disabled`, Chinese: `用户已经被禁用了`}
)

func (this *Item) Message(request *http.Request) string {

	if request == nil {
		return this.English
	}

	lang, _ := request.Cookie(LANG_KEY)
	formLangStr := request.FormValue(LANG_KEY)
	acceptLangStr := request.Header.Get("Accept-Language")
	var cookieLangStr string
	if lang != nil {
		cookieLangStr = lang.Value
	}
	tag, _ := language.MatchStrings(matcher, cookieLangStr, formLangStr, acceptLangStr)

	tagBase, _ := tag.Base()
	chineseBase, _ := language.Chinese.Base()

	if tagBase == chineseBase {
		return this.Chinese
	} else {
		return this.English
	}

}
