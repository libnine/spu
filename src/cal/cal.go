package cal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Calendar() {
	var temp []string
	var weekset [][]string
	offset := 0

	r, err := http.Get("https://us.econoday.com/byweek.asp?cust=us")
	if err != nil {
		log.Fatal("Couldn't get calendar data.")
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	events, times := regexp.MustCompile(`[econo(events|alerts)][\s]{0,1}[a-zA-Z]{0,5}[a-zA-Z+]{0,1}\"\>\<[a-zA-Z\s\=\"\.\?\d\&\/]+\#top\"\>(.+?)\<`), regexp.MustCompile(`\d+\:\d+\s[AP][M]`)
	e, t := events.FindAllStringSubmatch(string(body), -1), times.FindAllStringSubmatch(string(body), -1)

	for i := range e {
		if strings.Contains(e[i][1], "Settlement") || strings.Contains(e[i][1], "Motor Vehicle") {
			temp = append(temp, fmt.Sprintf("%s", e[i][1]))
			offset++
			continue
		}

		temp = append(temp, fmt.Sprintf("%s %s", t[i-offset][0], e[i][1]))

		if i == len(e)-1 {
			weekset = append(weekset, temp)
			break
		}

		x, err := strconv.Atoi(strings.Split(t[i-offset][0], ":")[0])
		if err != nil {
			log.Fatal("Error populating calendar data.")
		}

		y, err := strconv.Atoi(strings.Split(t[i-offset+1][0], ":")[0])
		if err != nil {
			log.Fatal("Error populating calendar data.")
		}

		if x > y && strings.Split(t[i-offset+1][0], " ")[1] == "AM" {
			weekset = append(weekset, temp)
			temp = nil
		} else if strings.Split(t[i-offset][0], " ")[1] == "PM" && (strings.Split(t[i-offset+1][0], " ")[1] == "AM" ||
			(strings.Contains(e[i+1][1], "Settlement") || strings.Contains(e[i+1][1], "Motor Vehicle"))) {
			weekset = append(weekset, temp)
			temp = nil
		}
	}

	dn := int(time.Now().Weekday())
	yr, mo, day := time.Now().Date()

	if dn == 6 || dn == 7 {
		dn = 0
	}

	fmt.Printf("\n%s %s %v, %v\n", time.Now().Weekday(), mo, day, yr)
	for n := range weekset[dn-1] {
		fmt.Printf("\n%s", weekset[dn-1][n])
	}

	fmt.Printf("\n")
}
