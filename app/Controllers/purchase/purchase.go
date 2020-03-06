package purchase

import (
	"github.com/auroraLZDF/gin-api-demo/app/Models"
	"github.com/gin-gonic/gin"
	"github.com/auroraLZDF/gin-api-demo/utils"
	. "github.com/auroraLZDF/gin-api-demo/app/Controllers"
)

var model Models.Purchase

func Index(c *gin.Context)  {
	userId := c.Query("user_id")
	params := map[string]interface{}{
		"user_id": utils.StringToInt(userId),
		"province_name": c.Query("province_name"),
		"city_name": c.Query("city_name"),
		"number_unit": utils.StringToInt(c.Query("number_unit")),
		"sort": utils.StringToInt(c.Query("sort")),
		"page": utils.StringToInt(c.DefaultQuery("page", "1")),
		"page_size": utils.StringToInt(c.DefaultQuery("page_size", "15")),
	}

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	data, count := model.Find(params)

	if count == 0 {
		c.JSON(200, gin.H{"code": 500, "msg": "该用户还没有关注任何采购信息"})
		return
	}

	var result = make(map[int]map[string]interface{})
	for key, value := range data {
		data[key].Number = value.Num
		data[key].NumberUnit = value.NumUnit
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

func Create(c *gin.Context)  {
	code := c.Param("code")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if _, err := model.FindOne(userId, code); err == nil {
		c.JSON(200, gin.H{"code": 200, "msg": "已关注"})
		return
	}

	params := map[string]interface{}{
		"product_name": c.PostForm("product_name"),
		"cas": c.PostForm("cas"),
		"number": c.PostForm("number"),
		"number_unit": c.PostForm("number_unit"),
		"purity": c.PostForm("purity"),
		"province_name": c.PostForm("province_name"),
		"city_name": c.PostForm("city_name"),
		"remarks": c.PostForm("remarks"),
		"state": c.PostForm("state"),
		"creation_time": c.PostForm("creation_time"),
	}

	model.Create(utils.StringToInt(userId), code, params)

	c.JSON(200, gin.H{"code": 200, "msg": "关注成功"})
	return
}

func BranchCancel(c *gin.Context)  {
	userId := c.Param("userId")
	codes := c.PostForm("codes")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(codes) {
		c.JSON(200, gin.H{"code": 500, "msg": "关注采购 ID 不能为空"})
		return
	}

	model.BatchCancel(userId, codes)

	c.JSON(200, gin.H{"code": 200, "msg": "批量取消关注成功"})
	return
}

func HasCollection(c *gin.Context)  {
	userId := c.Param("userId")
	codes := c.Query("codes")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(codes) {
		c.JSON(200, gin.H{"code": 500, "msg": "关注采购 ID 不能为空"})
	}

	data, _ := model.HasCollection(userId, codes)

	c.JSON(200, gin.H{"code": 200, "msg": "批量获取关注状态成功", "data": data})
	return
}
