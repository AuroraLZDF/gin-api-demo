package Models

import (
	"auroraLZDF/member_api/config"
	"auroraLZDF/member_api/utils"
	"log"
	"strings"
	"database/sql"
)

type MolData struct {
	Id           int        `gorm:"type:int(11);unsigned;primary_key"`
	UserId       int        `gorm:"type:int(11);unsigned;not null"`
	MolId        int        `gorm:"type:int(11);not null;default:0"`
	CategoryId   int        `gorm:"type:int(11);not null;default:''"`
	CategoryName string     `gorm:"varchar(255);not null;default:''"`
	Extend       string     `gorm:"type:text;not null;default:''"`
	CreatedAt    utils.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	NameZh       string
	NameEn       string
	CasNo        string
	Img          string
}

func (MolData) TableName() string {
	return "_collection_moldata"
}

func (m MolData) Create(molId int, userId int, extend string) bool {
	m.MolId = molId
	m.UserId = userId
	m.Extend = extend

	db := config.Db.Create(&m)
	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return false
	}

	return true
}

func (m MolData) Find(params map[string]interface{}) ([]MolData, int64) {
	var count int64
	var molData []MolData

	db := config.Db.Model(m)

	if params["user_id"].(int) > 0 {
		db = db.Where("user_id = ?", params["user_id"])
	}

	if params["mol_id"].(int) > 0 {
		db = db.Where("mol_id = ?", params["mol_id"])
	}

	if params["category_id"].(int) > 0 {
		db = db.Where("category_id = ?", params["category_id"])
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
	if params["limit"] == 0 {
		db = db.Order("id desc").Find(&molData)
	} else {
		db = db.Order("id desc").Limit(params["limit"]).Offset((params["page"].(int) - 1) * params["limit"].(int)).Find(&molData)
	}

	return molData, count
}

func (m MolData) FindOne(molId string, userId string) (data MolData, err error) {
	db := config.Db.Where("mol_id=?", molId).Where("user_id=?", userId).First(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return data, err
	}

	return data, nil
}

func (m MolData) BatchCancel(userId string, molIds string) error {
	var molId []string
	molId = strings.Split(molIds, ",")

	db := config.Db.Where("user_id = ?", userId).Where("mol_id in (?)", molId).Delete(&m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}

func (m MolData) FindCategoriesByUserId(userId string) map[int]map[string]interface{} {
	var rows *sql.Rows
	var data = make(map[int]map[string]interface{})

	rows, _ = config.Db.Model(m).Select("category_id,category_name,COUNT(id) AS total").Where("user_id = ?", userId).Group("category_id").Rows()

	var categoryId int
	var categoryName string
	var total int
	var count = 0

	for rows.Next() {
		_ = rows.Scan(&categoryId, &categoryName, &total)

		data[count] = map[string]interface{}{
			"category_id": categoryId,
			"category_name": categoryName,
			"total":      total,
		}
		count++
	}

	return data
}

func (m MolData) HasCollection(userId string, molIds string) (map[int]map[string]interface{}, error) {
	var data []MolData
	var result = make(map[int]map[string]interface{})
	var molId []string
	molId = strings.Split(molIds, ",")

	db := config.Db.Where("user_id = ?", userId).Where("mol_id in (?)", molId).Find(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return nil, err
	}

	for _, value := range data {
		status := 2
		if utils.InArray(utils.IntToString(value.MolId), molId) {
			status = 1
		}

		result[value.MolId] = map[string]interface{}{
			"mol_id": value.MolId,
			"status": status,
		}
	}

	return result, nil
}


