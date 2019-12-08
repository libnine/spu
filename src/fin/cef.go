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

type dataset struct {
	GTO []fund `json:"gt_one_percent"`
	LTO []fund `json:"lt_one_percent"`
}

type fund struct {
	ID        string     `json:"_id"`
	Ftwl      float64    `json:"52wk_low,omitempty"`
	Ftwh      float64    `json:"52wk_high,omitempty"`
	Chg       float64    `json:"chg"`
	ChgPct    float64    `json:"chg_pct"`
	Date      *time.Time `json:"date"`
	Last      float64    `json:"last"`
	Name      string     `json:"name"`
	QuoteTime string     `json:"quoteTime,omitempty"`
	Ticker    string     `json:"ticker"`
}

func CefCompile() {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("/usr/bin/sh", "./scripts/sh/cef.sh")
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func Cef(arg string) {
	var data []dataset

	f, err := os.Open("./data/dumps/cef.json")
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
		fmt.Printf("\nTicker\t\tLast\tChange\n")
		for n := range data[0].LTO {
			t := strings.ToUpper(data[0].LTO[n].Ticker)

			if len(data[0].LTO[n].Ticker) >= 8 {
				t = fmt.Sprintf("%s\t", strings.ToUpper(data[0].LTO[n].Ticker))
			} else {
				t = fmt.Sprintf("%s\t\t", strings.ToUpper(data[0].LTO[n].Ticker))
			}

			fmt.Printf("%s%0.2f\t%0.2f\n", t, data[0].LTO[n].Last, data[0].LTO[n].ChgPct)
		}
	case "up":
		fmt.Printf("\nUp > 1%%")
		fmt.Printf("\nTicker\t\tLast\tChange\n")
		for n := range data[0].GTO {
			t := strings.ToUpper(data[0].GTO[n].Ticker)

			if len(data[0].GTO[n].Ticker) >= 8 {
				t = fmt.Sprintf("%s\t", strings.ToUpper(data[0].GTO[n].Ticker))
			} else {
				t = fmt.Sprintf("%s\t\t", strings.ToUpper(data[0].GTO[n].Ticker))
			}

			fmt.Printf("%s%0.2f\t%0.2f\t\n", t, data[0].GTO[n].Last, data[0].GTO[n].ChgPct)
		}
	}
}
