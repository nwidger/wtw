package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/nwidger/wtw"
)

func main() {
	a := &wtw.Answer{
		Gender:     "m",
		TempInt:    60,
		Conditions: "c",
		Wind:       "nw",
		Time:       "current",
		Intensity:  "n",
		Feel:       "ib",
	}
	location := ""
	verbose := false

	flag.StringVar(&location, "location", location, "get current conditions for location, overrides -temp, -conditions and -wind")
	flag.StringVar(&a.Gender, "gender", a.Gender, "m (male) or f (female)")
	flag.IntVar(&a.TempInt, "temp", a.TempInt, "temp (°F)")
	flag.StringVar(&a.Conditions, "conditions", a.Conditions, "c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing)")
	flag.StringVar(&a.Wind, "wind", a.Wind, "nw (no wind), lw (light wind), hw (heavy wind)")
	flag.StringVar(&a.Time, "time", a.Time, "dawn, day, dusk, night or current")
	flag.StringVar(&a.Intensity, "intensity", a.Intensity, "n (easy run), lr (long run), h (hard workout) or r (race)")
	flag.StringVar(&a.Feel, "feel", a.Feel, "c (cool), ib (in between) or w (warm)")
	flag.BoolVar(&verbose, "v", verbose, "print conditions before answer")

	flag.Parse()

	if a.Time == "current" {
		a.Time = wtw.GetTime()
	}

	if len(location) > 0 {
		err := wtw.GetWeather(location, a)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading current conditions: %s\n", err.Error())
			os.Exit(1)
		}
	}

	neg := 1
	if a.TempInt < 0 {
		neg = -1
	}
	a.TempInt = neg * (5 * (int(math.Abs(float64(a.TempInt))) / 5))
	a.Temp = strconv.Itoa(a.TempInt)
	if a.Temp == "0" {
		a.Temp = "zero"
	}

	if verbose {
		fmt.Printf("gender %s temp %s conditions %s wind %s time %s intensity %s feel %s\n",
			a.Gender, a.Temp, a.Conditions,
			a.Wind, a.Time, a.Intensity, a.Feel)
	}

	clothes, err := wtw.GetAnswer(a)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, c := range clothes {
		fmt.Printf("%s\n", c)
	}
}
