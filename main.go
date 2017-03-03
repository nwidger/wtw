package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type answer struct {
	Gender     string   `json:"gender"`
	Temp       string   `json:"temp"`
	TempInt    int      `json:"-"`
	Conditions string   `json:"conditions"`
	Wind       string   `json:"wind"`
	Time       string   `json:"time"`
	Intensity  string   `json:"intensity"`
	Feel       string   `json:"feel"`
	Clothes    []string `json:"clothes"`
}

func loadAnswers(answers map[string]*answer) error {
	buf, err := Asset("answers.gz")
	if err != nil {
		return err
	}

	gr, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	dec := json.NewDecoder(gr)

	err = dec.Decode(&answers)
	if err != nil {
		return err
	}

	return nil
}

func getTime() string {
	now := time.Now()
	hour := now.Hour()
	switch {
	// 5am - 6am
	case hour >= 4 && hour <= 5:
		return "dawn"
	// 7am - 5pm
	case hour >= 6 && hour <= 16:
		return "day"
	// 6pm - 7pm
	case hour >= 17 && hour <= 18:
		return "dusk"
	// 8pm - 12pm, 1am - 4am
	case (hour >= 19 && hour <= 23) || (hour >= 0 && hour <= 3):
		return "night"
	// should never get here
	default:
		return "unknown"
	}
}

func getWind(speed int) string {
	switch {
	case speed >= 0 && speed <= 3:
		return "nw"
	case speed >= 4 && speed <= 8:
		return "lw"
	default:
		return "hw"
	}
}

// Yahoo Weather API condition codes, see
// https://developer.yahoo.com/weather/documentation.html#codes
//
// c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing)
var conditionCodes = map[int]string{
	0:  "r",  // tornado
	1:  "r",  // tropical storm
	2:  "r",  // hurricane
	3:  "r",  // severe thunderstorms
	4:  "r",  // thunderstorms
	5:  "s",  // mixed rain and snow
	6:  "r",  // mixed rain and sleet
	7:  "s",  // mixed snow and sleet
	8:  "lr", // freezing drizzle
	9:  "lr", // drizzle
	10: "r",  // freezing rain
	11: "r",  // showers
	12: "r",  // showers
	13: "s",  // snow flurries
	14: "s",  // light snow showers
	15: "s",  // blowing snow
	16: "s",  // snow
	17: "s",  // hail
	18: "s",  // sleet
	19: "c",  // dust
	20: "c",  // foggy
	21: "c",  // haze
	22: "c",  // smoky
	23: "c",  // blustery
	24: "c",  // windy
	25: "c",  // cold
	26: "o",  // cloudy
	27: "o",  // mostly cloudy (night)
	28: "o",  // mostly cloudy (day)
	29: "pc", // partly cloudy (night)
	30: "pc", // partly cloudy (day)
	31: "c",  // clear (night)
	32: "c",  // sunny
	33: "c",  // fair (night)
	34: "c",  // fair (day)
	35: "r",  // mixed rain and hail
	36: "c",  // hot
	37: "r",  // isolated thunderstorms
	38: "r",  // scattered thunderstorms
	39: "r",  // scattered thunderstorms
	40: "lr", // scattered showers
	41: "s",  // heavy snow
	42: "s",  // scattered snow showers
	43: "s",  // heavy snow
	44: "pc", // partly cloudy
	45: "r",  // thundershowers
	46: "s",  // snow showers
	47: "r",  // isolated thundershowers
	// 3200: "", // not available
}

func getConditions(code int) (string, error) {
	if code == 3200 { // not available
		return "", fmt.Errorf("conditions not available")
	}
	condition, ok := conditionCodes[code]
	if !ok {
		return "", fmt.Errorf("unknown condition code %d", code)
	}
	return condition, nil
}

func getWeather(location string, a *answer) error {
	u, err := url.Parse("https://query.yahooapis.com/v1/public/yql")
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("q", fmt.Sprintf(`select * from weather.forecast where woeid in (select woeid from geo.places(1) where text="%s")`, location))
	v.Set("format", "json")
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", http.StatusText(resp.StatusCode))
	}

	data := struct {
		Query struct {
			Results struct {
				Channel struct {
					Wind struct {
						Speed int `json:"speed,string"`
					} `json:"wind"`
					Item struct {
						Condition struct {
							Code int `json:"code,string"`
							Temp int `json:"temp,string"`
						} `json:"condition"`
					} `json:"item"`
				} `json:"channel"`
			} `json:"results"`
		} `json:"query"`
	}{}

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&data)
	if err != nil {
		return err
	}

	a.TempInt = data.Query.Results.Channel.Item.Condition.Temp
	a.Wind = getWind(data.Query.Results.Channel.Wind.Speed)
	a.Conditions, err = getConditions(data.Query.Results.Channel.Item.Condition.Code)
	if err != nil {
		return err
	}

	return nil
}

//go:generate go-bindata -o answers.go answers.gz
func main() {
	a := &answer{
		Gender:     "m",
		TempInt:    60,
		Conditions: "c",
		Wind:       "nw",
		Time:       "current",
		Intensity:  "n",
		Feel:       "ib",
	}
	location := "Portsmouth, NH"

	flag.StringVar(&location, "location", location, "get current conditions for location, overrides -temp, -conditions, -wind")
	flag.StringVar(&a.Gender, "gender", a.Gender, "m (male) or f (female)")
	flag.IntVar(&a.TempInt, "temp", a.TempInt, "temp (°F)")
	flag.StringVar(&a.Conditions, "conditions", a.Conditions, "c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing)")
	flag.StringVar(&a.Wind, "wind", a.Wind, "nw (now win), lw (light wind), hw (heavy wind)")
	flag.StringVar(&a.Time, "time", a.Time, "dawn, day, dusk, night or current")
	flag.StringVar(&a.Intensity, "intensity", a.Intensity, "n (easy run), lr (long run), h (hard workout) or r (race)")
	flag.StringVar(&a.Feel, "feel", a.Feel, "c (cool), ib (in between) or w (warm)")
	verbose := flag.Bool("v", false, "print conditions before answer")

	flag.Parse()

	if a.Time == "current" {
		a.Time = getTime()
	}

	if len(location) > 0 {
		err := getWeather(location, a)
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

	answers := map[string]*answer{}
	err := loadAnswers(answers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading answers: %s\n", err.Error())
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("wtw -gender %s -temp %s -conditions %s -wind %s -time %s -intensity %s -feel %s\n",
			a.Gender, a.Temp, a.Conditions,
			a.Wind, a.Time, a.Intensity, a.Feel)
	}

	answer, ok := answers[strings.Join([]string{
		a.Gender, a.Temp, a.Conditions,
		a.Wind, a.Time, a.Intensity, a.Feel,
	}, ",")]
	if !ok {
		fmt.Println("no answer")
		os.Exit(1)
	}

	for _, c := range answer.Clothes {
		fmt.Printf("%s\n", c)
	}
}
