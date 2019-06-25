package config

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"auroraLZDF/member_api/utils"
)

/*const (
	TimeFormat = "2006-01-02 15:04:05"
)*/

type config struct {
	AppPath string `json:"app_path"`

	DbHost string `json:"db_host"`
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbPort string `json:"db_port"`
	DbName string `json:"db_name"`

	Env             string `json:"env";default:"production"`
	MpsApiUrl       string `json:"mps_api_url";default:"http://mps.molbase.org"`
	PmsApiUrl       string `json:"pms_api_url";default:"http://pms.molbase.org"`
	PmsApiAccessKey string `json:"pms_api_access_key";default:"molbase-erp-!@#$%^&*()-66666666666666666"`
	PmsApiFromCrm   string `json:"pms_api_from_crm";default:"molbase_cms"`
	MonitorApiUrl   string `json:"monitor_api_url";default:"http://monitor.cron.molbase.org"`
	WapApiUrl       string `json:"wap_api_url";default:"http://m.molbase.com"`
	SpApiUrl        string `json:"sp_api_url";default:"http://sp.molbase.org"`
	HhwApiUrl       string `json:"hhw_api_url";default:"http://hhw.molbase.org"`
	ErpApiUrl       string `json:"erp_api_url";default:"http://erp.molbase.org"`
	SearchApiUrl    string `json:"search_api_url";default:"http://search.molbase.org"`
}

var Config config

func SetConfig() {
	path := utils.AppPath() + ".env.json"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic("读取配置文件失败")
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("[loadConfig]: %s\n", err.Error())
	}
}
