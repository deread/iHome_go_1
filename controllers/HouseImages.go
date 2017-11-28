package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/cache"
	//_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//  "time"
	"path"
	"strconv"
)

type HouseImagesUrl struct {
	Url string `json:"url"`
}

// 上传房源图片的返回结构
type HouseImagesResp struct {
	Errno  string         `json:"errno"`
	Errmsg string         `json:"errmsg"`
	Data   HouseImagesUrl `json:"data"`
}

type HouseImagesController struct {
	beego.Controller
}

func (this *HouseImagesController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *HouseImagesController) UploadHouseImages() {

	resp := HouseImagesResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//获取房源图片数据
	file, header, err := this.GetFile("house_image")

	if err != nil {
		resp.Errno = models.RECODE_SERVERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("upload house images error")
		return
	}
	defer file.Close()

	//将图片二进制数据存储到fastDFS中，得到fileID
	//创建一个文件的缓冲
	fileBuffer := make([]byte, header.Size)

	_, err = file.Read(fileBuffer)
	if err != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("read house images error")
		return
	}

	//home1.jpg
	suffix := path.Ext(header.Filename) // suffix = ".jpg"
	groupName, fileId, err1 := models.FDFSUploadByBuffer(fileBuffer, suffix[1:])
	if err1 != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("fdfs upload  file error")
		return
	}

	beego.Info("groupName : ", groupName, " fileId : ", fileId)

	//从请求 url 中得到 house_id //http:IP:PORT/api/v1.0/houses/3/images 备注:3 表示 房源ID
	house_id_str := this.Ctx.Input.Param(":id")

	beego.Info("============", house_id_str, "=============")
	//查看该房屋中的 index_image_url 显示图片是否为空
	o := orm.NewOrm()
	house_id, err := strconv.Atoi(house_id_str)
	if err != nil {
		fmt.Println("HouseImage strcon.Atoi is err")
		return
	}
	//如果为空，将 index_image_url 设置为此图片的 fileID
	//将图片的 fileID 追加到 HouseImage 字段并入库
	house := models.House{Id: house_id}
	if o.Read(&house) == nil {
		if house.Index_image_url == "" {
			house.Index_image_url = fileId
			if _, err := o.Update(&house, "Index_image_url"); err != nil {
				beego.Info("Updata Index_image_url err")
			}
		} else {
			//添加 id\url\house_id 字段到 house_image 数据库中
			house_images := models.HouseImage{Url: fileId, House: &house}
			if _, err := o.Insert(&house_images); err != nil {
				resp.Errno = models.RECODE_DBERR
				resp.Errmsg = models.RecodeText(resp.Errno)
				beego.Info("Insert HouseImage err")
			}
		}
	}

	//拼接一个完整的路径
	house_images_url := models.AddDomain2Url(fileId)
	beego.Info(house_images_url)
	//avatar_url := "http://39.106.110.44:8080/" + fileId

	resp.Data.Url = house_images_url
	return
}
