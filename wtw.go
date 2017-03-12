package wtw

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Conditions struct {
	Gender     string `json:"gender"`
	Temp       int    `json:"temp"`
	Conditions string `json:"conditions"`
	Wind       string `json:"wind"`
	Time       string `json:"time"`
	Intensity  string `json:"intensity"`
	Feel       string `json:"feel"`
}

func (a *Conditions) Validate() error {
	switch a.Gender {
	case "m", "f":
	default:
		return fmt.Errorf("invalid value for -gender: must be m or f")
	}

	switch a.Conditions {
	case "c", "pc", "o", "r", "lr", "s":
	default:
		return fmt.Errorf("invalid value for -conditions: must be c, pc, o, r, lr or s")
	}

	switch a.Wind {
	case "nw", "lw", "hw":
	default:
		return fmt.Errorf("invalid value for -wind: must be nw, lw or hw")
	}

	switch a.Time {
	case "dawn", "day", "dusk", "night", "current":
	default:
		return fmt.Errorf("invalid value for -time: must be dawn, day, dusk or night")
	}

	switch a.Intensity {
	case "n", "lr", "h", "r":
	default:
		return fmt.Errorf("invalid value for -intensity: must be n, lr, h or r")
	}

	switch a.Feel {
	case "c", "ib", "w":
	default:
		return fmt.Errorf("invalid value for -feel: must be c, ib or w")
	}

	return nil
}

var answerRegexp = regexp.MustCompile(`<strong><a href="[^"]+">([^<]+)</a></strong>`)

func GetClothesURL(a *Conditions) (*url.URL, error) {
	u, err := url.Parse("http://www.runnersworld.com/what-to-wear")
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("gender", a.Gender)
	v.Set("temp", GetTemp(a.Temp))
	v.Set("conditions", a.Conditions)
	v.Set("wind", a.Wind)
	v.Set("time", a.Time)
	v.Set("intensity", a.Intensity)
	v.Set("feel", a.Feel)
	u.RawQuery = v.Encode()

	return u, nil
}

func GetClothes(a *Conditions) ([]string, error) {
	if a.Time == "current" {
		a.Time = GetTime()
	}

	u, err := GetClothesURL(a)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("gender", a.Gender)
	v.Set("temp", GetTemp(a.Temp))
	v.Set("conditions", a.Conditions)
	v.Set("wind", a.Wind)
	v.Set("time", a.Time)
	v.Set("intensity", a.Intensity)
	v.Set("feel", a.Feel)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no answer: %s", http.StatusText(resp.StatusCode))
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	clothes := []string{}

	ms := answerRegexp.FindAllStringSubmatch(string(buf), -1)
	if ms == nil {
		return nil, fmt.Errorf("no answer: no matches")
	}

	for _, m := range ms {
		clothes = append(clothes, string(m[1]))
	}

	return clothes, nil
}

func GetTemp(temp int) string {
	neg := 1
	if temp < 0 {
		neg = -1
	}
	temp = neg * (5 * (int(math.Abs(float64(temp))) / 5))
	tempStr := strconv.Itoa(temp)
	if tempStr == "0" {
		tempStr = "zero"
	}

	return tempStr
}

func GetTime() string {
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

func GetWeather(location string, a *Conditions) error {
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

	a.Temp = data.Query.Results.Channel.Item.Condition.Temp
	a.Wind = getWind(data.Query.Results.Channel.Wind.Speed)
	a.Conditions, err = getConditions(data.Query.Results.Channel.Item.Condition.Code)
	if err != nil {
		return err
	}

	return nil
}
