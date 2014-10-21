package yahooweather

import (
	"net/http"
	"net/url"
	"go-simplejson-master"
	"io/ioutil"
	"fmt"
)

//Even in struct, only words that start with upper-case letter can be used outside the package
type WeatherInfo struct {
	Temp string
	Humidity string
	Weth string
	Units
}

type Units struct {
	Tp string
}

func BuildUrl(place string) (urlParsed string) {
	Url, _ := url.Parse("https://query.yahooapis.com/v1/public/yql")
	parameters := url.Values{}
	parameters.Add("q","select * from weather.forecast where woeid = " + place)
	parameters.Add("format","json")
	Url.RawQuery = parameters.Encode()
	urlParsed = Url.String()
	return 
}

func MakeQuery(weatherUrl string) (w *WeatherInfo){
	resp, err := http.Get(weatherUrl)
	if err != nil {
		fmt.Println("Connected Error")
		return nil
	}

	defer resp.Body.Close()
	body, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		fmt.Println("Cannot Read Information")
		return nil
	}

	js, e := simplejson.NewJson(body)
	if e != nil {
		fmt.Println("Parsing Json Error")
		return nil
	}

	//parse json
	w = new(WeatherInfo)	
	w.Tp, _ = js.Get("query").Get("results").Get("channel").Get("units").Get("temperature").String()
	w.Temp, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("temp").String()
	w.Weth,_ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("text").String()
	w.Humidity,_ = js.Get("query").Get("results").Get("channel").Get("atmosphere").Get("humidity").String()
	return 
}