package Models

import (
	"auroraLZDF/member_api/utils"
	"auroraLZDF/member_api/config"
	"log"
	"strings"
)

type Supply struct {
	Id           int        `gorm:"type:int(11);unsigned;primary_key"`
	UserId       int        `gorm:"type:int(11);unsigned;not null"`
	SupplyId     int        `gorm:"type:int(11);not null"`
	ProductName  string     `gorm:"type:varchar(200);default:null"`
	Purity       string     `gorm:"type:varchar(60);default:null"`
	SpecCount    string     `gorm:"type:varchar(45);default:0"`
	SpecUnit     string     `gorm:"type:varchar(15);default:null"`
	SpecPackage  string     `gorm:"type:varchar(15);default:null"`
	Price        string     `gorm:"type:varchar(45);default:null"`
	PriceUnit    string     `gorm:"type:varchar(45);default:null"`
	ProductLevel int        `gorm:"type:tinyint(4);not null;default:0"`
	ProvinceName string     `gorm:"type:varchar(20);default:null"`
	CityName     string     `gorm:"type:varchar(20);default:null"`
	Images       string     `gorm:"type:text"`
	Period       int     `gorm:"type:tinyint(4);not null;default:0"`
	CreationTime utils.Time `gorm:"type:timestamp;default:null"`
	CreatedAt    utils.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Supply) TableName() string {
	return "_collection_supply"
}

func (m Supply) Create(userId int, supplyId int, params map[string]interface{}) bool {
	m.UserId = userId
	m.SupplyId = supplyId
	m.ProductName = params["product_name"].(string)
	m.Purity = params["purity"].(string)
	m.SpecCount = params["spec_count"].(string)
	m.SpecUnit = params["spec_unit"].(string)
	m.SpecPackage = params["spec_package"].(string)
	m.Price = params["price"].(string)
	m.PriceUnit = params["price_unit"].(string)
	m.ProductLevel = params["product_level"].(int)
	m.ProvinceName = params["province_name"].(string)
	m.CityName = params["city_name"].(string)
	m.Images = params["images"].(string)
	m.Period = params["period"].(int)
	m.CreationTime = params["creation_time"].(utils.Time)

	db := config.Db.Create(&m)
	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return false
	}

	return true
}

func (m Supply) Find(params map[string]interface{}) ([]Supply, int64) {
	var count int64
	var purchase []Supply

	db := config.Db.Model(m)

	if params["user_id"].(int) > 0 {
		db = db.Where("user_id = ?", params["user_id"])
	}

	if utils.Required(params["province_name"].(string)) {
		db = db.Where("province_name = ?", params["province_name"])
	}

	if utils.Required(params["city_name"].(string)) {
		db = db.Where("city_name = ?", params["city_name"])
	}

	if utils.Required(params["spec_unit"].(string)) {
		db = db.Where("spec_unit = ?", params["spec_unit"])
	}

	// 获取记录数
	db.Count(&count)

	if params["sort"].(int) > 0 && params["sort"].(int) == 1 {
		db = db.Order("creation_time DESC")
	} else {
		db = db.Order("created_at DESC")
	}

	db = db.Limit(params["limit"]).Offset((params["page"].(int) - 1) * params["limit"].(int)).Find(&purchase)

	return purchase, count
}

func (m Supply) FindOne(userId string, supplyId string) (data Supply, err error) {
	db := config.Db.Where("user_id = ?", userId).Where("supply_id = ?", supplyId).First(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return data, err
	}

	return data, nil
}

func (m Supply) BatchCancel(userId string, supplyIds string) error {
	var supplyId []string
	supplyId = strings.Split(supplyIds, ",")

	db := config.Db.Where("user_id = ?", userId).Where("supply_id in (?)", supplyId).Delete(&m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}

func (m Supply) HasCollection(userId string, supplyIds string) (map[int]map[string]interface{}, error) {
	var data []Supply
	var result = make(map[int]map[string]interface{})
	var supplyId []string
	supplyId = strings.Split(supplyIds, ",")

	db := config.Db.Where("user_id = ?", userId).Where("supply_id in (?)", supplyId).Find(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return nil, err
	}

	for _, value := range data {
		status := 2
		if utils.InArray(utils.IntToString(value.SupplyId), supplyId) {
			status = 1
		}

		result[value.SupplyId] = map[string]interface{}{
			"supply_id": value.SupplyId,
			"status": status,
		}
	}

	return result, nil
}
