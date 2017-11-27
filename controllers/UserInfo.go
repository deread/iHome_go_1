package controllers

import (
	//	"encoding/json"
	//	"fmt"
	"github.com/astaxie/beego"
	//	//"github.com/astaxie/beego/cache"
	//	//_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//	//	"time"
	//	"github.com/astaxie/beego/config
	//"path"
)

type UserInfoResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}
type UserInfoController struct {
	beego.Controller
}

func (this *UserInfoController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}
func (this *UserInfoController) UserInfoGet() {
	beego.Debug("get /api/v1.0/user....UserInfoGet")

	resp := UserInfoResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)
	user_id := this.GetSession("user_id")
	if user_id == nil {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return

	}
	//查询数据库
	o := orm.NewOrm()
	var user models.User
	if err := o.QueryTable("user").Filter("id", user_id).One(&user); err == orm.ErrNoRows {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	resp.Data = user
	return
}
