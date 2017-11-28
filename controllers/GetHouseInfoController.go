package controllers

import (
	_ "encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//"strconv"
	_ "time"
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
	Data   interface{} `json:"data"`
}

type Area struct {
	Id   int    `json:"aid"`                  //区域编号
	Name string `orm:"size(32)" json:"aname"` //区域名字
}

//根据areaid-->AreaName
func GetAreaName(myid interface{}) (RAreaName string) {
	o := orm.NewOrm()
	num := myid.(int)
	//num, _ := strconv.Atoi(myid.(string))
	area := models.Area{Id: num}
	o.Read(&area)

	RAreaName = area.Name
	return
}
func Struct2house(this *models.House) interface{} {
	house_info := map[string]interface{}{
		"address":   this.Address,
		"area_name": GetAreaName(this.Area.Id),
		//"ctime":       this.Ctime.Format("2006-01-02 15:04:05"),
		"ctime":       this.Ctime,
		"house_id":    this.Id,
		"img_url":     this.Images,
		"order_count": this.Order_count,
		"price":       this.Price,
		"room_count":  this.Room_count,
		"title":       this.Title,
		"user_avatar": models.AddDomain2Url(this.User.Avatar_url),
	}
	return house_info
}
func (this *GetHouseInfoController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *GetHouseInfoController) GetHouseInfo() {
	resp := RespHouseInfo{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
	//通过session得到user_id
	resp_user_id := this.GetSession("user_id")

	fmt.Println("===>houses", resp_user_id)
	//查询house表，找到user_id所有的房屋

	o := orm.NewOrm()

	houses := []models.House{}
	house_list := []interface{}{}
	qs := o.QueryTable("house")
	qs.Filter("user_id", resp_user_id).All(&houses)
	//fmt.Printf("resqs ===>%+v", houses)
	for _, house := range houses {
		housedata := Struct2house(&house)
		house_list = append(house_list, housedata)

	}

	data := map[string]interface{}{}
	data["houses"] = house_list
	fmt.Printf("housesInfo--->%+v", data)
	//返回json数据
	resp.Data = data
	return
}
