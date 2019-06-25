package Requests

import (
	"net/http"
	"strings"
	"log"
	"io/ioutil"
)

func Get(url string, headers map[string]string) (string, error) {
	return request("get", url, headers, "")
}

func Post(url string, headers map[string]string, requestData string) (string, error) {
	return request("post", url, headers, requestData)
}

func Put() {

}

func Delete() {

}

func Patch() {

}

func request(method string, url string, headers map[string]string, requestData string) (result string, err error) {
	var request *http.Request

	client := &http.Client{}

	method = strings.ToUpper(method)
	if method == "GET" {
		request, err = http.NewRequest(method, url, nil)
	} else if method == "POST" {
		request, err = http.NewRequest(method, url, strings.NewReader(requestData))
	}

	if err != nil {
		log.Println("http NewRequest err: ", err)
		return
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("http Do err: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println("http response err: ", err)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	result = string(body)
	//fmt.Println(bodystr)

	return
}
