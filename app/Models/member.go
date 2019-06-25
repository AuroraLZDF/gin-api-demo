package Models

import (
	//"auroraLZDF/member_api/config"
	time "auroraLZDF/member_api/utils"
	"log"
	"auroraLZDF/member_api/config"
)

type Member struct {
	UserId        int       `gorm:"primary_key"`
	UserName      string    `gorm:"type:varchar(100);not null;default:''"`
	Email         string    `gorm:"type:varchar(100);not null;default:''"`
	Mobile        string    `gorm:"type:varchar(60);not null"`
	Password      string    `gorm:"type:varchar(32);not null;default:''"` // json:"-" 不显示该字段;  json:"omitempty"，字段值为空不显示
	Gender        int       `gorm:"type:tinyint(4);not null;default:0"`
	Birthday      time.Time `gorm:"type:date"`
	Nickname      string    `gorm:"type:varchar(200)"`
	MailVerify    int       `gorm:"type:tinyint(1);unsigned;not null;default:0"`
	MobileVerify  int       `gorm:"type:tinyint(1);unsigned;not null;default:0"`
	CountryId     int       `gorm:"type:int(11);not null;default:0"`
	ProvinceId    int       `gorm:"type:int(11);not null;default:0"`
	CityId        int       `gorm:"type:int(11);not null;default:0"`
	AllowLogin    int       `gorm:"type:tinyint(1);not null;default:1"`
	RegSource     int       `gorm:"type:tinyint(1);not null;default:1"`
	RememberToken string    `gorm:"type:varchar(200);not null"`
	WxUnionid     string    `gorm:"type:varchar(100);not null;default:''"`
	RegTime       string    `gorm:"type:timestamp;not null;default:'0000-00-00 00:00:00'"`
	RegIp         string    `gorm:"type:varchar(15);not null"`
	RegAddr       string    `gorm:"type:varchar(30);not null;default:''"`
	LastLoginAddr string    `gorm:"type:varchar(30);not null;default:''"`
	LastIp        string    `gorm:"type:varchar(15);not null"`
}


func (Member) TableName() string {
	return "_member"
}

func GetMemberInfoByFields(field string, value string) (member Member, err error) {
	db := config.Db.Where(field + "=?", value).First(&member)
	if err = db.Error; err != nil {
		log.Printf("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
		log.Println("获取用户数据失败： ", err)
	}

	return
}