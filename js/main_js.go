// +build js

package main

import (
	"fmt"

	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nwidger/wtw"
)

func getClothes() bool {
	go func() {
		var (
			tempInt int
			err     error
		)

		alert := func(message string) {
			js.Global.Call("alert", message)
		}

		document := js.Global.Get("document")

		locationElem := document.Call("getElementById", "location")
		location := locationElem.Get("value").String()

		if len(location) == 0 {
			alert(fmt.Sprintf("Please enter your location"))
			return
		}

		genderElem := document.Call("getElementById", "gender")
		gender := genderElem.Get("value").String()

		tempElem := document.Call("getElementById", "temp")
		temp := tempElem.Get("value").String()

		conditionsElem := document.Call("getElementById", "conditions")
		conditions := conditionsElem.Get("value").String()

		windElem := document.Call("getElementById", "wind")
		wind := windElem.Get("value").String()

		timeElem := document.Call("getElementById", "time")
		t := timeElem.Get("value").String()

		intensityElem := document.Call("getElementById", "intensity")
		intensity := intensityElem.Get("value").String()

		feelElem := document.Call("getElementById", "feel")
		feel := feelElem.Get("value").String()

		if len(temp) > 0 {
			tempInt, err = strconv.Atoi(temp)
			if err != nil {
				alert(fmt.Sprintf("Invalid temperature: %s", err.Error()))
				return
			}
		}

		a := &wtw.Conditions{
			Gender:     gender,
			Temp:       tempInt,
			Conditions: conditions,
			Wind:       wind,
			Time:       t,
			Intensity:  intensity,
			Feel:       feel,
		}

		if a.Time == "current" {
			a.Time = wtw.GetTime()
		}

		if len(location) > 0 {
			err = wtw.GetWeather(location, a)
			if err != nil {
				alert(fmt.Sprintf("Error getting weather data: %s", err.Error()))
				return
			}
		}

		u, err := wtw.GetClothesURL(a)
		if err != nil {
			alert(fmt.Sprintf("Error constructing URL: %s", err.Error()))
			return
		}

		document.Set("location", u.String())
	}()

	return false
}

func main() {
	js.Global.Set("getClothes", getClothes)
}
