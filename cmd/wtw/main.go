package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nwidger/wtw"
)

func main() {
	a := &wtw.Conditions{
		Gender:     "m",
		Temp:       60,
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
	flag.IntVar(&a.Temp, "temp", a.Temp, "temp (°F)")
	flag.StringVar(&a.Conditions, "conditions", a.Conditions, "c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing)")
	flag.StringVar(&a.Wind, "wind", a.Wind, "nw (no wind), lw (light wind), hw (heavy wind)")
	flag.StringVar(&a.Time, "time", a.Time, "dawn, day, dusk, night or current")
	flag.StringVar(&a.Intensity, "intensity", a.Intensity, "n (easy run), lr (long run), h (hard workout) or r (race)")
	flag.StringVar(&a.Feel, "feel", a.Feel, "c (cool), ib (in between) or w (warm)")
	flag.BoolVar(&verbose, "v", verbose, "print conditions before answer")

	flag.Parse()

	err := a.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		flag.Usage()
		os.Exit(1)
	}

	if len(location) > 0 {
		err := wtw.GetWeather(location, a)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading current conditions: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if verbose {
		fmt.Printf("wtw -gender %s -temp %d -conditions %s -wind %s -time %s -intensity %s -feel %s -v\n",
			a.Gender, a.Temp, a.Conditions,
			a.Wind, a.Time, a.Intensity, a.Feel)
	}

	clothes, err := wtw.GetClothes(a)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, c := range clothes {
		fmt.Printf("%s\n", c)
	}
}
