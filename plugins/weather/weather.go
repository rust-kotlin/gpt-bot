package weather

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetWeather(location, apiKey string) (string, error) {
	locationID, locationName, err := getGeo(location, apiKey)
	if err != nil {
		//fmt.Println("获取地理位置信息失败：", err)
		return "", err
	}
	Url := fmt.Sprintf("https://devapi.qweather.com/v7/weather/3d?location=%s&key=%s", locationID, apiKey)
	// 发送HTTP请求
	resp, err := http.Get(Url)
	if err != nil {
		//fmt.Println("发送HTTP请求失败：", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭HTTP响应Body失败：", err)
		}
	}(resp.Body)
	var data map[string]interface{}
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		// 如果是gzip压缩，则使用gzip.Reader解压缩响应数据
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("解压缩失败：%v", err)
		}
		defer func(gzipReader *gzip.Reader) {
			err := gzipReader.Close()
			if err != nil {
				fmt.Println("关闭gzip.Reader失败：", err)
			}
		}(gzipReader)
		jsonData, err := ioutil.ReadAll(gzipReader)
		if err != nil {
			return "", fmt.Errorf("读取API响应失败：%v", err)
		}
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return "", fmt.Errorf("解析JSON失败：%v", err)
		}
	} else {
		// 如果响应未经过gzip压缩，则直接读取并解析JSON
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", fmt.Errorf("解析JSON失败：%v", err)
		}
	}
	// 获取天气信息
	if data["code"].(string) == "200" {
		daily := data["daily"].([]interface{})
		s := locationName + " 的天气预报：\n"
		for _, item := range daily {
			//fmt.Println(i)
			day := item.(map[string]interface{})
			fxDate := day["fxDate"].(string)
			sunrise := day["sunrise"].(string)
			sunset := day["sunset"].(string)
			tempMax := day["tempMax"].(string)
			tempMin := day["tempMin"].(string)
			textDay := day["textDay"].(string)
			textNight := day["textNight"].(string)
			s += fmt.Sprintf("%s，日出：%s，日落：%s，最高温度：%s℃，最低温度：%s℃，白天天气状况：%s，夜间天气状况：%s\n", fxDate, sunrise, sunset, tempMax, tempMin, textDay, textNight)
		}
		return s, nil
	} else {
		return "", fmt.Errorf("获取天气失败：%s", location)
	}
}

func getGeo(location, apiKey string) (string, string, error) {
	// 构建请求URL
	Url := fmt.Sprintf("https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s", url.QueryEscape(location), apiKey)
	// 发送HTTP GET请求
	resp, err := http.Get(Url)
	if err != nil {
		return "", "", fmt.Errorf("发送HTTP请求失败：%v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭HTTP响应Body失败：", err)
		}
	}(resp.Body)
	// 读取响应数据并解析JSON
	var result map[string]interface{}
	// 检查响应是否为gzip压缩
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		// 如果是gzip压缩，则使用gzip.Reader解压缩响应数据
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("解压缩失败：%v", err)
		}
		defer func(gzipReader *gzip.Reader) {
			err := gzipReader.Close()
			if err != nil {
				fmt.Println("关闭gzip.Reader失败：", err)
			}
		}(gzipReader)
		jsonData, err := ioutil.ReadAll(gzipReader)
		if err != nil {
			return "", "", fmt.Errorf("读取API响应失败：%v", err)
		}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			return "", "", fmt.Errorf("解析JSON失败：%v", err)
		}
	} else {
		// 如果响应未经过gzip压缩，则直接读取并解析JSON
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return "", "", fmt.Errorf("解析JSON失败：%v", err)
		}
	}
	if result["code"] == "404" {
		return "", "", fmt.Errorf("未找到城市：%s", location)
	}
	locations := result["location"].([]interface{})
	var cityList []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		Adm1 string `json:"adm1"`
		Adm2 string `json:"adm2"`
	}
	for _, loc := range locations {
		l := loc.(map[string]interface{})
		var city struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Adm1 string `json:"adm1"`
			Adm2 string `json:"adm2"`
		}
		city.Name = l["name"].(string)
		city.ID = l["id"].(string)
		city.Adm1 = l["adm1"].(string)
		city.Adm2 = l["adm2"].(string)
		cityList = append(cityList, city)
	}
	// 返回城市ID
	return cityList[0].ID, cityList[0].Adm1 + cityList[0].Adm2 + " " + cityList[0].Name, nil
}
