package moldata

import (
	"github.com/auroraLZDF/gin-api-demo/app/Models"
	"github.com/gin-gonic/gin"
	"github.com/auroraLZDF/gin-api-demo/utils"
	. "github.com/auroraLZDF/gin-api-demo/app/Controllers"
)

var model Models.MolData

func Index(c *gin.Context) {
	params := map[string]interface{}{
		"user_id": utils.StringToInt(c.Query("user_id")),
		"mol_id": utils.StringToInt(c.Query("mol_id")),
		"category_id": utils.StringToInt(c.Query("category_id")),
		"page": utils.StringToInt(c.DefaultQuery("page", "1")),
		"page_size": utils.StringToInt(c.DefaultQuery("page_size", "15")),
		"begin_time": c.Query("begin_time"),
		"end_time": c.Query("end_time"),
	}

	data, count := model.Find(params)

	var result = make(map[int]map[string]interface{})
	for key, value := range data {
		extend := utils.JsonToMap(value.Extend)

		data[key].NameZh = extend["name_zh"].(string)
		data[key].NameEn = extend["name_en"].(string)
		data[key].CasNo = extend["cas_no"].(string)
		data[key].Img = extend["img"].(string)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"data":  result,
			"total": count,
		},
	})
	return
}

func Create(c *gin.Context) {
	molId := c.Param("molId")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	molInfo := GetMolInfo(molId)
	if molInfo == nil {
		c.JSON(200, gin.H{"code": 500, "msg": "获取百科详情失败"})
		return
	}

	extend := map[string]interface{}{
		"name_zh": molInfo["cnName"],
		"name_en": molInfo["enName"],
		"cas_no": molInfo["cas"],
		"img": molInfo["structImage"],
	}

	if _, err := model.FindOne(molId, userId); err != nil {
		model.Create(utils.StringToInt(molId), utils.StringToInt(userId), utils.MapToJson(extend))
	}

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func Cancel(c *gin.Context) {
	molId := c.Param("molId")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	model.BatchCancel(userId, molId)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func BatchCancel(c *gin.Context) {
	molIds := c.PostForm("mol_ids")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	model.BatchCancel(userId, molIds)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func Categories(c *gin.Context) {
	userId := c.PostForm("user_id")
	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	data := model.FindCategoriesByUserId(userId)

	c.JSON(200, gin.H{"code": 200, "data": data})
	return
}

func HasCollection(c *gin.Context) {
	userId := c.Param("userId")
	molIds := c.Query("mol_ids")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(molIds) {
		c.JSON(200, gin.H{"code": 500, "msg": "关注百科 ID 不能为空"})
	}

	data, _ := model.HasCollection(userId, molIds)

	c.JSON(200, gin.H{"code": 200, "msg": "批量获取关注状态成功", "data": data})
	return
}
