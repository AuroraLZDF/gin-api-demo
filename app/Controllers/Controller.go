package Controllers

import (
	. "github.com/auroraLZDF/gin-api-demo/app/Models"
	. "github.com/auroraLZDF/gin-api-demo/config"
	"github.com/auroraLZDF/gin-api-demo/utils/Requests"
	"github.com/auroraLZDF/gin-api-demo/utils"
	"log"
	"strings"
	"encoding/json"
)

func IsValidMember(userId string) bool {
	if !utils.Required(userId) {
		return false
	}

	_, err := GetMemberInfoByFields("user_id", userId)
	if err != nil {
		return false
	}

	return true
}

func GetStoreInfo(storeId string) map[string]interface{} {
	url := Config.ErpApiUrl + "/api/store/get_store_data?id=" + storeId + "&company_info=true"

	var headers map[string]string
	response, _ := Requests.Get(url, headers)
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["status_code"].(float64)) != 1 {
		log.Println("获取店铺信息出错：" + responseData["message"].(string))
		return nil
	}

	data := responseData["data"].(map[string]interface{})

	return data
}

func GetStoreVip(storeIds []string) map[string]interface{} {
	url := Config.ErpApiUrl + "/api/store/vip_check"

	ids := strings.Join(storeIds, ",")

	var headers map[string]string
	response, _ := Requests.Get(url+"?store_ids="+ids, headers)
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["status_code"].(float64)) != 1 {
		log.Println("获取店铺 VIP 信息出错：" + responseData["message"].(string))
		return nil
	}

	data := responseData["data"].(map[string]interface{})

	return data
}

func GetStoreSkuStatus(requestData map[int]map[string]int) map[int]interface{} {
	url := Config.SpApiUrl + "/api/store/product/getStoreSkuStatus"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	jsonBytes, _ := json.Marshal(requestData)
	response, _ := Requests.Post(url, headers, string(jsonBytes))
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["code"].(float64)) != 0 {
		log.Println("获取商品状态出错：" + responseData["msg"].(string))
		return nil
	}

	data := responseData["data"].(map[int]interface{})

	return data
}

func GetPoolProduct(goodsId string, storeId string) map[string]interface{} {
	url := Config.SpApiUrl + "/api/store/product/getStoreProduct"

	var headers map[string]string
	response, _ := Requests.Get(url+"?storeId="+storeId+"&productId="+goodsId, headers)
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["code"].(float64)) != 0 {
		log.Println("获取联营商品出错：：" + responseData["msg"].(string))
		return nil
	}

	data := responseData["data"].(map[string]interface{})

	return data
}

func GetProprietaryProduct(goodsId string) map[string]interface{} {
	url := Config.PmsApiUrl + "/api/sproduct/info"

	var headers map[string]string
	requestData := "id="+goodsId
	response, _ := Requests.Post(url, headers, requestData)
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["code"].(float64)) != 0 {
		log.Println("获取自营商品出错：" + responseData["msg"].(string))
		return nil
	}

	data := responseData["data"].(map[string]interface{})

	return data
}

func GetProprietaryStoreInfo(goodsId string) map[int]interface{} {
	url := Config.PmsApiUrl + "/api/goods/getShopInfoBySku"

	var headers map[string]string
	response, _ := Requests.Get(url+"?ids="+goodsId, headers)
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["code"].(float64)) != 200 {
		log.Println("获取自营商品店铺信息出错：" + responseData["msg"].(string))
		return nil
	}

	data := responseData["data"].(map[int]interface{})

	return data
}

func GetPoolProductBySearch(goodsId string, storeId string) map[string]interface{} {
	url := Config.SearchApiUrl + "/pms/getProductByStore"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	var goodsIds = []string{goodsId}
	//requestData := "pageSize=100&order=1&fromMall=true&storeSelf=true&storeId="+storeId+"&pmsSkuIds="+
	requestData := map[string]interface{}{
		"pageSize": 100,
		"order": 1,
		"fromMall": true,
		"storeSelf": true,
		"storeId": storeId,
		"pmsSkuIds": goodsIds,
	}

	jsonBytes, _ := json.Marshal(requestData)
	response, _ := Requests.Post(url, headers, string(jsonBytes))
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["code"].(float64)) != 0 {
		log.Println("根据自营商品获取联营商品状态出错：" + responseData["msg"].(string))
		return nil
	}

	data := responseData["data"].(map[string]interface{})

	return data


}

func GetMolInfo(molId string) map[string]interface{} {
	url := Config.HhwApiUrl + "/api/compound/findByCompoundId"

	var headers map[string]string
	response, _ := Requests.Get(url+"?id="+molId, headers)
	responseData := utils.JsonToMap(response)

	if utils.FloatToInt(responseData["code"].(float64)) != 0 {
		log.Println("获取百科详情出错：" + responseData["msg"].(string))
		return nil
	}

	data := responseData["data"].(map[string]interface{})

	return data
}
