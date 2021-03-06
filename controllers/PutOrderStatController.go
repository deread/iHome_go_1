package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"strconv"
)

type PutOrderStatController struct {
	beego.Controller
}

type PutOrderStatResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
	//Data   interface{} `json:"data"`
}

//返回函数
func (this *PutOrderStatController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//api/v1.0/orders/:id/status-->struct
/*
type StructAction struct {
	Action string `json:"action"` //accept  ,reject
}
*/

func (this *PutOrderStatController) PutOrderStat() {
	resp := PutOrderStatResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
	//通过Session得到当前user_id
	v := this.GetSession("user_id")
	if v == nil {
		resp.Errno = models.RECODE_SESSIONERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//通过url得到当前订单ID
	//order_id := this.Input().Get("id")
	order_id := this.Ctx.Input.Param(":id")
	int_order_id, err := strconv.Atoi(order_id)
	fmt.Println("int_order_id-->", int_order_id)
	//解析json数据，得到action

	var req map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	action := req["action"]
	fmt.Println("action-->", action)
	//校验action合法性
	if action != "accept" && action != "reject" {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	//查找订单表order_house，得到订单并确定订单状态是WAIT_ACCEPT
	OrderHouse := models.OrderHouse{}
	o := orm.NewOrm()
	err = o.QueryTable("order_house").Filter("id", int_order_id).One(&OrderHouse)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	/*
		//校验订单的user_id是否是当前用户的user_id
		if v.(int) != OrderHouse.User.Id {
			resp.Errno = models.RECODE_USERERR
			resp.Errmsg = models.RecodeText(resp.Errno)
			return
		}
	*/
	//校验houseOrder状态
	Order_stat := OrderHouse.Status
	if Order_stat != models.ORDER_STATUS_WAIT_ACCEPT {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	//if action=accept-->更换订单status状态WAIT_COMMIT等待用户评价
	if action == "accept" {
		OrderHouse.Status = models.ORDER_STATUS_WAIT_COMMENT
		//else  action=reject-->更换订单状态status状态REJECT
	} else if action == "reject" {
		OrderHouse.Status = models.ORDER_STATUS_REJECTED
		//从url中获取reason参数字段
		Rjct_reason := req["reason"]
		//将reason字段添加到order的评价Comment字段中
		OrderHouse.Comment = Rjct_reason.(string)
		fmt.Println("OrderHouse_comment-->", OrderHouse)
	}

	//更新mysql数据库中的订单
	_, err = o.Update(&OrderHouse)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	//返回前端
	return
}
