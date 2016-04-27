package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	url_province = "http://www.weather.com.cn/data/list3/city.xml?level=1"
	url_city     = "http://www.weather.com.cn/data/list3/city%s.xml?level=2"
	url_county   = "http://www.weather.com.cn/data/list3/city%s.xml?level=3"
	url_weather  = "http://www.weather.com.cn/adat/cityinfo/%s.html"
	code_nation  = "101" // 国家代码
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	provincesStr, err := handler(url_province)
	if err != nil {
		log.Fatalf("加载省份数据出错:%v", err)
	}

	provinces := strings.Split(provincesStr, ",")
	for _, province := range provinces {
		province_code_name := strings.Split(province, "|")
		fmt.Println("province", "code:", province_code_name[0], "name:", province_code_name[1])
		// 城市
		citysStr, err := handler(fmt.Sprintf(url_city, province_code_name[0]))
		if err != nil {
			log.Fatalf("加载城市数据出错:%v", err)
		}

		citys := strings.Split(citysStr, ",")
		for _, city := range citys {
			city_code_name := strings.Split(city, "|")
			fmt.Println("  ===>city", "code:", city_code_name[0], "name:", city_code_name[1])

			// 县
			countysStr, err := handler(fmt.Sprintf(url_county, city_code_name[0]))
			if err != nil {
				log.Fatalf("加载县城数据出错:%v", err)
			}
			countys := strings.Split(countysStr, ",")
			for _, county := range countys {
				county_code_name := strings.Split(county, "|")
				fmt.Println("  === ===>county", "code:", county_code_name[0], "name:", county_code_name[1], code_nation+county_code_name[0])

				weather, err := handler(fmt.Sprintf(url_weather, code_nation+county_code_name[0]))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("=========> weather", weather)
			}
		}
	}
}

func handler(url string) (data string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	data = string(body[:])
	return
}
