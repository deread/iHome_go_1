package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	"iHome_go_1/models"
	"strconv"
	"time"
)

type PostOrderController struct {
	beego.Controller
}

// 返回客户端结构体
type PostOrderReap struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

// 接收客户端请求数据结构体
type OrderPostInfo struct {
	House_id   string `json:"house_id"`
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
}

// 将数据返回
func (this *PostOrderController) ResData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// 发布订单
func (this *PostOrderController) PostOrder() {

	// 1. 创建JSON对象, 并设置在函数结束时调用ResData来返回数据
	resp := PostOrderReap{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.ResData(&resp)

	// 2. 从session获取用户Id
	uid := this.GetSession("user_id")
	if uid == nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3. 从Post请求数据中拿到订单部分数据
	var request_data OrderPostInfo
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	beego.Info("request data: ", request_data)

	// 4. 判断json数据是否有效
	hid := request_data.House_id
	startDate, err1 := time.Parse("2006-01-02", request_data.Start_date)
	endDate, err2 := time.Parse("2006-01-02", request_data.End_date)
	//startDate, err1 := time.Parse("2006-01-02 15:04:05", request_data.Start_date)
	//endDate, err2 := time.Parse("2006-01-02 15:04:05", request_data.End_date)
	if hid == "" || err1 != nil || err2 != nil {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Error("Json数据错误")
		return
	}

	// 5. 判断end_date是否在start_date之后
	beego.Info("租房开始时间: ", startDate)
	beego.Info("租房结束时间: ", endDate)
	if startDate.After(endDate) {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Error("租房开始时间>租房结束时间")
		return
	}

	// 6. 得到一共入住的天数
	diff := endDate.Sub(startDate)
	arrDay := diff.Hours()/24 + 1
	beego.Info("入住天数: ", arrDay)
	if arrDay <= 0 {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Error("入住天数<=0,", diff.Hours())
		return
	}

	// 7. 根据House_Id从Mysql查询房屋信息
	// 7.1 获取数据库句柄
	o := orm.NewOrm()

	// 7.2 查询house这张表, 将查询结果放到house中
	house := models.House{}
	house.Id, _ = strconv.Atoi(hid)

	if err := o.Read(&house); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Error("查询数据库出错")
		return
	}

	// 8. 判断当前的user_id和房屋所有者Id是否相同
	if uid.(int) == house.User.Id {
		resp.Errno = models.RECODE_ROLEERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Error("房客和房东一致")
		return
	}
	beego.Info("用户user_id为: ", uid, "房东user_id为:", house.User.Id)

	// 9. 确保用户选择的房屋未被预定, 日期上无冲突

	// 10. 封装订单数据对象信息
	order := models.OrderHouse{}

	order.User = &models.User{Id: uid.(int)}
	order.House = &models.House{Id: house.Id}
	order.Begin_date = startDate
	order.End_date = endDate
	order.Days = int(arrDay)
	order.House_price = house.Price
	order.Amount = house.Price * order.Days
	order.Ctime = time.Now()
	order.Status = models.ORDER_STATUS_WAIT_ACCEPT

	// 11. 将订单数据插入到Mysql中
	oid, err := o.Insert(&order)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Error("插入订单数据到Mysql错误")
		return
	}
	beego.Info("将订单数据插入到Mysql成功, 返回的order_id = ", oid)

	// 12. 返回order_id
	oidMap := map[string]string{
		"order_id": strconv.FormatInt(oid, 10),
	}
	resp.Data = oidMap

	return
}
