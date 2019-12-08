package fin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func PffCompile() {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("/usr/bin/sh", "./scripts/sh/pff.sh")
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func Pff(arg string) {
	var data []DS

	f, err := os.Open("./data/dumps/pff.json")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)

	_ = json.Unmarshal(bytes, &data)

	yr, mo, day := time.Now().Date()

	fmt.Printf("\n%s %s %v, %v\n", time.Now().Weekday(), mo, day, yr)

	switch arg {
	case "down":
		fmt.Printf("\nDown > 1%%")
		fmt.Printf("\nTicker\t\tChange\tYield\tVolume (1k)\tRelative Volume\n")
		for n := range data[0].LTO {
			t := strings.ToUpper(data[0].LTO[n].Ticker)
	
			if len(data[0].LTO[n].Ticker) >= 8 {
				t = fmt.Sprintf("%s\t", strings.ToUpper(data[0].LTO[n].Ticker))
			} else {
				t = fmt.Sprintf("%s\t\t", strings.ToUpper(data[0].LTO[n].Ticker))
			}
	
			fmt.Printf("%s%0.2f%%\t%0.2f\t%0.2f\t\t%0.2f\n", t, data[0].LTO[n].ChgPct, data[0].LTO[n].Yield, data[0].LTO[n].Volume, data[0].LTO[n].RelVol)
		}
	case "up":
		fmt.Printf("\nUp > 1%%")
		fmt.Printf("\nTicker\t\tChange\tYield\tVolume (1k)\tRelative Volume\n")
		for n := range data[0].GTO {
			t := strings.ToUpper(data[0].GTO[n].Ticker)
	
			if len(data[0].GTO[n].Ticker) >= 8 {
				t = fmt.Sprintf("%s\t", strings.ToUpper(data[0].GTO[n].Ticker))
			} else {
				t = fmt.Sprintf("%s\t\t", strings.ToUpper(data[0].GTO[n].Ticker))
			}
	
			fmt.Printf("%s%0.2f%%\t%0.2f\t%0.2f\t\t%0.2f\n", t, data[0].GTO[n].ChgPct, data[0].GTO[n].Yield, data[0].GTO[n].Volume, data[0].GTO[n].RelVol)
		}	
	case "rel":
		fmt.Printf("\nOver 2x Average Volume")
		fmt.Printf("\nTicker\t\tChange\tYield\tVolume (1k)\tRelative Volume\n")
		for n := range data[0].RVol {
			t := strings.ToUpper(data[0].RVol[n].Ticker)
	
			if len(data[0].RVol[n].Ticker) >= 8 {
				t = fmt.Sprintf("%s\t", strings.ToUpper(data[0].RVol[n].Ticker))
			} else {
				t = fmt.Sprintf("%s\t\t", strings.ToUpper(data[0].RVol[n].Ticker))
			}
	
			fmt.Printf("%s%0.2f%%\t%0.2f\t%0.2f\t\t%0.2f\n", t, data[0].RVol[n].ChgPct, data[0].RVol[n].Yield, data[0].RVol[n].Volume, data[0].RVol[n].RelVol)
		}
	case "yield":
		fmt.Printf("\nHigh Yield")
		fmt.Printf("\nTicker\t\tChange\tYield\tVolume (1k)\tRelative Volume\n")
		for n := range data[0].HY {
			t := strings.ToUpper(data[0].HY[n].Ticker)
	
			if len(data[0].HY[n].Ticker) >= 8 {
				t = fmt.Sprintf("%s\t", strings.ToUpper(data[0].HY[n].Ticker))
			} else {
				t = fmt.Sprintf("%s\t\t", strings.ToUpper(data[0].HY[n].Ticker))
			}
	
			fmt.Printf("%s%0.2f%%\t%0.2f%%\t%0.2f\t\t%0.2f\n", t, data[0].HY[n].ChgPct, data[0].HY[n].Yield, data[0].HY[n].Volume, data[0].HY[n].RelVol)
		}
	}
}
