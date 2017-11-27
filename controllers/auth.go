package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome/models"
)

type Baseinfo struct {
	User_id   int    `json:"user_id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Mobile    string `json:"mobile"`
	Real_name string `json:"real_name"`
	ID_card   string `json:"id_card"`
	Url       string `json:"avatar_url"`
}

type AuthcheckResp struct {
	Errno  string   `json:"errno"`
	Errmsg string   `json:"errmsg"`
	Data   Baseinfo `json:"data"`
}

type AuthInfo struct {
	Real_name string `json:"real_name"`
	ID_card   string `json:"id_card"`
}

type AuthController struct {
	beego.Controller
}

func (this *AuthController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *AuthController) AuthCheck() {
	authcheckresp := AuthcheckResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&authcheckresp)

	user_id := this.GetSession("user_id")

	var userinfo models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("id", user_id).One(&userinfo)

	if err != nil {
		//查询出错
		beego.Debug("AuthCheck error: ", err)
		authcheckresp.Errno = models.RECODE_NODATA
		authcheckresp.Errmsg = models.RecodeText(authcheckresp.Errno)
		return
	}

	if userinfo.Real_name == "" || userinfo.Id_card == "" {
		//表示未认证
		beego.Debug("用户未认证")
		authcheckresp.Errno = models.RECODE_ROLEERR
		authcheckresp.Errmsg = models.RecodeText(authcheckresp.Errno)
		return
	}

	var baseinfo Baseinfo
	baseinfo.User_id = userinfo.Id
	baseinfo.Name = userinfo.Name
	baseinfo.Password = userinfo.Password_hash
	baseinfo.Mobile = userinfo.Mobile
	baseinfo.Real_name = userinfo.Real_name
	baseinfo.ID_card = userinfo.Id_card
	baseinfo.Url = userinfo.Avatar_url

	authcheckresp.Data = baseinfo
	return
}

func (this *AuthController) UpdateAuthinfo() {
	var auth_info AuthInfo
	json.Unmarshal(this.Ctx.Input.RequestBody, &auth_info)

	authinforesp := AuthcheckResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&authinforesp)

	user_id := this.GetSession("user_id")
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("id", user_id).Update(orm.Params{"real_name": auth_info.Real_name, "id_card": auth_info.ID_card})

	if err != nil {
		authinforesp.Errno = models.RECODE_DATAERR
		authinforesp.Errmsg = models.RecodeText(authinforesp.Errno)
		return
	}

	return
}
