package store

import (
	"github.com/gin-gonic/gin"
	. "auroraLZDF/member_api/app/Controllers"
	"auroraLZDF/member_api/utils"
	"auroraLZDF/member_api/app/Models"
)

var model Models.Store

func Index(c *gin.Context) {
	query := c.Request.URL.Query() // 获取所有 get 参数

	//var params map[string]int	// 报错: assignment to entry in nil map
	var params = make(map[string]int)	// make 声明，初始化变量，给变量赋了默认值
	for key, value := range query {
		params[key] = utils.StringToInt(value[0])
	}

	data, count := model.Find(params)

	var storeIds []string
	var result = make(map[int]map[string]interface{})
	for key, value := range data {
		extend := utils.JsonToMap(value.Extend)

		result[key]["id"] = value.Id
		result[key]["user_id"] = value.UserId
		result[key]["store_id"] = value.StoreId
		result[key]["store_type"] = value.StoreType
		result[key]["created_at"] = value.CreatedAt
		// extend
		result[key]["store_name"] = extend["store_name"]
		result[key]["en_store_name"] = extend["en_store_name"]
		result[key]["type"] = extend["type"]
		result[key]["region_id"] = extend["region_id"]
		result[key]["region_name"] = extend["region_name"]
		result[key]["region_en_name"] = extend["region_en_name"]
		result[key]["address"] = extend["address"]
		result[key]["en_address"] = extend["en_address"]
		result[key]["contact_name"] = extend["contact_name"]
		result[key]["contact_phone"] = extend["contact_phone"]
		result[key]["level_name"] = extend["level_name"]
		result[key]["company_tel"] = extend["company_tel"]
		result[key]["headpic_id"] = extend["headpic_id"]
		result[key]["sgs"] = extend["sgs"]
		result[key]["sell_type"] = extend["sell_type"]
		result[key]["bussinessType"] = extend["bussinessType"]
		result[key]["bussinessTypeCn"] = extend["bussinessTypeCn"]

		storeIds[key] = utils.IntToString(value.StoreId)
	}

	// 批量判断店铺是否vip
	vips := GetStoreVip(storeIds)
	for key, value := range result {
		vip := vips[value["store_id"].(string)].(string)
		if utils.Required(vip) {
			result[key]["vip"] = vip
		} else {
			result[key]["vip"] = 0
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg": "success",
		"data": map[string]interface{}{
			"data": result,
			"total": count,
		},
	})
	return
}

func Types(c *gin.Context) {
	userId := c.Query("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code":500, "msg": "会员不存在"})
		return
	}

	data := model.FindStoreTypeByUserId(userId)

	c.JSON(200, gin.H{"code": 200, "msg": "success", "data":data})
	return
}

func Create(c *gin.Context) {
	storeId := c.Param("storeId")
	userId := c.PostForm("user_id")
	storeType := c.PostForm("type")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(storeType) {
		c.JSON(200, gin.H{"code": 500, "msg": "店铺类型不能为空"})
		return
	}

	storeInfo := GetStoreInfo(storeId)
	if storeInfo == nil {
		c.JSON(200, gin.H{"code": 500, "msg": "获取店铺信息失败"})
		return
	}

	var companyInfo map[string]interface{}
	//log.Println("storeInfo: ", storeInfo)
	if storeInfo["company_info"] != "" {
		companyInfo = storeInfo["company_info"].(map[string]interface{})
	}

	extend := map[string]interface{}{
		"store_name":      storeInfo["store_name"],
		"en_store_name":   storeInfo["en_store_name"],
		"type":            storeInfo["type"],
		"region_id":       storeInfo["region_id"],
		"region_name":     storeInfo["region_name"],
		"region_en_name":  storeInfo["region_en_name"],
		"address":         storeInfo["address"],
		"en_address":      storeInfo["en_address"],
		"contact_name":    storeInfo["contact_name"],
		"contact_phone":   storeInfo["contact_phone"],
		"level_name":      storeInfo["level_name"],
		"company_tel":     storeInfo["company_tel"],
		"headpic_id":      storeInfo["headpic_id"],
		"sgs":             storeInfo["sgs"],
		"sell_type":       storeInfo["sell_type"],
		"bussinessType":   companyInfo["bussinessType"],
		"bussinessTypeCn": companyInfo["bussinessTypeCn"],
	}

	if _, err := model.FindOne(storeId, userId); err != nil {
		model.Create(storeId, userId, storeType, extend)
	}

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func Cancel(c *gin.Context) {
	storeId := c.Param("storeId")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	//model.Cancel(utils.StringToInt(userId), utils.StringToInt(storeId))
	model.BatchCancel(userId, storeId)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
}

func BatchCancel(c *gin.Context) {
	userId := c.PostForm("user_id")
	storeIds := c.PostForm("store_ids")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(storeIds) {
		c.JSON(200, gin.H{"code": 500, "msg": "店铺 ID 不能为空"})
		return
	}

	model.BatchCancel(userId, storeIds)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}
