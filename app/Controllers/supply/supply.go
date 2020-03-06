package supply

import (
	"github.com/auroraLZDF/gin-api-demo/app/Models"
	"github.com/gin-gonic/gin"
	"github.com/auroraLZDF/gin-api-demo/utils"
	. "github.com/auroraLZDF/gin-api-demo/app/Controllers"
)

var model Models.Supply

func Index(c *gin.Context)  {
	userId := c.Query("user_id")
	params := map[string]interface{}{
		"user_id": utils.StringToInt(userId),
		"province_name": c.Query("province_name"),
		"city_name": c.Query("city_name"),
		"spec_unit": utils.StringToInt(c.Query("spec_unit")),
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
		c.JSON(200, gin.H{"code": 500, "msg": "该用户还没有关注任何供应信息"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"data":  data,
			"total": count,
		},
	})
	return
}

func Create(c *gin.Context)  {
	supplyId := c.Param("supplyId")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if _, err := model.FindOne(userId, supplyId); err == nil {
		c.JSON(200, gin.H{"code": 200, "msg": "已关注"})
		return
	}

	params := map[string]interface{}{
		"product_name": c.PostForm("product_name"),
		"purity": c.PostForm("purity"),
		"spec_count": c.PostForm("spec_count"),
		"spec_unit": c.PostForm("spec_unit"),
		"spec_package": c.PostForm("spec_package"),
		"price": c.PostForm("price"),
		"price_unit": c.PostForm("price_unit"),
		"product_level": c.PostForm("product_level"),
		"province_name": c.PostForm("province_name"),
		"city_name": c.PostForm("city_name"),
		"images": c.PostForm("images"),
		"period": c.PostForm("period"),
		"creation_time": c.PostForm("creation_time"),
	}

	model.Create(utils.StringToInt(userId), utils.StringToInt(supplyId), params)

	c.JSON(200, gin.H{"code": 200, "msg": "关注成功"})
	return
}

func BranchCancel(c *gin.Context)  {
	userId := c.Param("userId")
	supplyIds := c.PostForm("supply_ids")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(supplyIds) {
		c.JSON(200, gin.H{"code": 500, "msg": "关注供应 ID 不能为空"})
		return
	}

	model.BatchCancel(userId, supplyIds)

	c.JSON(200, gin.H{"code": 200, "msg": "批量取消关注成功"})
	return
}

func HasCollection(c *gin.Context)  {
	userId := c.Param("userId")
	supplyIds := c.Query("supply_ids")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(supplyIds) {
		c.JSON(200, gin.H{"code": 500, "msg": "关注供应 ID 不能为空"})
	}

	data, _ := model.HasCollection(userId, supplyIds)

	c.JSON(200, gin.H{"code": 200, "msg": "批量获取关注状态成功", "data": data})
	return
}
