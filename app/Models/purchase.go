package Models

import (
	"auroraLZDF/member_api/utils"
	"auroraLZDF/member_api/config"
	"log"
	"strings"
)

type Purchase struct {
	Id           int        `gorm:"type:int(11);unsigned;primary_key"`
	UserId       int        `gorm:"type:int(11);unsigned;not null"`
	Code         string     `gorm:"type:varchar(20);not null"`
	ProductName  string     `gorm:"type:varchar(200);default:null"`
	Cas          string     `gorm:"type:varchar(45);default:null"`
	Num          int        `gorm:"type:int(11);default:null"`
	NumUnit      int        `gorm:"type:tinyint(4);default:null"`
	Purity       string     `gorm:"type:varchar(60);default:null"`
	ProvinceName string     `gorm:"type:varchar(20);default:null"`
	CityName     string     `gorm:"type:varchar(20);default:null"`
	Remarks      string     `gorm:"type:text"`
	State        int        `gorm:"type:tinyint(4);default:null"`
	CreationTime utils.Time `gorm:"type:timestamp;default:null"`
	CreatedAt    utils.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Number       int
	NumberUnit   int
}

func (Purchase) TableName() string {
	return "_collection_purchase"
}

func (m Purchase) Create(userId int, code string, params map[string]interface{}) bool {
	m.UserId = userId
	m.Code = code
	m.ProductName = params["product_name"].(string)
	m.Cas = params["cas"].(string)
	m.Num = params["number"].(int)
	m.NumUnit = params["number_unit"].(int)
	m.Purity = params["purity"].(string)
	m.ProvinceName = params["province_name"].(string)
	m.CityName = params["city_name"].(string)
	m.Remarks = params["remarks"].(string)
	m.State = params["state"].(int)
	m.CreationTime = params["creation_time"].(utils.Time)

	db := config.Db.Create(&m)
	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return false
	}

	return true
}

func (m Purchase) Find(params map[string]interface{}) ([]Purchase, int64) {
	var count int64
	var purchase []Purchase

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

	if params["number_unit"].(int) > 0 {
		db = db.Where("number_unit = ?", params["number_unit"])
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

func (m Purchase) FindOne(userId string, code string) (data Purchase, err error) {
	db := config.Db.Where("user_id = ?", userId).Where("code = ?", code).First(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return data, err
	}

	return data, nil
}

func (m Purchase) BatchCancel(userId string, codes string) error {
	var code []string
	code = strings.Split(codes, ",")

	db := config.Db.Where("user_id = ?", userId).Where("code in (?)", code).Delete(&m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}

func (m Purchase) HasCollection(userId string, codes string) (map[string]map[string]interface{}, error) {
	var data []Purchase
	var result = make(map[string]map[string]interface{})
	var code []string
	code = strings.Split(codes, ",")

	db := config.Db.Where("user_id = ?", userId).Where("code in (?)", code).Find(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return nil, err
	}

	for _, value := range data {
		status := 2
		if utils.InArray(value.Code, code) {
			status = 1
		}

		result[value.Code] = map[string]interface{}{
			"mol_id": value.Code,
			"status": status,
		}
	}

	return result, nil
}
