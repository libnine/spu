package fin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type DS struct {
	GTO  []dump `json:"gt_one_percent"`
	HY   []dump `json:"high_yield"`
	LTO  []dump `json:"lt_one_percent"`
	RVol []dump `json:"rel_vol"`
}

type dump struct {
	ID        string     `json:"_id"`
	Ftwl      float64    `json:"52wk_low,omitempty"`
	Ftwh      float64    `json:"52wk_high,omitempty"`
	Chg       float64    `json:"chg"`
	ChgPct    float64    `json:"chg_pct"`
	Date      *time.Time `json:"date"`
	DayHigh   float64    `json:"day_high"`
	DayLow    float64    `json:"day_low"`
	Div       float64    `json:"div,omitempty"`
	Exdiv     string     `json:"ex_div,omitempty"`
	Last      float64    `json:"last"`
	Name      string     `json:"name"`
	QuoteTime string     `json:"quoteTime,omitempty"`
	RelVol    float64    `json:"relative_volume"`
	Ticker    string     `json:"ticker"`
	Volume    float64    `json:"volume"`
	Yield     float64    `json:"yield"`
}

func Pff() {
	var data []DS
	// var stdout, stderr bytes.Buffer

	// cmd := exec.Command("/usr/bin/sh", "./scripts/sh/pff.sh")
	// cmd.Stdout, cmd.Stderr = &stdout, &stderr
	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }

	f, err := os.Open("./data/dumps/dump.json")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)

	_ = json.Unmarshal(bytes, &data)

	yr, mo, day := time.Now().Date()

	fmt.Printf("\n%s %s %v, %v\n\n", time.Now().Weekday(), mo, day, yr)
	fmt.Printf("\nDown > 1%%")
	fmt.Printf("\nTicker\t\tChange\tVolume (1000)\tRelative Volume\n")
	for n := range data[0].LTO {
		t := strings.ToUpper(data[0].LTO[n].Ticker)
		fmt.Printf("%s\t\t%0.2f%%\t%0.2f\t%0.2f\n", t, data[0].LTO[n].ChgPct, data[0].LTO[n].Volume, data[0].LTO[n].RelVol)
	}

	fmt.Printf("\nUp > 1%%")
	fmt.Printf("\nTicker\t\tChange\tVolume (1000)\tRelative Volume\n")
	for n := range data[0].GTO {
		t := strings.ToUpper(data[0].GTO[n].Ticker)
		fmt.Printf("%s\t\t%0.2f%%\t%0.2f\t%0.2f\n", t, data[0].GTO[n].ChgPct, data[0].GTO[n].Volume, data[0].GTO[n].RelVol)
	}

	fmt.Printf("\nOver 2x Average Volume")
	fmt.Printf("\nTicker\t\tChange\tVolume (1000)\tRelative Volume\n")
	for n := range data[0].RVol {
		t := strings.ToUpper(data[0].RVol[n].Ticker)
		fmt.Printf("%s\t\t%0.2f%%\t%0.2f\t%0.2f\n", t, data[0].RVol[n].ChgPct, data[0].RVol[n].Volume, data[0].RVol[n].RelVol)
	}
}
