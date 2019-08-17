/*
   Yahooweather is my first program in Go. It is still under development.
   Dependencies:
    github.com/bitly/go-simplejson

   Author: psaux0
   Version: 0.01
   Last Revised Time: Oct 21, 2014
*/
package yahooweather

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bitly/go-simplejson"
)

//WeatherInfo is returned with weather data
type WeatherInfo struct {
	Temp     string
	Humidity string
	Weth     string
	Units
}

//Units selects a unit
type Units struct {
	Tp string
}

//Location specifies a location, you need EITHER a zipcode or a city/state
type Location struct {
	City    string
	State   string
	Zipcode string
}

//BuildLocation builds a Location struct
func BuildLocation(city, state, zipcode string) (loc *Location) {
	return &Location{
		city,
		state,
		zipcode,
	}
}

//BuildUrl returns a parsed URL
func BuildUrl(loc *Location) (urlParsed string) {
	Url, _ := url.Parse("https://query.yahooapis.com/v1/public/yql")
	parameters := url.Values{}

	if loc.Zipcode != "" {
		parameters.Add("q", "select * from weather.forecast where woeid in (select woeid from geo.places(1) where text=\""+loc.Zipcode+"\")")

	} else {
		parameters.Add("q", "select * from weather.forecast where woeid in (select woeid from geo.places(1) where text=\""+loc.City+", "+loc.State+"\")")

	}

	parameters.Add("format", "json")
	Url.RawQuery = parameters.Encode()
	urlParsed = Url.String()
	return
}

//MakeQuery queries the servers
func MakeQuery(weatherUrl string) (w *WeatherInfo) {
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
	w.Weth, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("text").String()
	w.Humidity, _ = js.Get("query").Get("results").Get("channel").Get("atmosphere").Get("humidity").String()
	return
}
