package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"strconv"
)

type OrdercommController struct {
	beego.Controller
}

type OrdercommResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
	//Data   interface{} `json:"data"`
}

//返回函数
func (this *OrdercommController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *OrdercommController) PutOrderComment() {
	resp := OrdercommResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

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

	comment := req["comment"].(string)
	fmt.Println("comment-->", comment)
	//校验action合法性
	if comment == "" {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	o := orm.NewOrm()
	res, err := o.Raw("UPDATE order_house SET comment = ? where id = ?", comment, int_order_id).Exec()
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	num, _ := res.RowsAffected()
	fmt.Println("mysql row affected nums: ", num)
	return
}
