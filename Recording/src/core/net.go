package core

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"net/http"
)

func Post(url string, params map[string]interface{}) (int, *simplejson.Json, error) {

	bytesParams, err := json.Marshal(params)
	if err != nil {
		return 0, nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(bytesParams))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}

	//respBytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return 0, nil, err
	//}

	//byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println(*str)

	json, err := simplejson.NewFromReader(resp.Body)

	return resp.StatusCode, json, err
}
