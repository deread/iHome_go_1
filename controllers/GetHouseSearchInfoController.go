package controllers

//-----------获取搜索房源信息
import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
)

type GetHouseSearchInfoController struct {
	beego.Controller
}
type GetStructHouses struct {
	Aid int    `json:"aid"`
	Sd  string `json:"sd"`
	Ed  string `json:"ed"`
	Sk  string `json:"sk"`
	P   string `json:"p"`
}

//用于搜索房屋返回结构体
type RespHouseSearchInfo struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func (this *GetHouseSearchInfoController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}
func (this *GetHouseSearchInfoController) GetHouseSearchInfo() {
	resp := RespHouseSearchInfo{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	//获取前端的参数
	//var ReqData GetStructHouses
	var aid int
	this.Ctx.Input.Bind(&aid, "aid")
	var sd string
	this.Ctx.Input.Bind(&sd, "sd")
	var ed string
	this.Ctx.Input.Bind(&ed, "ed")
	var sk string
	this.Ctx.Input.Bind(&sk, "sk")
	var page int
	this.Ctx.Input.Bind(&page, "p")

	//把时间从str转换成字符串格式

	//校验开始时间一定要早于结束时间

	//判断page的合法性 一定是大于0的整数
	if page <= 0 {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(models.RECODE_PARAMERR)
		return
	}
	//尝试从redis中获取数据, 有则返回
	bm, err := cache.NewCache("redis", `{"key":"ihome_go_1","conn":":6400","dbNum":"0"}`)
	if err != nil {
		//resp.Errno = models.RECODE_DBERR
		//resp.Errmsg = models.RecodeText(models.RECODE_DBERR)
		//return
	}
	if bm.IsExist("house") {
		resp.Errno = models.RECODE_DATAEXIST
		resp.Errmsg = models.RecodeText(models.RECODE_DATAEXIST)
		//bm.Get(
		//return
	}
	//获取redis数据(未完成)

	//如果没有 从mysql中查询
	housearray := []models.House{}
	o := orm.NewOrm()
	num, err := o.QueryTable("house").Filter("area_id", aid).All(&housearray)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(models.RECODE_DBERR)
		return
	}
	total_page := int(num)/models.HOUSE_LIST_PAGE_CAPACITY + 1
	current_page := 1

	house_list := []interface{}{}
	for _, house := range housearray {
		housedata := Struct2house(&house)
		house_list = append(house_list, housedata)
	}

	data := map[string]interface{}{}
	data["houses"] = house_list
	data["total_page"] = total_page
	data["current_page"] = current_page
	resp.Data = data
	//返回前端
	fmt.Printf("this is house info--==========>%+v\n", data)

	defer this.RetData(&resp)
}
