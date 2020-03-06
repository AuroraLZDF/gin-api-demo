package config

import (
	"encoding/json"
	"github.com/auroraLZDF/gin-api-demo/utils"
	"io/ioutil"
	"log"
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
	MpsApiUrl       string `json:"mps_api_url";default:"http://mps.xxxx.org"`
	PmsApiUrl       string `json:"pms_api_url";default:"http://pms.xxx.org"`
	PmsApiAccessKey string `json:"pms_api_access_key";default:"xxx-erp-!@#$%^&*()-66666666666666666"`
	PmsApiFromCrm   string `json:"pms_api_from_crm";default:"xxx"`
	MonitorApiUrl   string `json:"monitor_api_url";default:"http://monitor.cron.xxx.org"`
	WapApiUrl       string `json:"wap_api_url";default:"http://m.xxx.com"`
	SpApiUrl        string `json:"sp_api_url";default:"http://sp.xxx.org"`
	HhwApiUrl       string `json:"hhw_api_url";default:"http://hhw.xxx.org"`
	ErpApiUrl       string `json:"erp_api_url";default:"http://erp.xxx.org"`
	SearchApiUrl    string `json:"search_api_url";default:"http://search.xxx.org"`
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
