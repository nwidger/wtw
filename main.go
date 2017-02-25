package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type answer struct {
	Gender     string   `json:"gender"`
	Temp       string   `json:"temp"`
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

//go:generate go-bindata -o answers.go answers.gz
func main() {
	gender, conditions, wind, time, intensity, feel :=
		"m", "c", "nw", "day", "n", "ib"
	tempInt := 60

	flag.StringVar(&gender, "gender", gender, "m (male) or f (female)")
	flag.IntVar(&tempInt, "temp", tempInt, "temp (°F)")
	flag.StringVar(&conditions, "conditions", conditions, "c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing)")
	flag.StringVar(&wind, "wind", wind, "nw (now win), lw (light wind), hw (heavy wind)")
	flag.StringVar(&time, "time", time, "dawn, day, dusk or night")
	flag.StringVar(&intensity, "intensity", intensity, "n (easy run), lr (long run), h (hard workout) or r (race)")
	flag.StringVar(&feel, "feel", feel, "c (cool), ib (in between) or w (warm)")

	flag.Parse()

	neg := 1
	if tempInt < 0 {
		neg = -1
	}
	tempInt = neg * (5 * (int(math.Abs(float64(tempInt))) / 5))
	temp := strconv.Itoa(tempInt)
	if temp == "0" {
		temp = "zero"
	}

	answers := map[string]*answer{}
	err := loadAnswers(answers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading answers: %s\n", err.Error())
		os.Exit(1)
	}

	answer, ok := answers[strings.Join([]string{
		gender, temp, conditions,
		wind, time, intensity, feel,
	}, ",")]
	if !ok {
		fmt.Println("no answer")
		os.Exit(1)
	}

	for _, c := range answer.Clothes {
		fmt.Printf("%s\n", c)
	}
}
