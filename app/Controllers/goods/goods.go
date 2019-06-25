package goods

import (
	"github.com/gin-gonic/gin"
	"auroraLZDF/member_api/app/Models"
	"auroraLZDF/member_api/utils"
	. "auroraLZDF/member_api/app/Controllers"
)

var model Models.Goods

func Index(c *gin.Context) {
	params := map[string]interface{}{
		"user_id":    utils.StringToInt(c.Query("user_id")),
		"goods_id":   utils.StringToInt(c.Query("goods_id")),
		"brand_id":   utils.StringToInt(c.Query("brand_id")),
		"page_size":  utils.StringToInt(c.DefaultQuery("page_size", "15")),
		"page":       utils.StringToInt(c.DefaultQuery("page", "1")),
		"begin_time": c.Query("begin_time"),
		"end_time":   c.Query("end_time"),
	}

	data, count := model.Find(params)

	var query = make(map[int]map[string]int)
	var result = make(map[int]map[string]interface{})
	for key, value := range data {
		extend := utils.JsonToMap(value.Extend)

		result[key]["id"] = value.Id
		result[key]["user_id"] = value.UserId
		result[key]["store_id"] = value.StoreId
		result[key]["store_name"] = value.StoreName
		result[key]["goods_id"] = value.GoodsId
		result[key]["sku_id"] = value.SkuId
		//result[key]["brand_id"] = value.BrandId
		//result[key]["brand_name"] = value.BrandName
		result[key]["price"] = value.Price
		result[key]["price_unit"] = value.PriceUnit
		result[key]["sell_type"] = value.SellType
		result[key]["created_at"] = value.CreatedAt
		// extend
		result[key]["cas_no"] = extend["cas_no"]
		result[key]["purity"] = extend["purity"]
		result[key]["spec"] = extend["spec"]
		result[key]["img"] = extend["img"]
		result[key]["goods_name"] = extend["goods_name"]
		result[key]["goods_name_en"] = extend["goods_name_en"]
		result[key]["brand_id"] = extend["brand_id"]
		result[key]["brand_name"] = extend["brand_name"]
		result[key]["brand_name_en"] = extend["brand_name_en"]
		result[key]["sale_price"] = extend["sale_price"]
		result[key]["sale_unit"] = extend["sale_unit"]
		result[key]["structImage"] = ""
		result[key]["is_valid"] = 0 // 1 => 上架 0 => 下架

		if value.StoreId > 0 && value.SkuId > 0 {
			query[key] = map[string]int{
				"storeId": value.StoreId,
				"skuId":   value.SkuId,
			}
		}
	}

	invalidGoods := GetStoreSkuStatus(query)

	var invalids = make(map[string]interface{})
	for _, value := range invalidGoods {
		val := value.(map[string]interface{})
		if utils.Required(val["skuId"].(string)) {
			invalids[val["skuId"].(string)] = value
		}
	}

	for key, value := range result {
		for k, _ := range invalids {
			if k == value["sku_id"] {
				invalidStat := invalids[k].(map[string]interface{})
				status := invalidStat["status"]

				result[key]["is_valid"] = status // status-> 0上架 1下架
				result[key]["sale_price"] = invalidStat["price"]
			}
		}
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
	goodsId := c.Param("goodsId")
	userId := c.PostForm("user_id")
	storeId := c.PostForm("store_id")
	casNo := c.PostForm("cas_no")
	purity := c.PostForm("purity")
	spec := c.PostForm("spec")
	img := c.PostForm("img")
	skuId := c.DefaultPostForm("sku_id", "0")
	brandId := c.DefaultPostForm("brand_id", "0")
	brandName := c.DefaultPostForm("brand_name", "")
	price := c.DefaultPostForm("price", "")
	priceUnit := c.DefaultPostForm("price_unit", "")
	sellType := c.DefaultPostForm("sell_type", "")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(price) || !utils.Required(priceUnit) {
		c.JSON(200, gin.H{"code": 500, "msg": "价格信息不存在"})
		return
	}

	if !utils.Required(sellType) || sellType != "1" || sellType != "2" {
		c.JSON(200, gin.H{"code": 500, "msg": "请填写正确的商品类型"})
		return
	}

	if sellType == "1" && !utils.Required(storeId) {
		c.JSON(200, gin.H{"code": 500, "msg": "联营商品店铺不存在"})
		return
	}

	if sellType == "1" && !utils.Required(skuId) {
		c.JSON(200, gin.H{"code": 500, "msg": "商品 skuId 不存在"})
		return
	}

	if _, err := model.FindOne(storeId, userId); err == nil {
		c.JSON(200, gin.H{"code": 200, "msg": "success"})
		return
	}

	var extend = make(map[string]interface{})
	var storeName string
	if sellType == "1" { // 联营商品
		goodsInfo := GetPoolProduct(goodsId, storeId)
		if goodsInfo == nil {
			c.JSON(200, gin.H{"code": 500, "msg": "获取联营商品信息出错"})
			return
		}

		var saleUnit string
		if goodsInfo["sellType"] == 1 {
			saleUnit = goodsInfo["numUnitTxt"].(string) + "/" + goodsInfo["packUnitTxt"].(string)
		} else {
			saleUnit = goodsInfo["numUnitTxt"].(string)
		}

		extend = map[string]interface{}{
			"cas_no":        casNo,
			"purity":        purity,
			"spec":          spec,
			"img":           img,
			"goods_name":    goodsInfo["productZhName"],
			"goods_name_en": goodsInfo["productEnName"],
			"brand_id":      goodsInfo["brandId"],
			"brand_name":    goodsInfo["brandZhName"],
			"brand_name_en": goodsInfo["brandEnName"],
			"sale_price":    goodsInfo["salePrice"],
			"sale_unit":     saleUnit,
		}

		storeInfo := GetStoreInfo(storeId)
		if storeInfo == nil {
			c.JSON(200, gin.H{"code": 500, "msg": "获联营取店铺信息出错"})
			return
		}
		storeName = storeInfo["store_name"].(string)
	} else {
		goodsInfo := GetProprietaryProduct(goodsId)
		if goodsInfo == nil {
			c.JSON(200, gin.H{"code": 500, "msg": "获取自营商品信息出错"})
			return
		}

		brand := goodsInfo["brand"].(map[string]interface{})
		brandId := brand["id"]
		brandName := brand["brand_name"]
		brandNameEn := brand["brand_name_en"]

		extend = map[string]interface{}{
			"cas_no":        casNo,
			"purity":        purity,
			"spec":          spec,
			"img":           img,
			"goods_name":    goodsInfo["name"],
			"goods_name_en": goodsInfo["name_en"],
			"sale_price":    goodsInfo["sproduct_price"],
			"sale_unit":     goodsInfo["sell_unit"],
			"brand_id":      brandId,
			"brand_name":    brandName,
			"brand_name_en": brandNameEn,
		}

		// 自营商品根据商品id（sku_id）获取店铺信息
		storeInfos := GetProprietaryStoreInfo(goodsId)
		if storeInfos == nil {
			c.JSON(200, gin.H{"code": 500, "msg": "获取自营店铺信息出错"})
			return
		}

		storeInfo := storeInfos[0].(map[string]interface{})
		storeId = storeInfo["shopId"].(string)
		storeName = storeInfo["shopName"].(string)

		// 2、根据 自营商品id（sku_id）+店铺id 获取自营商品对应联营商品信息（goods_id,sku_id）
		if utils.StringToInt(storeId) > 0 {
			poolGoodsInfo := GetPoolProductBySearch(goodsId, storeId)
			if poolGoodsInfo == nil {
				c.JSON(200, gin.H{"code": 500, "msg": "获取联营商品信息出错"})
				return
			}

			info := poolGoodsInfo["data"].(map[int]interface{})[0].(map[string]string)
			//poolGoodsId := info["goodsId"]
			skuId = info["skuId"]
		}
	}

	params := map[string]interface{}{
		"goodsId":   utils.StringToInt(goodsId),
		"skuId":     utils.StringToInt(skuId),
		"userId":    utils.StringToInt(userId),
		"brandId":   utils.StringToInt(brandId),
		"brandName": brandName,
		"price":     price,
		"priceUnit": priceUnit,
		"sellType":  sellType,
		"storeId":   utils.StringToInt(storeId),
		"storeName": storeName,
		"extend":    utils.MapToJson(extend),
	}
	model.Create(params)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func Cancel(c *gin.Context) {
	goodsId := c.Param("goodsId")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	model.BatchCancel(userId, goodsId)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func BatchCancel(c *gin.Context) {
	goodsIds := c.PostForm("goods_ids")
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	model.BatchCancel(userId, goodsIds)

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}

func Brands(c *gin.Context) {
	userId := c.PostForm("user_id")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	data := model.FindBrandsByUserId(userId)

	c.JSON(200, gin.H{"code": 200, "data": data})
	return
}

func EmptyGoods(c *gin.Context) {
	userId := c.PostForm("user_id")
	_type := c.PostForm("type")

	if !IsValidMember(userId) {
		c.JSON(200, gin.H{"code": 500, "msg": "会员不存在"})
		return
	}

	if !utils.Required(_type) {
		c.JSON(200, gin.H{"code": 500, "msg": "操作类型不能为空"})
		return
	}

	params := map[string]interface{}{
		"user_id":userId,
	}
	data,_ := model.Find(params)

	var query = make(map[int]map[string]int)
	for key, value := range data {
		if value.StoreId == 0 && value.SkuId == 0 {	// 非法数据，视为无效商品
			continue
		}

		query[key] = map[string]int{
			"storeId": value.StoreId,
			"skuId": value.SkuId,
		}
	}

	if _type == "invalid" {
		goodsStatus := GetStoreSkuStatus(query)

		var querySku map[int]int
		for key, value := range query {
			querySku[key] = value["skuId"]
		}

		var goodsSku []int
		var invalids []int
		for key, value := range goodsStatus {
			val := value.(map[string]interface{})
			if utils.Required(val["skuId"].(string)) {
				goodsSku[key] = val["skuId"].(int)
				if val["status"].(int) == 0 {
					invalids[key] = val["skuId"].(int)	// status-> 0上架 1下架
				}
			}
		}

		for key, value := range querySku {
			for _, val := range goodsSku {
				if val == value {
					delete(querySku, key)	// 删除查询到的商品，没有查询到结果的 skuId 即为失效商品
				}
			}
		}

		for _, value := range querySku {
			invalids = append(invalids, value)
		}

		model.DeleteInvalidGoods(invalids, userId)
	}

	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}
