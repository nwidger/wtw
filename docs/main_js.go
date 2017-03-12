// +build js

package main

import (
	"log"

	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nwidger/wtw"
)

func main() {
	document := js.Global.Get("document")

	submitElem := document.Call("getElementById", "submit")
	submitElem.Call("addEventListener", "click", func() {
		go func() {
			var (
				tempInt int
				err     error
			)

			log.Println("getting element values")

			locationElem := document.Call("getElementById", "location")
			location := locationElem.Get("value").String()

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
					log.Fatal(err)
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
				wtw.GetWeather(location, a)
			}

			u, err := wtw.GetClothesURL(a)
			if err != nil {
				log.Fatal(err)
			}

			document.Set("location", u.String())
		}()
	})
}
