package controllers

import (
	_ "encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
)

type GetHouseInfoController struct {
	beego.Controller
}
type HouseInfo struct {
	Address     string `json:"address"`
	Area_name   string `json:"area_name"`
	Ctime       string `json:"ctime"`
	House_id    int    `json:"house_id"`
	Img_url     string `json:"img_url"`
	Order_count int    `json:"order_count"`
	Price       int    `json:"price"`
	Room_count  int    `json:"room_count"`
	Title       string `json:"title"`
	User_avatar string `json:"user_avatar"`
}
type RespHouseInfo struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   []HouseInfo `json:"data"`
}

func (this *GetHouseInfoController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}
func (this *GetHouseInfoController) GetHouseInfo() {
	resp := RegResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
	//通过session得到user_id
	resp_user_id := this.GetSession("user_id")

	fmt.Println("===>houses", resp_user_id)
	//查询house表，找到user_id所有的房屋

	o := orm.NewOrm()

	qs := o.QueryTable("house")
	var resqs orm.QuerySeter
	resqs = qs.Filter("user_id", resp_user_id)
	fmt.Printf("resqs ===>%+v", resqs)

	//返回json数据

	return
}
