package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"time"
)

type GetHouseIndexController struct {
	beego.Controller
}

type GetHouseIndexResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

//返回函数
func (this *GetHouseIndexController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *GetHouseIndexController) GetHouseIndex() {
	resp := GetHouseIndexResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//从redis中获取数据
	rcache, err := cache.NewCache("redis", `{"key":"ihome_go_1","conn":":6400","dbNum":"0"}`)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(models.RECODE_DBERR)
	}
	housedata := rcache.Get("home_page_data")
	rData := []interface{}{}
	if housedata != nil {
		fmt.Println("============get house data from redis============")
		json.Unmarshal(housedata.([]byte), &rData)
		resp.Data = rData
		return
	}

	//如果没有，从mysqll数据库中获取数据
	houses := []models.House{}
	o := orm.NewOrm()
	_, err = o.QueryTable("house").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses)
	if err == nil {
		for _, house := range houses {
			o.LoadRelated(&house, "Area")
			o.LoadRelated(&house, "User")
			o.LoadRelated(&house, "Images")
			o.LoadRelated(&house, "Facilities")
			//添加image_url的ip:port
			fmt.Println("Index_image_url-->", house.Index_image_url)
			//struct -->house类型
			housedata := Struct2house(&house)

			fmt.Println("--=--==-=-=", housedata)
			rData = append(rData, housedata)
		}
	}
	//返回前端
	resp.Data = rData
	//将请求字段放入到redis中，以便下次使用
	house_value, _ := json.Marshal(rData)
	rcache.Put("house_page_data", house_value, 3600*time.Second)
	return
}
