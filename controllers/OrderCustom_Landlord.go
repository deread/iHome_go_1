//beego.Router("api/v1.0/user/orders", &controllers.GetOrderController{}, "get:GetOrder")
package controllers

import (
	//	_ "encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//	_ "time"
)

type GetOrderController struct {
	beego.Controller
}

type OrderInfo struct {
	Amount     int    `json:"amount"`
	Comment    string `json:"comment"`
	Ctime      string `json:"ctime"`
	Days       int    `json:"days"`
	End_date   string `json:"end_date"`
	Img_url    string `json:"img_url"`
	Order_id   int    `json:"order_id"`
	Start_date string `json:"start_date"`
	Status     string `json:"status"`
	Title      string `json:"title"`
}

type RespOrderInfo struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
	//	Data   OrderInfo `json:"data"`

	Data interface{} `json:"data"`
}

func changestruct(this *models.OrderHouse) interface{} {
	orderInfo := map[string]interface{}{
		"amount":     this.Amount,
		"comment":    this.Comment,
		"ctime":      this.Ctime,
		"days":       this.Days,
		"end_date":   this.End_date,
		"img_url":    models.AddDomain2Url(this.House.Index_image_url),
		"order_id":   this.Id,
		"start_date": this.Begin_date,
		"status":     this.Status,
		"title":      this.House.Title,
	}
	return orderInfo
}
func (this *GetOrderController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}
func (this *GetOrderController) GetOrder() {
	resp := RespOrderInfo{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	resp_user_id := this.GetSession("user_id")
	fmt.Println("use_id:", resp_user_id)
	var role string
	this.Ctx.Input.Bind(&role, "role")
	if role == "" {
		resp.Errno = models.RECODE_ROLEERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	o := orm.NewOrm()
	oreders := []models.OrderHouse{}
	houselist := []interface{}{}
	if "custom" == role {
		beego.Info("我是客户")
		o.QueryTable("order_house").Filter("user_id", resp_user_id).OrderBy("ctime").All(&oreders)

	} else {

		beego.Info("我是房东")
		houses := []models.House{}
		o.QueryTable("house").Filter("user_id", resp_user_id).OrderBy("ctime").All(&houses)
		var houseids []int
		for houseid, _ := range houses {
			houseids = append(houseids, houseid)
		}
		o.QueryTable("order_house").Filter("house_id__in", houseids).OrderBy("ctime").All(&oreders)
	}
	for _, oreder := range oreders {
		housedata := changestruct(&oreder)
		houselist = append(houselist, housedata)

	}
	data := map[string]interface{}{}
	data["orders"] = houselist
	beego.Info("oreders:", data)
	resp.Data = data
	return
}
