package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	//	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	"iHome_go_1/models"
	"strconv"
	"time"
)

type GetHouseDetailInfoController struct {
	beego.Controller
}

type GetHouseDetailInfoReap struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type HouseDetailInfo struct {
	House   interface{} `json:"house"`
	User_id int         `json:"user_id"`
}

// 将数据返回
func (this *GetHouseDetailInfoController) ResData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// 得到一个房屋的信息信息
func (this *GetHouseDetailInfoController) GetHouseDetailInfo() {

	// 1. 创建JSON对象, 并设置在函数结束时调用ResData来返回数据
	resp := GetHouseDetailInfoReap{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.ResData(&resp)

	// 2. 从session获取用户Id
	uid := this.GetSession("user_id")
	if uid == nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3. 从url中获取房屋Id
	hid := this.Ctx.Input.Param(":id")
	if hid == "" {
		resp.Errno = models.RECODE_DATAEXIST
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 4. 从redis通过房屋Id获取房屋信息
	redis_conn, err := cache.NewCache("redis", `{"key": "iHome_go", "conn": "127.0.0.1:6400", "dbNum": "0"}`)
	if err != nil {

		// 4.1 没有获取到Redis缓存表示redis有问题
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	} else {

		// 4.2 从redis缓冲中查看是否有字段
		var house_redis interface{}
		house_redis = redis_conn.Get("house_info_" + hid)
		if house_redis != nil {
			beego.Debug("从redis获取到房屋数据")

			// 4.2.1 把Json转换成数据
			var house_redis_json interface{}
			err := json.Unmarshal(house_redis.([]byte), &house_redis_json)
			if err != nil {
				resp.Errno = models.RECODE_DATAERR
				resp.Errmsg = models.RecodeText(resp.Errno)
				return
			}

			// 4.2.2 将数据返回
			resp.Data = HouseDetailInfo{
				House:   house_redis_json,
				User_id: uid.(int),
			}
			return
		}
	}

	// 5. 从Mysql获取房屋信息
	// 5.1 获取数据库句柄
	o := orm.NewOrm()
	house := models.House{}

	// 5.2 查询house这张表, 将查询结果放到house中
	house.Id, _ = strconv.Atoi(hid)

	if err := o.Read(&house); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	if _, err := o.LoadRelated(&house, "Area"); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	if _, err := o.LoadRelated(&house, "User"); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	if _, err := o.LoadRelated(&house, "Images"); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	if _, err := o.LoadRelated(&house, "Facilities"); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 6. 将数据写入到redis中(保存60分钟)
	house_json := house.To_one_house_desc()
	house_mysql_json, _ := json.Marshal(house_json)
	err = redis_conn.Put("houser_info_"+hid, house_mysql_json, time.Second*3600)

	// 7. 将数据赋给resp, 将数据返回
	resp.Data = HouseDetailInfo{
		House:   house_json,
		User_id: uid.(int),
	}

	return
}
