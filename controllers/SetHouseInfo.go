package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	//	//"github.com/astaxie/beego/cache"
	//	//_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"strconv"
	//	//	"time"
	//	"path"
)

//房源信息请求的数据
type SetHouseInfoRequest struct {
	Title      string   `json:"title"`
	Price      string   `json:"price"`
	Area_id    string   `json:"area_id"`
	Address    string   `json:"address"`
	Room_count string   `json:"room_count"`
	Acreage    string   `json:"acreage"`
	Unit       string   `json:"unit"`
	Capacity   string   `json:"capacity"`
	Beds       string   `json:"beds"`
	Deposit    string   `json:"deposit"`
	Min_days   string   `json:"min_days"`
	Max_days   string   `json:"max_days"`
	Facility   []string `json:"facility"`
}
type SetHouseInfo struct {
	House_id int64 `json:"house_id"`
}

//房源信息业务回复
type SetHouseInfoResp struct {
	Errno  string       `json:"errno"`
	Errmsg string       `json:"errmsg"`
	Data   SetHouseInfo `json:"data"`
}
type HousesController struct {
	beego.Controller
}

func (this *HousesController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *HousesController) ReleaseHouses() {
	resp := SetHouseInfoResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)
	//得到房源信息
	//request
	var request_data SetHouseInfoRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	fmt.Printf("request data : %+v\n", request_data)
	//插入房源数据到house表中
	///*
	house := models.House{}
	house.Title = request_data.Title
	house.Price, _ = strconv.Atoi(request_data.Price)
	house.Price = house.Price * 100
	area_id, _ := strconv.Atoi(request_data.Area_id)
	area := models.Area{Id: area_id}
	house.Area = &area
	house.Address = request_data.Address
	house.Room_count, _ = strconv.Atoi(request_data.Room_count)
	house.Acreage, _ = strconv.Atoi(request_data.Acreage)
	house.Unit = request_data.Unit
	house.Capacity, _ = strconv.Atoi(request_data.Capacity)
	house.Beds = request_data.Beds
	house.Deposit, _ = strconv.Atoi(request_data.Deposit)
	house.Deposit = house.Deposit * 100
	house.Min_days, _ = strconv.Atoi(request_data.Min_days)
	house.Max_days, _ = strconv.Atoi(request_data.Max_days)
	user := models.User{Id: this.GetSession("user_id").(int)}
	house.User = &user
	//	house.Facilities = request_data.Facility
	o := orm.NewOrm()

	houseid, err := o.Insert(&house)
	if err != nil {
		fmt.Println("insert houseinfo error = ", err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info("SetHouseInfo  insert succ id = ", houseid)
	//插入facility和house的多对多关系到表中
	facilities := []*models.Facility{}
	for _, fid := range request_data.Facility {
		id, _ := strconv.Atoi(fid)
		facility := &models.Facility{Id: id}
		facilities = append(facilities, facility)

	}

	m2mhouse_facility := o.QueryM2M(&house, "Facilities")
	num, err := m2mhouse_facility.Add(facilities)
	if err != nil {
		fmt.Println("insert facilities error = ", err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info("num :", num)
	resp.Data = SetHouseInfo{House_id: houseid}
	return
}
