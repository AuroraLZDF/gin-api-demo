package Models

import (
	"auroraLZDF/member_api/utils"
	"auroraLZDF/member_api/config"
	"log"
	"strings"
	"database/sql"
)

type Goods struct {
	Id         int    `gorm:"unsigned;primary_key"`
	UserId     int    `gorm:"type:int(11);unsigned;not null"`
	StoreId    int    `gorm:"type:int(11);not null;default:0"`
	StoreName  string    `gorm:"type:int(11);not null;default:''"`
	GoodsId   int    `gorm:"type:int(11);unsigned;not null"`
	SkuId     int    `gorm:"type:int(11);not null;default:0"`
	BrandId   int    `gorm:"type:int(11);default:0"`
	BrandName string `gorm:"type:varchar(250);not null;default:''"`
	Price     string `gorm:"type:varchar(200);not null;default:''"`
	PriceUnit string `gorm:"type:varchar(200);not null;default:''"`
	SellType  string `gorm:"type:tinyint(2);unsigned;not null;default:0"`
	Extend    string `gorm:"type:text;not null;default:''"`
	CreatedAt utils.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Goods) TableName() string {
	return "_collection_goods"
}

func (m Goods) Create(params map[string]interface{}) bool {
	m.GoodsId = params["goodsId"].(int)
	m.SkuId = params["skuId"].(int)
	m.UserId = params["userId"].(int)
	m.BrandId = params["brandId"].(int)
	m.BrandName = params["brandName"].(string)
	m.Price = params["price"].(string)
	m.PriceUnit = params["priceUnit"].(string)
	m.SellType = params["sellType"].(string)
	m.StoreId = params["storeId"].(int)
	m.StoreName = params["storeName"].(string)
	m.Extend = params["extend"].(string)


	db := config.Db.Create(&m)
	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return false
	}

	return true
}

func (m Goods) Find(params map[string]interface{}) ([]Goods, int64) {
	var count int64
	var goods []Goods

	db := config.Db.Model(m)

	/*if params["page"] == 0 {
		params["page"] = 1
	}

	if params["limit"] == 0 {
		params["limit"] = 15
	}*/

	if params["user_id"].(int) > 0 {
		db = db.Where("user_id = ?", params["user_id"])
	}

	if params["goods_id"].(int) > 0 {
		db = db.Where("goods_id = ?", params["goods_id"])
	}

	if params["brand_id"].(int) > 0 {
		db = db.Where("brand_id = ?", params["brand_id"])
	}

	if utils.Required(params["begin_time"].(string)) {
		db = db.Where("created_at >= ?", params["begin_time"])
	}

	if utils.Required(params["end_time"].(string)) {
		db = db.Where("created_at <= ?", params["end_time"])
	}

	// 获取记录数
	db.Count(&count)

	// 获取记录数据
	db = db.Order("id desc").Limit(params["limit"]).Offset((params["page"].(int) - 1) * params["limit"].(int)).Find(&goods)

	return goods, count
}

func (m Goods) FindOne(goodsId string, userId string) (data Store, err error) {
	db := config.Db.Where("goods_id=?", goodsId).Where("user_id=?", userId).First(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return data, err
	}

	return data, nil
}

func (m Goods) FindBrandsByUserId(userId string) map[int]map[string]interface{} {
	var rows *sql.Rows
	var data = make(map[int]map[string]interface{})

	rows, _ = config.Db.Model(m).Select("brand_id,brand_name,COUNT(id) AS total").Where("user_id = ?", userId).Group("brand_id").Rows()

	var brandId int
	var brandName string
	var total int
	var count = 0

	for rows.Next() {
		_ = rows.Scan(&brandId, &brandName, &total)

		data[count] = map[string]interface{}{
			"brand_id": brandId,
			"brand_name": brandName,
			"total":      total,
		}
		count++
	}

	return data
}

func (m Goods) BatchCancel(userId string, goodsIds string) error {
	var goodsId []string
	goodsId = strings.Split(goodsIds, ",")

	db := config.Db.Where("user_id = ?", userId).Where("goods_id in (?)", goodsId).Delete(&m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}

func (m Goods) DeleteInvalidGoods(invalids []int, userId string) error {
	db := config.Db.Where("user_id = ?", userId).Where("sku_id in (?)", invalids).Delete(&m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}
