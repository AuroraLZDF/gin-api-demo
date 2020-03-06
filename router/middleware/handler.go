package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/auroraLZDF/gin-api-demo/utils"
	"log"
	"fmt"
)

var secretKey = "d1f3472bae4aed8a6736a3cef01fd136"

func Handler() gin.HandlerFunc {

	return func(c *gin.Context) {
		var data = make(map[string]interface{})

		//c.Request.ParseForm()
		for key, value := range c.Request.Form {
			fmt.Println("key", key, "value", value)
			data[key] = value
		}

		fmt.Println(data)

		if !utils.Required(data["sign"].(string)) {
			c.String(500, "Undefined index: sign.")
			return
		}

		sign := data["sign"]
		delete(data, "sign")
		data["secretKey"] = secretKey

		var prefix = "?"
		for key, value := range data {
			prefix += "&" + key + "=" + value.(string)
		}

		if sign != utils.Md5(prefix + secretKey) {
			c.String(401, "Unauthorized.")
			return
		}

		log.Println("sign: " + sign.(string), "prefix: " + prefix)

		c.Next()
	}
}