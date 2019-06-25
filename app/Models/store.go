package Models

import (
	"auroraLZDF/member_api/utils"
	"auroraLZDF/member_api/config"
	"log"
	"database/sql"
	"strings"
)

type Store struct {
	Id        int        `gorm:"primary_key"`
	UserId    int        `gorm:"type:int(11);unsigned;not null"`
	StoreId   int        `gorm:"type:int(11);unsigned;not null"`
	StoreType int        `gorm:"type:int(11);unsigned;not null;default:0"`
	Extend    string     `gorm:"type:text;not null"`
	CreatedAt utils.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Store) TableName() string {
	return "_collection_store"
}

func (m Store) Create(storeId string, userId string, storeType string, extend map[string]interface{}) bool {
	m.UserId = utils.StringToInt(userId)
	m.StoreId = utils.StringToInt(storeId)
	m.StoreType = utils.StringToInt(storeType)
	m.Extend = utils.MapToJson(extend)

	db := config.Db.Create(&m)
	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return false
	}

	return true
}

func (m Store) Find(params map[string]int) ([]Store, int64) {
	var count int64
	var stores []Store

	db := config.Db.Model(m)

	if params["page"] == 0 {
		params["page"] = 1
	}

	if params["limit"] == 0 {
		params["limit"] = 15
	}

	if params["user_id"] > 0 {
		db = db.Where("user_id = ?", params["user_id"])
	}

	if params["store_id"] > 0 {
		db = db.Where("store_id = ?", params["store_id"])
	}

	if params["store_type"] > 0 {
		db = db.Where("store_type = ?", params["store_type"])
	}

	// 获取记录数
	db.Count(&count)

	// 获取记录数据
	db = db.Order("id desc").Limit(params["limit"]).Offset((params["page"] - 1) * params["limit"]).Find(&stores)

	return stores, count
}

func (m Store) FindOne(storeId string, userId string) (data Store, err error) {
	db := config.Db.Where("store_id=?", storeId).Where("user_id=?", userId).First(&data)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return data, err
	}

	return data, nil
}

func (m Store) FindStoreTypeByUserId(userId string) (map[int]map[string]int) {
	var rows *sql.Rows
	var data = make(map[int]map[string]int)

	rows, _ = config.Db.Model(m).Select("store_type,COUNT(id) AS total").Where("user_id = ?", userId).Group("store_type").Rows()

	var _type int
	var total int
	var count = 0

	for rows.Next() {
		_ = rows.Scan(&_type, &total)

		data[count] = map[string]int{
			"store_type": _type,
			"total":      total,
		}
		count++
	}

	return data
}

/*func (m Store) Cancel(userId int, storeId int) error {
	db := config.Db.Where("user_id = ?", userId).Where("store_id = ?", storeId).Delete(m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}*/

func (m Store) BatchCancel(userId string, storeIds string) error {
	var storeId []string
	storeId = strings.Split(storeIds, ",")

	db := config.Db.Where("user_id = ?", userId).Where("store_id in (?)", storeId).Delete(&m)

	if err := db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		return err
	}

	return nil
}
