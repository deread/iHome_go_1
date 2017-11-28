package routers

import (
	"github.com/astaxie/beego"
	"iHome_go_1/controllers"
)

func init() {

	beego.Router("/", &controllers.MainController{})
	//注册[post]
	beego.Router("api/v1.0/users", &controllers.UserController{}, "post:Reg")
	//登陆[post]
	beego.Router("api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	//请求session[get]
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:Get")
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "delete:Delete")
	//请求地域[get]
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
	//请求用户基本信息[get]
	beego.Router("api/v1.0/user", &controllers.UserInfoController{}, "get:UserInfoGet")
	//请求上传头像[post]
	beego.Router("api/v1.0/user/avatar", &controllers.UserController{}, "post:GetAvatar")
	//请求更新用户名[put]
	beego.Router("api/v1.0/user/name", &controllers.UserController{}, "put:UpdateUsername")
	//实名认证检查[get]
	beego.Router("api/v1.0/user/auth", &controllers.AuthController{}, "get:AuthCheck")
	//更新实名认证信息[post]
	beego.Router("api/v1.0/user/auth", &controllers.AuthController{}, "post:UpdateAuthinfo")
	//发布房源信息[post]
	beego.Router("api/v1.0/houses", &controllers.HousesController{}, "post:ReleaseHouses")
	//上传房源图片信息[post]？为房源ID
	beego.Router("api/v1.0/houses/?:id/images", &controllers.HouseImagesController{}, "post:UploadHouseImages")
	//请求当前用户已经发布的房源信息[get]
	beego.Router("api/v1.0/user/houses", &controllers.GetHouseInfoController{}, "get:GetHouseInfo")
	//请求房源详细信息[get] ?-->房源id
	beego.Router("api/v1.0/houses/:id", &controllers.GetHouseDetailInfoController{}, "get:GetHouseDetailInfo")
	//请求房屋搜索信息[get]
	beego.Router("api/v1.0/houses", &controllers.GetHouseSearchInfoController{}, "get:GetHouseSearchInfo")
	//请求查看房东/租客订单信息[get]
	//beego.Router("api/v1.0/orders", &controllers.GetOrderController{}, "get:GetOrder")
	//房东同意/拒绝订单[put] ? -->订单id
	//beego.Router("api/v1.0/orders/?/status", &controllers.PutOrderStatController{}, "put:PutOrderStat")
	//用户评价订单信息 [put] ?-->订单Id
	//beego.Router("api/v1.0/orders/?/comment", &controllers.PutOrderCommController{}, "put:PutOrderComm")

}
